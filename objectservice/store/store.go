package store

import (
	"context"
	"fmt"
	dagpoolcli "github.com/filedag-project/filedag-storage/dag/pool/client"
	"github.com/filedag-project/filedag-storage/objectservice/consts"
	"github.com/filedag-project/filedag-storage/objectservice/datatypes"
	"github.com/filedag-project/filedag-storage/objectservice/lock"
	"github.com/filedag-project/filedag-storage/objectservice/objmetadb"
	"github.com/filedag-project/filedag-storage/objectservice/utils"
	"github.com/filedag-project/filedag-storage/objectservice/utils/s3utils"
	"github.com/google/uuid"
	"github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	logging "github.com/ipfs/go-log/v2"
	"github.com/ipfs/go-merkledag"
	ufsio "github.com/ipfs/go-unixfs/io"
	"github.com/klauspost/readahead"
	pool "github.com/libp2p/go-buffer-pool"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/syndtr/goleveldb/leveldb"
	"golang.org/x/xerrors"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var log = logging.Logger("store")

//storageSys store sys
type storageSys struct {
	Db              objmetadb.ObjStoreMetaDBAPI
	DagPool         ipld.DAGService
	CidBuilder      cid.Builder
	nsLock          *lock.NsLockMap
	newBucketNSLock func(bucket string) lock.RWLocker
	hasBucket       func(ctx context.Context, bucket string) bool

	gcPeriod  time.Duration
	gcTimeout time.Duration
}

//canCreateFolder check if can create folder
func (s *storageSys) canCreateFolder(ctx context.Context, bucket string, folder string) bool {
	lk := s.newNSLock(bucket, folder)
	lkctx, err := lk.GetRLock(ctx, globalOperationTimeout)
	if err != nil {
		return false
	}
	ctx = lkctx.Context()
	defer lk.RUnlock(lkctx.Cancel)
	if strings.Count(folder, "/") > 1 {
		folders := strings.SplitAfterN(folder, "/", 2)
		last := folders[len(folders)-1]
		aa := folder[:len(folder)-len(last)]
		if _, err := s.getObjectInfo(ctx, bucket, aa); err != nil {
			return false
		}
		folder = folders[len(folders)-1]
	}
	_, err = s.getObjectInfo(ctx, bucket, folder[:strings.Index(folder, "/")])
	if err != nil {
		return true
	} else {
		return false
	}
}

// StoreStats store system stats
func (s *storageSys) StoreStats(ctx context.Context) (DataUsageInfo, error) {
	var dataUsageInfo DataUsageInfo
	buckets := s.GetAllBucket(ctx)
	dataUsageInfo.BucketsCount = uint64(len(buckets))
	dataUsageInfo.TotalCaptivity = defaultTotalCaptivity
	for _, bucket := range buckets {
		info, err := s.GetBucketInfo(ctx, bucket)
		if err != nil {
			continue
		}
		dataUsageInfo.BucketsUsage = append(dataUsageInfo.BucketsUsage, info)
		dataUsageInfo.ObjectsTotalCount += info.Objects
		dataUsageInfo.ObjectsTotalSize += info.Size
	}
	return dataUsageInfo, nil
}

//NewStorageSys new a storage sys
func NewStorageSys(ctx context.Context, dagService ipld.DAGService, db objmetadb.ObjStoreMetaDBAPI) ObjectStoreSystemAPI {
	cidBuilder, _ := merkledag.PrefixForCidVersion(0)
	s := &storageSys{
		Db:         db,
		DagPool:    dagService,
		CidBuilder: cidBuilder,
		nsLock:     lock.NewNSLock(),
		gcPeriod:   15 * time.Minute,
		gcTimeout:  30 * time.Minute,
	}
	go func() {
		s.processObjectGC(ctx)
	}()
	return s
}

func getObjectKey(bucket, object string) string {
	return fmt.Sprintf(objectKeyFormat, bucket, object)
}

func getUploadKey(bucket, object, uploadID string) string {
	return fmt.Sprintf(uploadKeyFormat, bucket, object, uploadID)
}

func newDelObjectKey() string {
	return fmt.Sprintf(deleteKeyFormat, mustGetUUID())
}

// newNSLock - initialize a new namespace RWLocker instance.
func (s *storageSys) newNSLock(bucket string, objects ...string) lock.RWLocker {
	return s.nsLock.NewNSLock(bucket, objects...)
}

func (s *storageSys) SetNewBucketNSLock(newBucketNSLock func(bucket string) lock.RWLocker) {
	s.newBucketNSLock = newBucketNSLock
}

func (s *storageSys) SetHasBucket(hasBucket func(ctx context.Context, bucket string) bool) {
	s.hasBucket = hasBucket
}

func (s *storageSys) store(ctx context.Context, reader io.ReadCloser, size int64) (cid.Cid, error) {
	data := io.Reader(reader)
	if size > bigFileThreshold {
		// We use 2 buffers, so we always have a full buffer of input.
		bufA := pool.Get(chunkSize)
		bufB := pool.Get(chunkSize)
		defer pool.Put(bufA)
		defer pool.Put(bufB)
		ra, err := readahead.NewReaderBuffer(data, [][]byte{bufA[:chunkSize], bufB[:chunkSize]})
		if err == nil {
			data = ra
			defer ra.Close()
		} else {
			log.Infof("readahead.NewReaderBuffer failed, error: %v", err)
		}
	}
	node, err := dagpoolcli.BalanceNode(data, s.DagPool, s.CidBuilder)
	if err != nil {
		return cid.Undef, err
	}
	select {
	case <-ctx.Done():
		return cid.Undef, ctx.Err()
	default:
	}
	return node.Cid(), nil
}

func (s *storageSys) checkAndDeleteObjectData(ctx context.Context, bucket, object string) {
	if oldObjInfo, err := s.getObjectInfo(ctx, bucket, object); err == nil {
		c, err := cid.Decode(oldObjInfo.ETag)
		if err != nil {
			log.Warnw("decode cid error", "cid", oldObjInfo.ETag)
		} else {
			if err = s.markObjetToDelete(c); err != nil {
				log.Errorw("mark Objet to delete error", "bucket", bucket, "object", object, "cid", oldObjInfo.ETag, "error", err)
			}
		}
	}
}

//StoreObject store object
func (s *storageSys) StoreObject(ctx context.Context, bucket, object string, reader io.ReadCloser, size int64, meta map[string]string, fileFolder bool) (ObjectInfo, error) {
	bktlk := s.newBucketNSLock(bucket)
	bktlkCtx, err := bktlk.GetRLock(ctx, globalOperationTimeout)
	if err != nil {
		return ObjectInfo{}, err
	}
	ctx = bktlkCtx.Context()
	defer bktlk.RUnlock(bktlkCtx.Cancel)

	if !s.hasBucket(ctx, bucket) {
		return ObjectInfo{}, BucketNotFound{Bucket: bucket}
	}
	var root cid.Cid
	if !fileFolder {
		root, err = s.store(ctx, reader, size)
		if err != nil {
			return ObjectInfo{}, err
		}
	} else {
		if !s.canCreateFolder(ctx, bucket, object) {
			return ObjectInfo{}, ErrCanNotCreatFolder
		}
	}
	var etag = func(root cid.Cid) string {
		if root != cid.Undef {
			return root.String()
		} else {
			return cid.Undef.String()
		}
	}(root)
	objInfo := ObjectInfo{
		Bucket:           bucket,
		Name:             object,
		ModTime:          time.Now().UTC(),
		Size:             size,
		IsDir:            false,
		ETag:             etag,
		VersionID:        "",
		IsLatest:         true,
		DeleteMarker:     false,
		ContentType:      meta[strings.ToLower(consts.ContentType)],
		ContentEncoding:  meta[strings.ToLower(consts.ContentEncoding)],
		SuccessorModTime: time.Now().UTC(),
	}
	// Update expires
	if exp, ok := meta[strings.ToLower(consts.Expires)]; ok {
		if t, e := time.Parse(http.TimeFormat, exp); e == nil {
			objInfo.Expires = t.UTC()
		}
	}

	lk := s.newNSLock(bucket, object)
	lkctx, err := lk.GetLock(ctx, globalOperationTimeout)
	if err != nil {
		return ObjectInfo{}, err
	}
	ctx = lkctx.Context()
	defer lk.Unlock(lkctx.Cancel)

	// Has old file?
	s.checkAndDeleteObjectData(ctx, bucket, object)

	err = s.Db.Put(getObjectKey(bucket, object), objInfo)
	if err != nil {
		return ObjectInfo{}, err
	}
	return objInfo, nil
}

//GetObject Get object
func (s *storageSys) GetObject(ctx context.Context, bucket, object string) (ObjectInfo, io.ReadCloser, error) {
	lk := s.newNSLock(bucket, object)
	lkctx, err := lk.GetRLock(ctx, globalOperationTimeout)
	if err != nil {
		return ObjectInfo{}, nil, err
	}
	ctx = lkctx.Context()
	defer lk.RUnlock(lkctx.Cancel)

	meta, err := s.getObjectInfo(ctx, bucket, object)
	if err != nil {
		return ObjectInfo{}, nil, err
	}
	cid, err := cid.Decode(meta.ETag)
	if err != nil {
		return ObjectInfo{}, nil, err
	}
	dagNode, err := s.DagPool.Get(ctx, cid)
	if err != nil {
		return ObjectInfo{}, nil, err
	}
	reader, err := ufsio.NewDagReader(ctx, dagNode, s.DagPool)
	if err != nil {
		return ObjectInfo{}, nil, err
	}
	return meta, reader, nil
}

func (s *storageSys) getObjectInfo(ctx context.Context, bucket, object string) (meta ObjectInfo, err error) {
	err = s.Db.Get(getObjectKey(bucket, object), &meta)
	if err != nil {
		if xerrors.Is(err, leveldb.ErrNotFound) {
			return meta, ErrObjectNotFound
		}
	}
	return
}

func (s *storageSys) GetObjectInfo(ctx context.Context, bucket, object string) (meta ObjectInfo, err error) {
	lk := s.newNSLock(bucket, object)
	lkctx, err := lk.GetRLock(ctx, globalOperationTimeout)
	if err != nil {
		return ObjectInfo{}, err
	}
	ctx = lkctx.Context()
	defer lk.RUnlock(lkctx.Cancel)

	return s.getObjectInfo(ctx, bucket, object)
}

//DeleteObject delete object
func (s *storageSys) DeleteObject(ctx context.Context, bucket, object string) error {
	lk := s.newNSLock(bucket, object)
	lkctx, err := lk.GetLock(ctx, deleteOperationTimeout)
	if err != nil {
		return err
	}
	ctx = lkctx.Context()
	defer lk.Unlock(lkctx.Cancel)

	meta, err := s.getObjectInfo(ctx, bucket, object)
	if err != nil {
		return err
	}
	if strings.HasSuffix(meta.Name, "/") {
		allChan, err := s.Db.ReadAllChan(ctx, fmt.Sprintf(objectKeyFormat, bucket, object), "")
		if err != nil {
			return err
		}
		for entry := range allChan {
			if err = s.Db.Delete(entry.GetKey()); err != nil {
				return err
			}
		}
	} else {
		cid, err := cid.Decode(meta.ETag)
		if err != nil {
			return err
		}
		if err = s.markObjetToDelete(cid); err != nil {
			log.Errorw("mark Objet to delete error", "bucket", bucket, "object", object, "cid", meta.ETag, "error", err)
		}
		if err = s.Db.Delete(getObjectKey(bucket, object)); err != nil {
			return err
		}
	}
	return nil
}

func (s *storageSys) CleanObjectsInBucket(ctx context.Context, bucket string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	prefixKey := fmt.Sprintf(allObjectPrefixFormat, bucket, "")
	all, err := s.Db.ReadAllChan(ctx, prefixKey, "")
	if err != nil {
		return err
	}
	for entry := range all {
		var o ObjectInfo
		if err = entry.UnmarshalValue(&o); err != nil {
			return err
		}
		if err = s.DeleteObject(ctx, bucket, o.Name); err != nil {
			return err
		}
	}
	return nil
}

//GetBucketInfo Get BucketInfo
func (s *storageSys) GetBucketInfo(ctx context.Context, bucket string) (bi BucketInfo, err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	seekKey := ""
	prefixKey := fmt.Sprintf(allObjectPrefixFormat, bucket, "")
	all, err := s.Db.ReadAllChan(ctx, prefixKey, seekKey)
	if err != nil {
		return bi, err
	}
	var size, objects uint64
	for entry := range all {
		var o ObjectInfo
		if err = entry.UnmarshalValue(&o); err != nil {
			return bi, err
		}
		size += uint64(o.Size)
		objects++

	}
	return BucketInfo{
		Name:    bucket,
		Size:    size,
		Objects: objects,
	}, nil
}

//ListObjects list user object
//TODO use more params
func (s *storageSys) ListObjects(ctx context.Context, bucket string, prefix string, marker string, delimiter string, maxKeys int) (loi ListObjectsInfo, err error) {
	if maxKeys == 0 {
		return loi, nil
	}

	if len(prefix) > 0 && maxKeys == 1 && delimiter == "" && marker == "" {
		// Optimization for certain applications like
		// - Cohesity
		// - Actifio, Splunk etc.
		// which send ListObjects requests where the actual object
		// itself is the prefix and max-keys=1 in such scenarios
		// we can simply verify locally if such an object exists
		// to avoid the need for ListObjects().
		objInfo, err := s.GetObjectInfo(ctx, bucket, prefix)
		if err == nil {
			loi.Objects = append(loi.Objects, objInfo)
			return loi, nil
		}
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	seekKey := ""
	if marker != "" {
		seekKey = fmt.Sprintf(allObjectSeekKeyFormat, bucket, marker)
	}
	prefixKey := fmt.Sprintf(allObjectPrefixFormat, bucket, prefix)
	all, err := s.Db.ReadAllChan(ctx, prefixKey, seekKey)
	if err != nil {
		return loi, err
	}
	index := 0
	m := make(map[string]struct{})
	for entry := range all {
		if index == maxKeys {
			loi.IsTruncated = true
			break
		}
		var o ObjectInfo
		if err = entry.UnmarshalValue(&o); err != nil {
			return loi, err
		}
		index++
		// bucket/aaa/ccc/
		getFolder(o, prefix, &loi, m)
	}
	if loi.IsTruncated {
		loi.NextMarker = loi.Objects[len(loi.Objects)-1].Name
	}
	for key := range m {
		loi.Prefixes = append(loi.Prefixes, key)
	}
	return loi, nil
}
func getFolder(o ObjectInfo, prefix string, loi *ListObjectsInfo, m map[string]struct{}) {
	// bucket/aaa/ccc/
	if strings.HasSuffix(o.Name, "/") {
		if len(o.Name) > len(prefix) {
			name := o.Name[len(prefix):]
			m[name[:strings.Index(name, "/")+1]] = struct{}{}
		}
	} else {
		loi.Objects = append(loi.Objects, o)
	}
}
func (s *storageSys) EmptyBucket(ctx context.Context, bucket string) (bool, error) {
	loi, err := s.ListObjects(ctx, bucket, "", "", "", 1)
	if err != nil {
		return false, err
	}
	return len(loi.Objects) == 0, nil
}

// ListObjectsV2 list objects
func (s *storageSys) ListObjectsV2(ctx context.Context, bucket string, prefix string, continuationToken string, delimiter string, maxKeys int, owner bool, startAfter string) (ListObjectsV2Info, error) {
	marker := continuationToken
	if marker == "" {
		marker = startAfter
	}
	loi, err := s.ListObjects(ctx, bucket, prefix, marker, delimiter, maxKeys)
	if err != nil {
		return ListObjectsV2Info{}, err
	}
	listV2Info := ListObjectsV2Info{
		IsTruncated:           loi.IsTruncated,
		ContinuationToken:     continuationToken,
		NextContinuationToken: loi.NextMarker,
		Objects:               loi.Objects,
		Prefixes:              loi.Prefixes,
	}
	return listV2Info, nil
}

// mustGetUUID - get a random UUID.
func mustGetUUID() string {
	u, err := uuid.NewRandom()
	if err != nil {
		log.Fatal(err)
	}

	return u.String()
}

func (s *storageSys) NewMultipartUpload(ctx context.Context, bucket string, object string, meta map[string]string) (MultipartInfo, error) {
	bktlk := s.newBucketNSLock(bucket)
	bktlkCtx, err := bktlk.GetRLock(ctx, globalOperationTimeout)
	if err != nil {
		return MultipartInfo{}, err
	}
	ctx = bktlkCtx.Context()
	defer bktlk.RUnlock(bktlkCtx.Cancel)

	if !s.hasBucket(ctx, bucket) {
		return MultipartInfo{}, BucketNotFound{Bucket: bucket}
	}

	// uploadId is random, so don't to lock it
	uploadId := mustGetUUID()
	info := MultipartInfo{
		Bucket:    bucket,
		Object:    object,
		UploadID:  uploadId,
		MetaData:  meta,
		Initiated: time.Now().UTC(),
	}

	err = s.Db.Put(getUploadKey(bucket, object, uploadId), info)
	if err != nil {
		return MultipartInfo{}, err
	}
	return info, nil
}

func (s *storageSys) GetMultipartInfo(ctx context.Context, bucket string, object string, uploadID string) (MultipartInfo, error) {
	bktlk := s.newBucketNSLock(bucket)
	bktlkCtx, err := bktlk.GetRLock(ctx, globalOperationTimeout)
	if err != nil {
		return MultipartInfo{}, err
	}
	ctx = bktlkCtx.Context()
	defer bktlk.RUnlock(bktlkCtx.Cancel)

	uploadIDLock := s.newNSLock(bucket, lock.PathJoin(object, uploadID))
	lkctx, err := uploadIDLock.GetRLock(ctx, globalOperationTimeout)
	if err != nil {
		return MultipartInfo{}, err
	}
	ctx = lkctx.Context()
	defer uploadIDLock.RUnlock(lkctx.Cancel)

	return s.getMultipartInfo(ctx, bucket, object, uploadID)
}

func (s *storageSys) getMultipartInfo(ctx context.Context, bucket string, object string, uploadID string) (MultipartInfo, error) {
	info := MultipartInfo{}
	err := s.Db.Get(getUploadKey(bucket, object, uploadID), &info)
	return info, err
}

func (s *storageSys) PutObjectPart(ctx context.Context, bucket string, object string, uploadID string, partID int, reader io.ReadCloser, size int64, meta map[string]string) (pi objectPartInfo, err error) {
	bktlk := s.newBucketNSLock(bucket)
	bktlkCtx, err := bktlk.GetRLock(ctx, globalOperationTimeout)
	if err != nil {
		return pi, err
	}
	ctx = bktlkCtx.Context()
	defer bktlk.RUnlock(bktlkCtx.Cancel)

	root, err := s.store(ctx, reader, size)
	if err != nil {
		return pi, err
	}

	partInfo := objectPartInfo{
		Number:  partID,
		ETag:    root.String(),
		Size:    size,
		ModTime: time.Now().UTC(),
	}

	uploadIDLock := s.newNSLock(bucket, lock.PathJoin(object, uploadID))
	ulkctx, err := uploadIDLock.GetLock(ctx, globalOperationTimeout)
	if err != nil {
		return pi, err
	}
	ctx = ulkctx.Context()
	defer uploadIDLock.Unlock(ulkctx.Cancel)

	mi, err := s.getMultipartInfo(ctx, bucket, object, uploadID)
	if err != nil {
		return pi, err
	}

	mi.Parts = append(mi.Parts, partInfo)
	err = s.Db.Put(getUploadKey(bucket, object, uploadID), mi)
	if err != nil {
		return pi, err
	}
	return partInfo, nil
}

func (s *storageSys) removeMultipartInfo(ctx context.Context, bucket string, object string, uploadID string) error {
	return s.Db.Delete(getUploadKey(bucket, object, uploadID))
}

// objectPartIndex - returns the index of matching object part number.
func objectPartIndex(parts []objectPartInfo, partNumber int) int {
	for i, part := range parts {
		if partNumber == part.Number {
			return i
		}
	}
	return -1
}

var etagRegex = regexp.MustCompile("\"*?([^\"]*?)\"*?$")

// canonicalizeETag returns ETag with leading and trailing double-quotes removed,
// if any present
func canonicalizeETag(etag string) string {
	return etagRegex.ReplaceAllString(etag, "$1")
}

func (s *storageSys) CompleteMultiPartUpload(ctx context.Context, bucket string, object string, uploadID string, parts []datatypes.CompletePart) (oi ObjectInfo, err error) {
	bktlk := s.newBucketNSLock(bucket)
	bktlkCtx, err := bktlk.GetRLock(ctx, globalOperationTimeout)
	if err != nil {
		return oi, err
	}
	ctx = bktlkCtx.Context()
	defer bktlk.RUnlock(bktlkCtx.Cancel)

	if !s.hasBucket(ctx, bucket) {
		return oi, BucketNotFound{Bucket: bucket}
	}

	uploadIDLock := s.newNSLock(bucket, lock.PathJoin(object, uploadID))
	ulkctx, err := uploadIDLock.GetLock(ctx, globalOperationTimeout)
	if err != nil {
		return oi, err
	}
	ctx = ulkctx.Context()
	defer uploadIDLock.Unlock(ulkctx.Cancel)

	mi, err := s.getMultipartInfo(ctx, bucket, object, uploadID)
	if err != nil {
		return oi, err
	}

	var objectSize int64
	var links []dagpoolcli.LinkInfo
	for i, part := range parts {
		partIndex := objectPartIndex(mi.Parts, part.PartNumber)
		if partIndex < 0 {
			invp := s3utils.InvalidPart{
				PartNumber: part.PartNumber,
				GotETag:    part.ETag,
			}
			return oi, invp
		}
		gotPart := mi.Parts[partIndex]

		// ensure that part ETag is canonicalized to strip off extraneous quotes
		part.ETag = canonicalizeETag(part.ETag)
		if gotPart.ETag != part.ETag {
			invp := s3utils.InvalidPart{
				PartNumber: part.PartNumber,
				ExpETag:    gotPart.ETag,
				GotETag:    part.ETag,
			}
			return oi, invp
		}

		// All parts except the last part has to be at least 5MB.
		if (i < len(parts)-1) && !(gotPart.Size >= consts.MinPartSize) {
			return oi, s3utils.PartTooSmall{
				PartNumber: part.PartNumber,
				PartSize:   gotPart.Size,
				PartETag:   part.ETag,
			}
		}

		// Save for total object size.
		objectSize += gotPart.Size

		c, err := cid.Decode(gotPart.ETag)
		if err != nil {
			return oi, err
		}
		linkInfo, err := dagpoolcli.CreateLinkInfo(ctx, s.DagPool, c)
		if err != nil {
			return oi, err
		}
		links = append(links, linkInfo)
	}
	root, err := dagpoolcli.BuildDataCidByLinks(ctx, s.DagPool, s.CidBuilder, links)
	if err != nil {
		return oi, err
	}
	objInfo := ObjectInfo{
		Bucket:           bucket,
		Name:             object,
		ModTime:          time.Now().UTC(),
		Size:             objectSize,
		IsDir:            false,
		ETag:             root.String(),
		VersionID:        "",
		IsLatest:         true,
		DeleteMarker:     false,
		ContentType:      mi.MetaData[strings.ToLower(consts.ContentType)],
		ContentEncoding:  mi.MetaData[strings.ToLower(consts.ContentEncoding)],
		SuccessorModTime: time.Now().UTC(),
	}
	// Update expires
	if exp, ok := mi.MetaData[strings.ToLower(consts.Expires)]; ok {
		if t, e := time.Parse(http.TimeFormat, exp); e == nil {
			objInfo.Expires = t.UTC()
		}
	}

	lk := s.newNSLock(bucket, object)
	lkctx, err := lk.GetLock(ctx, globalOperationTimeout)
	if err != nil {
		return ObjectInfo{}, err
	}
	ctx = lkctx.Context()
	defer lk.Unlock(lkctx.Cancel)

	// Has old file?
	s.checkAndDeleteObjectData(ctx, bucket, object)

	err = s.Db.Put(getObjectKey(bucket, object), objInfo)
	if err != nil {
		return ObjectInfo{}, err
	}

	// remove MultipartInfo
	err = s.removeMultipartInfo(ctx, bucket, object, uploadID)
	if err != nil {
		log.Errorw("remove MultipartInfo error", "bucket", bucket, "object", object, "uploadID", uploadID, "error", err)
	}
	return objInfo, nil
}

func (s *storageSys) AbortMultipartUpload(ctx context.Context, bucket string, object string, uploadID string) error {
	bktlk := s.newBucketNSLock(bucket)
	bktlkCtx, err := bktlk.GetRLock(ctx, globalOperationTimeout)
	if err != nil {
		return err
	}
	ctx = bktlkCtx.Context()
	defer bktlk.RUnlock(bktlkCtx.Cancel)

	if !s.hasBucket(ctx, bucket) {
		return BucketNotFound{Bucket: bucket}
	}

	uploadIDLock := s.newNSLock(bucket, lock.PathJoin(object, uploadID))
	ulkctx, err := uploadIDLock.GetLock(ctx, globalOperationTimeout)
	if err != nil {
		return err
	}
	ctx = ulkctx.Context()
	defer uploadIDLock.Unlock(ulkctx.Cancel)

	mi, err := s.getMultipartInfo(ctx, bucket, object, uploadID)
	if err != nil {
		return err
	}

	for _, part := range mi.Parts {
		c, err := cid.Decode(part.ETag)
		if err != nil {
			return err
		}

		if err = s.markObjetToDelete(c); err != nil {
			log.Errorw("mark Objet to delete error", "bucket", bucket, "object", object, "cid", part.ETag, "error", err)
		}
	}

	// remove MultipartInfo
	err = s.removeMultipartInfo(ctx, bucket, object, uploadID)
	if err != nil {
		log.Errorw("remove MultipartInfo error", "bucket", bucket, "object", object, "uploadID", uploadID, "error", err)
	}
	return nil
}

func (s *storageSys) ListObjectParts(ctx context.Context, bucket, object, uploadID string, partNumberMarker, maxParts int) (result ListPartsInfo, err error) {
	bktlk := s.newBucketNSLock(bucket)
	bktlkCtx, err := bktlk.GetRLock(ctx, globalOperationTimeout)
	if err != nil {
		return result, err
	}
	ctx = bktlkCtx.Context()
	defer bktlk.RUnlock(bktlkCtx.Cancel)

	if !s.hasBucket(ctx, bucket) {
		return result, BucketNotFound{Bucket: bucket}
	}

	uploadIDLock := s.newNSLock(bucket, lock.PathJoin(object, uploadID))
	ulkctx, err := uploadIDLock.GetRLock(ctx, globalOperationTimeout)
	if err != nil {
		return result, err
	}
	ctx = ulkctx.Context()
	defer uploadIDLock.RUnlock(ulkctx.Cancel)

	mi, err := s.getMultipartInfo(ctx, bucket, object, uploadID)
	if err != nil {
		return result, err
	}

	if maxParts == 0 {
		return result, nil
	}

	if partNumberMarker < 0 {
		partNumberMarker = 0
	}

	// Limit output to maxPartsList.
	if maxParts > consts.MaxPartsList-partNumberMarker {
		maxParts = consts.MaxPartsList - partNumberMarker
	}

	result.Bucket = bucket
	result.Object = object
	result.UploadID = uploadID
	result.MaxParts = maxParts
	result.PartNumberMarker = partNumberMarker
	result.Metadata = utils.CloneMapSS(mi.MetaData)

	start := partNumberMarker + 1
	end := start + maxParts
	if len(mi.Parts) <= start {
		return result, nil
	}
	if end > len(mi.Parts) {
		end = len(mi.Parts)
	}
	parts := mi.Parts[start:end]

	if len(parts) == 0 || maxParts == 0 {
		return result, nil
	}

	result.Parts = parts

	// If listed entries are more than maxParts, we set IsTruncated as true.
	if len(mi.Parts)-1 > end {
		result.IsTruncated = true
		// Make sure to fill next part number marker if IsTruncated is
		// true for subsequent listing.
		result.NextPartNumberMarker = result.Parts[len(result.Parts)-1].Number
	}
	return
}

// Lookup - returns if uploadID is valid
func (lm ListMultipartsInfo) Lookup(uploadID string) bool {
	for _, upload := range lm.Uploads {
		if upload.UploadID == uploadID {
			return true
		}
	}
	return false
}

func (s *storageSys) ListMultipartUploads(ctx context.Context, bucket, prefix, keyMarker, uploadIDMarker, delimiter string, maxUploads int) (result ListMultipartsInfo, err error) {
	bktlk := s.newBucketNSLock(bucket)
	bktlkCtx, err := bktlk.GetRLock(ctx, globalOperationTimeout)
	if err != nil {
		return result, err
	}
	ctx = bktlkCtx.Context()
	defer bktlk.RUnlock(bktlkCtx.Cancel)

	if !s.hasBucket(ctx, bucket) {
		return result, BucketNotFound{Bucket: bucket}
	}

	result.MaxUploads = maxUploads
	result.KeyMarker = keyMarker
	result.UploadIDMarker = uploadIDMarker
	result.Prefix = prefix
	result.Delimiter = delimiter

	if maxUploads == 0 {
		return result, nil
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	seekKey := ""
	if keyMarker != "" {
		seekKey = fmt.Sprintf(allUploadSeekKeyFormat, bucket, keyMarker, uploadIDMarker)
	}
	all, err := s.Db.ReadAllChan(ctx, fmt.Sprintf(allUploadPrefixFormat, bucket, prefix), seekKey)
	if err != nil {
		return result, err
	}
	index := 0
	for entry := range all {
		if index == maxUploads {
			result.IsTruncated = true
			break
		}
		var mi MultipartInfo
		if err = entry.UnmarshalValue(&mi); err != nil {
			return result, err
		}
		index++
		result.Uploads = append(result.Uploads, mi)
	}
	if result.IsTruncated {
		next := result.Uploads[len(result.Uploads)-1]
		result.NextKeyMarker = next.Object
		result.NextUploadIDMarker = next.UploadID
	}

	return result, nil
}

func (s *storageSys) markObjetToDelete(c cid.Cid) error {
	return s.Db.Put(newDelObjectKey(), c.String())
}

// GetAllBucket - get all bucket
func (s *storageSys) GetAllBucket(ctx context.Context) []string {
	var m []string
	all, err := s.Db.ReadAllChan(ctx, bucketPrefix, "")
	if err != nil {
		return nil
	}
	for entry := range all {
		data := BucketMetadata{}
		if err = entry.UnmarshalValue(&data); err != nil {
			continue
		}
		m = append(m, data.Name)
	}
	return m
}
func (s *storageSys) deleteObjets(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, s.gcTimeout)
	defer cancel()

	all, err := s.Db.ReadAllChan(ctx, allDeletePrefixFormat, "")
	if err != nil {
		return err
	}
	for entry := range all {
		var root string
		if err = entry.UnmarshalValue(&root); err != nil {
			return err
		}
		c, err := cid.Decode(root)
		if err != nil {
			log.Warnw("decode cid error", "cid", root)
			if err = s.Db.Delete(entry.GetKey()); err != nil {
				return err
			}
			continue
		}
		if err = dagpoolcli.RemoveDAG(ctx, s.DagPool, c); err != nil {
			log.Errorw("remove DAG error", "cid", c.String(), "error", err)
			break
		}
		if err = s.Db.Delete(entry.GetKey()); err != nil {
			return err
		}
	}
	return nil
}

// processObjectGC is a goroutine to do object GC
func (s *storageSys) processObjectGC(ctx context.Context) {
	timer := time.NewTimer(s.gcPeriod)
	defer timer.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			if checkSystemIdle() {
				log.Debug("starting object GC...")
				if err := s.deleteObjets(ctx); err != nil {
					log.Errorf("object GC err: %v", err)
				}
				log.Debug("object GC completed")
			}
			timer.Reset(s.gcPeriod)
		}
	}
}

func getCpuPercent() (float64, error) {
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return 0, err
	}
	return percent[0], nil
}

func checkSystemIdle() bool {
	if p, err := getCpuPercent(); err != nil {
		log.Errorf("get cpu percent error: %v", err)
	} else if p > maxCpuPercent {
		return false
	}
	if v, err := mem.VirtualMemory(); err != nil {
		log.Errorf("get memory used percent error: %v", err)
	} else if v.UsedPercent > maxUsedMemoryPercent {
		return false
	}
	return true
}
