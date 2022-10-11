package store

import (
	"context"
	"errors"
	"fmt"
	"github.com/dustin/go-humanize"
	dagpoolcli "github.com/filedag-project/filedag-storage/dag/pool/client"
	"github.com/filedag-project/filedag-storage/objectservice/consts"
	"github.com/filedag-project/filedag-storage/objectservice/datatypes"
	"github.com/filedag-project/filedag-storage/objectservice/lock"
	"github.com/filedag-project/filedag-storage/objectservice/uleveldb"
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

const (
	// bigFileThreshold is the point where we add readahead to put operations.
	bigFileThreshold = 64 * humanize.MiByte
	// equals unixfsChunkSize
	chunkSize int = 1 << 20

	objectKeyFormat        = "obj-%s-%s/"
	allObjectPrefixFormat  = "obj-%s-%s"
	allObjectSeekKeyFormat = "obj-%s-%s/"

	uploadKeyFormat        = "uploadObj-%s-%s-%s"
	allUploadPrefixFormat  = "uploadObj-%s-%s"
	allUploadSeekKeyFormat = "uploadObj-%s-%s-%s"

	deleteKeyFormat       = "delObj-%s"
	allDeletePrefixFormat = "delObj-"

	globalOperationTimeout = 5 * time.Minute
	deleteOperationTimeout = 1 * time.Minute

	maxCpuPercent        = 60
	maxUsedMemoryPercent = 80
)

var ErrObjectNotFound = errors.New("object not found")
var ErrBucketNotEmpty = errors.New("bucket not empty")

//StorageSys store sys
type StorageSys struct {
	Db              *uleveldb.ULevelDB
	DagPool         ipld.DAGService
	CidBuilder      cid.Builder
	nsLock          *lock.NsLockMap
	newBucketNSLock func(bucket string) lock.RWLocker
	hasBucket       func(ctx context.Context, bucket string) bool

	gcPeriod  time.Duration
	gcTimeout time.Duration
}

//NewStorageSys new a storage sys
func NewStorageSys(ctx context.Context, dagService ipld.DAGService, db *uleveldb.ULevelDB) *StorageSys {
	cidBuilder, _ := merkledag.PrefixForCidVersion(0)
	s := &StorageSys{
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

// NewNSLock - initialize a new namespace RWLocker instance.
func (s *StorageSys) NewNSLock(bucket string, objects ...string) lock.RWLocker {
	return s.nsLock.NewNSLock(bucket, objects...)
}

func (s *StorageSys) SetNewBucketNSLock(newBucketNSLock func(bucket string) lock.RWLocker) {
	s.newBucketNSLock = newBucketNSLock
}

func (s *StorageSys) SetHasBucket(hasBucket func(ctx context.Context, bucket string) bool) {
	s.hasBucket = hasBucket
}

func (s *StorageSys) store(ctx context.Context, reader io.ReadCloser, size int64) (cid.Cid, error) {
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

func (s *StorageSys) checkAndDeleteObjectData(ctx context.Context, bucket, object string) {
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
func (s *StorageSys) StoreObject(ctx context.Context, bucket, object string, reader io.ReadCloser, size int64, meta map[string]string) (ObjectInfo, error) {
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

	root, err := s.store(ctx, reader, size)
	if err != nil {
		return ObjectInfo{}, err
	}

	objInfo := ObjectInfo{
		Bucket:           bucket,
		Name:             object,
		ModTime:          time.Now().UTC(),
		Size:             size,
		IsDir:            false,
		ETag:             root.String(),
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

	lk := s.NewNSLock(bucket, object)
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
func (s *StorageSys) GetObject(ctx context.Context, bucket, object string) (ObjectInfo, io.ReadCloser, error) {
	lk := s.NewNSLock(bucket, object)
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

func (s *StorageSys) getObjectInfo(ctx context.Context, bucket, object string) (meta ObjectInfo, err error) {
	err = s.Db.Get(getObjectKey(bucket, object), &meta)
	if err != nil {
		if xerrors.Is(err, leveldb.ErrNotFound) {
			return meta, ErrObjectNotFound
		}
	}
	return
}

func (s *StorageSys) GetObjectInfo(ctx context.Context, bucket, object string) (meta ObjectInfo, err error) {
	lk := s.NewNSLock(bucket, object)
	lkctx, err := lk.GetRLock(ctx, globalOperationTimeout)
	if err != nil {
		return ObjectInfo{}, err
	}
	ctx = lkctx.Context()
	defer lk.RUnlock(lkctx.Cancel)

	return s.getObjectInfo(ctx, bucket, object)
}

//DeleteObject delete object
func (s *StorageSys) DeleteObject(ctx context.Context, bucket, object string) error {
	lk := s.NewNSLock(bucket, object)
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
	cid, err := cid.Decode(meta.ETag)
	if err != nil {
		return err
	}

	if err = s.Db.Delete(getObjectKey(bucket, object)); err != nil {
		return err
	}

	if err = s.markObjetToDelete(cid); err != nil {
		log.Errorw("mark Objet to delete error", "bucket", bucket, "object", object, "cid", meta.ETag, "error", err)
	}
	return nil
}

// ListObjectsInfo - container for list objects.
type ListObjectsInfo struct {
	// Indicates whether the returned list objects response is truncated. A
	// value of true indicates that the list was truncated. The list can be truncated
	// if the number of objects exceeds the limit allowed or specified
	// by max keys.
	IsTruncated bool

	// When response is truncated (the IsTruncated element value in the response is true),
	// you can use the key name in this field as marker in the subsequent
	// request to get next set of objects.
	//
	// NOTE: AWS S3 returns NextMarker only if you have delimiter request parameter specified,
	NextMarker string

	// List of objects info for this request.
	Objects []ObjectInfo

	// List of prefixes for this request.
	Prefixes []string
}

//ListObjects list user object
//TODO use more params
func (s *StorageSys) ListObjects(ctx context.Context, bucket string, prefix string, marker string, delimiter string, maxKeys int) (loi ListObjectsInfo, err error) {
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
	log.Infow("ListObjects ReadAllChan", "prefixKey", prefixKey, "seekKey", seekKey)
	all, err := s.Db.ReadAllChan(ctx, prefixKey, seekKey)
	if err != nil {
		return loi, err
	}
	index := 0
	for entry := range all {
		if index == maxKeys {
			loi.IsTruncated = true
			break
		}
		var o ObjectInfo
		if err = entry.UnmarshalValue(&o); err != nil {
			return loi, err
		}
		log.Infow("ListObjects", "index", index, "key", entry.Key, "object name", o.Name)
		index++
		loi.Objects = append(loi.Objects, o)
	}
	if loi.IsTruncated {
		loi.NextMarker = loi.Objects[len(loi.Objects)-1].Name
		log.Infow("ListObjects", "last object name", loi.Objects[len(loi.Objects)-1].Name)
		log.Infow("ListObjects", "first object name", loi.Objects[0].Name)
	}

	return loi, nil
}

func (s *StorageSys) EmptyBucket(ctx context.Context, bucket string) (bool, error) {
	loi, err := s.ListObjects(ctx, bucket, "", "", "", 1)
	if err != nil {
		return false, err
	}
	return len(loi.Objects) == 0, nil
}

// ListObjectsV2Info - container for list objects version 2.
type ListObjectsV2Info struct {
	// Indicates whether the returned list objects response is truncated. A
	// value of true indicates that the list was truncated. The list can be truncated
	// if the number of objects exceeds the limit allowed or specified
	// by max keys.
	IsTruncated bool

	// When response is truncated (the IsTruncated element value in the response
	// is true), you can use the key name in this field as marker in the subsequent
	// request to get next set of objects.
	//
	// NOTE: This element is returned only if you have delimiter request parameter
	// specified.
	ContinuationToken     string
	NextContinuationToken string

	// List of objects info for this request.
	Objects []ObjectInfo

	// List of prefixes for this request.
	Prefixes []string
}

// ListObjectsV2 list objects
func (s *StorageSys) ListObjectsV2(ctx context.Context, bucket string, prefix string, continuationToken string, delimiter string, maxKeys int, owner bool, startAfter string) (ListObjectsV2Info, error) {
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

func (s *StorageSys) NewMultipartUpload(ctx context.Context, bucket string, object string, meta map[string]string) (MultipartInfo, error) {
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

func (s *StorageSys) GetMultipartInfo(ctx context.Context, bucket string, object string, uploadID string) (MultipartInfo, error) {
	bktlk := s.newBucketNSLock(bucket)
	bktlkCtx, err := bktlk.GetRLock(ctx, globalOperationTimeout)
	if err != nil {
		return MultipartInfo{}, err
	}
	ctx = bktlkCtx.Context()
	defer bktlk.RUnlock(bktlkCtx.Cancel)

	uploadIDLock := s.NewNSLock(bucket, lock.PathJoin(object, uploadID))
	lkctx, err := uploadIDLock.GetRLock(ctx, globalOperationTimeout)
	if err != nil {
		return MultipartInfo{}, err
	}
	ctx = lkctx.Context()
	defer uploadIDLock.RUnlock(lkctx.Cancel)

	return s.getMultipartInfo(ctx, bucket, object, uploadID)
}

func (s *StorageSys) getMultipartInfo(ctx context.Context, bucket string, object string, uploadID string) (MultipartInfo, error) {
	info := MultipartInfo{}
	err := s.Db.Get(getUploadKey(bucket, object, uploadID), &info)
	return info, err
}

func (s *StorageSys) PutObjectPart(ctx context.Context, bucket string, object string, uploadID string, partID int, reader io.ReadCloser, size int64, meta map[string]string) (pi objectPartInfo, err error) {
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

	uploadIDLock := s.NewNSLock(bucket, lock.PathJoin(object, uploadID))
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

func (s *StorageSys) removeMultipartInfo(ctx context.Context, bucket string, object string, uploadID string) error {
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

func (s *StorageSys) CompleteMultiPartUpload(ctx context.Context, bucket string, object string, uploadID string, parts []datatypes.CompletePart) (oi ObjectInfo, err error) {
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

	uploadIDLock := s.NewNSLock(bucket, lock.PathJoin(object, uploadID))
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

	lk := s.NewNSLock(bucket, object)
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

func (s *StorageSys) AbortMultipartUpload(ctx context.Context, bucket string, object string, uploadID string) error {
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

	uploadIDLock := s.NewNSLock(bucket, lock.PathJoin(object, uploadID))
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

// ListPartsInfo - represents list of all parts.
type ListPartsInfo struct {
	// Name of the bucket.
	Bucket string

	// Name of the object.
	Object string

	// Upload ID identifying the multipart upload whose parts are being listed.
	UploadID string

	// Part number after which listing begins.
	PartNumberMarker int

	// When a list is truncated, this element specifies the last part in the list,
	// as well as the value to use for the part-number-marker request parameter
	// in a subsequent request.
	NextPartNumberMarker int

	// Maximum number of parts that were allowed in the response.
	MaxParts int

	// Indicates whether the returned list of parts is truncated.
	IsTruncated bool

	// List of all parts.
	Parts []objectPartInfo

	// Any metadata set during InitMultipartUpload, including encryption headers.
	Metadata map[string]string

	// ChecksumAlgorithm if set
	ChecksumAlgorithm string
}

func (s *StorageSys) ListObjectParts(ctx context.Context, bucket, object, uploadID string, partNumberMarker, maxParts int) (result ListPartsInfo, err error) {
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

	uploadIDLock := s.NewNSLock(bucket, lock.PathJoin(object, uploadID))
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

// ListMultipartsInfo - represnets bucket resources for incomplete multipart uploads.
type ListMultipartsInfo struct {
	// Together with upload-id-marker, this parameter specifies the multipart upload
	// after which listing should begin.
	KeyMarker string

	// Together with key-marker, specifies the multipart upload after which listing
	// should begin. If key-marker is not specified, the upload-id-marker parameter
	// is ignored.
	UploadIDMarker string

	// When a list is truncated, this element specifies the value that should be
	// used for the key-marker request parameter in a subsequent request.
	NextKeyMarker string

	// When a list is truncated, this element specifies the value that should be
	// used for the upload-id-marker request parameter in a subsequent request.
	NextUploadIDMarker string

	// Maximum number of multipart uploads that could have been included in the
	// response.
	MaxUploads int

	// Indicates whether the returned list of multipart uploads is truncated. A
	// value of true indicates that the list was truncated. The list can be truncated
	// if the number of multipart uploads exceeds the limit allowed or specified
	// by max uploads.
	IsTruncated bool

	// List of all pending uploads.
	Uploads []MultipartInfo

	// When a prefix is provided in the request, The result contains only keys
	// starting with the specified prefix.
	Prefix string

	// A character used to truncate the object prefixes.
	// NOTE: only supported delimiter is '/'.
	Delimiter string

	// CommonPrefixes contains all (if there are any) keys between Prefix and the
	// next occurrence of the string specified by delimiter.
	CommonPrefixes []string

	EncodingType string // Not supported yet.
}

func (s *StorageSys) ListMultipartUploads(ctx context.Context, bucket, prefix, keyMarker, uploadIDMarker, delimiter string, maxUploads int) (result ListMultipartsInfo, err error) {
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

func (s *StorageSys) markObjetToDelete(c cid.Cid) error {
	return s.Db.Put(newDelObjectKey(), c.String())
}

func (s *StorageSys) deleteObjets(ctx context.Context) error {
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
			if err = s.Db.Delete(entry.Key); err != nil {
				return err
			}
			continue
		}
		if err = dagpoolcli.RemoveDAG(ctx, s.DagPool, c); err != nil {
			log.Errorw("remove DAG error", "cid", c.String(), "error", err)
			break
		}
		if err = s.Db.Delete(entry.Key); err != nil {
			return err
		}
	}
	return nil
}

// processObjectGC is a goroutine to do object GC
func (s *StorageSys) processObjectGC(ctx context.Context) {
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
