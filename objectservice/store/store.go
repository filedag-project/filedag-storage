package store

import (
	"context"
	"errors"
	"fmt"
	dagpoolcli "github.com/filedag-project/filedag-storage/dag/pool/client"
	"github.com/filedag-project/filedag-storage/objectservice/consts"
	"github.com/filedag-project/filedag-storage/objectservice/uleveldb"
	"github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	logging "github.com/ipfs/go-log/v2"
	"github.com/ipfs/go-merkledag"
	ufsio "github.com/ipfs/go-unixfs/io"
	"github.com/syndtr/goleveldb/leveldb"
	"golang.org/x/xerrors"
	"io"
	"net/http"
	"strings"
	"time"
)

var log = logging.Logger("store")

//StorageSys store sys
type StorageSys struct {
	Db         *uleveldb.ULevelDB
	DagPool    ipld.DAGService
	CidBuilder cid.Builder
}

const objectPrefixFormat = "obj-%s-%s-%s/"
const allObjectPrefixFormat = "obj-%s-%s-%s"
const allObjectSeekPrefixFormat = "obj-%s-%s-%s"

var ErrObjectNotFound = errors.New("object not found")

func getObjectKey(user, bucket, object string) string {
	return fmt.Sprintf(objectPrefixFormat, user, bucket, object)
}

//StoreObject store object
func (s *StorageSys) StoreObject(ctx context.Context, user, bucket, object string, reader io.ReadCloser, size int64, meta map[string]string) (ObjectInfo, error) {
	if strings.HasPrefix(object, "/") {
		object = object[1:]
	}
	node, err := dagpoolcli.BalanceNode(reader, s.DagPool, s.CidBuilder)
	if err != nil {
		return ObjectInfo{}, err
	}
	objInfo := ObjectInfo{
		Bucket:           bucket,
		Name:             object,
		ModTime:          time.Now().UTC(),
		Size:             size,
		IsDir:            false,
		ETag:             node.Cid().String(),
		VersionID:        "",
		IsLatest:         true,
		DeleteMarker:     false,
		ContentType:      meta[strings.ToLower(consts.ContentType)],
		ContentEncoding:  meta[strings.ToLower(consts.ContentEncoding)],
		Parts:            nil,
		SuccessorModTime: time.Now().UTC(),
	}
	// Update expires
	if exp, ok := meta[strings.ToLower(consts.Expires)]; ok {
		if t, e := time.Parse(http.TimeFormat, exp); e == nil {
			objInfo.Expires = t.UTC()
		}
	}
	// Has old file?
	if oldObjInfo, exist := s.HasObject(ctx, user, bucket, object); exist {
		c, err := cid.Decode(oldObjInfo.ETag)
		if err != nil {
			log.Warnw("decode cid error", "cid", oldObjInfo.ETag)
		} else if err = dagpoolcli.RemoveDAG(ctx, s.DagPool, c); err != nil {
			log.Errorw("remove DAG error", "cid", oldObjInfo.ETag)
		}
	}

	err = s.Db.Put(getObjectKey(user, bucket, object), objInfo)
	if err != nil {
		return ObjectInfo{}, err
	}
	return objInfo, nil
}

//GetObject Get object
func (s *StorageSys) GetObject(ctx context.Context, user, bucket, object string) (ObjectInfo, ufsio.DagReader, error) {
	meta := ObjectInfo{}
	if strings.HasPrefix(object, "/") {
		object = object[1:]
	}
	err := s.Db.Get(getObjectKey(user, bucket, object), &meta)
	if err != nil {
		if xerrors.Is(err, leveldb.ErrNotFound) {
			return ObjectInfo{}, nil, ErrObjectNotFound
		}
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

// HasObject has Object ?
func (s *StorageSys) HasObject(ctx context.Context, user, bucket, object string) (ObjectInfo, bool) {
	meta := ObjectInfo{}
	if strings.HasPrefix(object, "/") {
		object = object[1:]
	}
	err := s.Db.Get(getObjectKey(user, bucket, object), &meta)
	if err != nil {
		return ObjectInfo{}, false
	}
	return meta, true
}

func (s *StorageSys) GetObjectInfo(ctx context.Context, user, bucket, object string) (meta ObjectInfo, err error) {
	if strings.HasPrefix(object, "/") {
		object = object[1:]
	}
	err = s.Db.Get(getObjectKey(user, bucket, object), &meta)
	return
}

//DeleteObject delete object
func (s *StorageSys) DeleteObject(ctx context.Context, user, bucket, object string) error {
	if strings.HasPrefix(object, "/") {
		object = object[1:]
	}
	meta := ObjectInfo{}
	err := s.Db.Get(getObjectKey(user, bucket, object), &meta)
	if err != nil {
		if xerrors.Is(err, leveldb.ErrNotFound) {
			return ErrObjectNotFound
		}
		return err
	}
	cid, err := cid.Decode(meta.ETag)
	if err != nil {
		return err
	}

	if err = s.Db.Delete(getObjectKey(user, bucket, object)); err != nil {
		return err
	}

	if err = dagpoolcli.RemoveDAG(ctx, s.DagPool, cid); err != nil {
		return err
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
func (s *StorageSys) ListObjects(ctx context.Context, user, bucket string, prefix string, marker string, delimiter string, maxKeys int) (loi ListObjectsInfo, err error) {
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
		objInfo, err := s.GetObjectInfo(ctx, user, bucket, prefix)
		if err == nil {
			loi.Objects = append(loi.Objects, objInfo)
			return loi, nil
		}
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	seekKey := ""
	if marker != "" {
		seekKey = fmt.Sprintf(allObjectSeekPrefixFormat, user, bucket, marker)
	}
	all, err := s.Db.ReadAllChan(ctx, fmt.Sprintf(allObjectPrefixFormat, user, bucket, prefix), seekKey)
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
		index++
		loi.Objects = append(loi.Objects, o)
	}
	if loi.IsTruncated {
		loi.NextMarker = loi.Objects[len(loi.Objects)-1].Name
	}

	return loi, nil
}

//MkBucket store object
func (s *StorageSys) MkBucket(parentDirectoryPath string, bucket string) error {
	return nil
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
func (s *StorageSys) ListObjectsV2(ctx context.Context, user, bucket string, prefix string, continuationToken string, delimiter string, maxKeys int, owner bool, startAfter string) (ListObjectsV2Info, error) {
	marker := continuationToken
	if marker == "" {
		marker = startAfter
	}
	loi, err := s.ListObjects(ctx, user, bucket, prefix, marker, delimiter, maxKeys)
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

//NewStorageSys new a storage sys
func NewStorageSys(dagService ipld.DAGService, db *uleveldb.ULevelDB) *StorageSys {
	cidBuilder, _ := merkledag.PrefixForCidVersion(0)
	return &StorageSys{
		Db:         db,
		DagPool:    dagService,
		CidBuilder: cidBuilder,
	}
}
