package store

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/pool"
	"github.com/filedag-project/filedag-storage/dag/pool/config"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	logging "github.com/ipfs/go-log/v2"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

var log = logging.Logger("store")

//StorageSys store sys
type StorageSys struct {
	Db      *uleveldb.ULevelDB
	DagPool pool.DAGPool
}

const (
	PoolStorePath        = "POOL_STORE_PATH"
	PoolBatchNum         = "POOL_BATCH_NUM"
	PoolCaskNum          = "POOL_CASK_NUM"
	defaultPoolStorePath = "./"
	defaultPoolBatchNum  = 4
	defaultPoolCaskNum   = 2
)
const objectPrefixTemplate = "object-%s-%s-%s/"
const allObjectPrefixTemplate = "object-%s-%s-"

//StoreObject store object
func (s *StorageSys) StoreObject(ctx context.Context, user, bucket, object string, reader io.ReadCloser) (ObjectInfo, error) {
	if strings.HasPrefix(object, "/") {
		object = object[1:]
	}
	cid, err := s.DagPool.Add(ctx, reader)
	if err != nil {
		return ObjectInfo{}, err
	}
	all, err := ioutil.ReadAll(reader)
	meta := ObjectInfo{
		Bucket:           bucket,
		Name:             object,
		ModTime:          time.Now().UTC(),
		Size:             int64(len(all)),
		IsDir:            false,
		ETag:             cid,
		VersionID:        "",
		IsLatest:         false,
		DeleteMarker:     false,
		ContentType:      "application/x-msdownload",
		ContentEncoding:  "",
		Expires:          time.Unix(0, 0).UTC(),
		Parts:            nil,
		AccTime:          time.Unix(0, 0).UTC(),
		SuccessorModTime: time.Now().UTC(),
	}
	err = s.Db.Put(fmt.Sprintf(objectPrefixTemplate, user, bucket, object), meta)
	if err != nil {
		return ObjectInfo{}, err
	}
	return meta, nil
}

//GetObject Get object
func (s *StorageSys) GetObject(ctx context.Context, user, bucket, object string) (ObjectInfo, io.ReadCloser, error) {
	meta := ObjectInfo{}
	if strings.HasPrefix(object, "/") {
		object = object[1:]
	}
	err := s.Db.Get(fmt.Sprintf(objectPrefixTemplate, user, bucket, object), &meta)
	if err != nil {
		return ObjectInfo{}, nil, err
	}
	reader, err := s.DagPool.Get(ctx, meta.ETag)
	return meta, reader, nil
}

// HasObject has Object ?
func (s *StorageSys) HasObject(ctx context.Context, user, bucket, object string) (ObjectInfo, bool) {
	meta := ObjectInfo{}
	if strings.HasPrefix(object, "/") {
		object = object[1:]
	}
	err := s.Db.Get(fmt.Sprintf(objectPrefixTemplate, user, bucket, object), &meta)
	if err != nil {
		return ObjectInfo{}, false
	}
	return meta, true
}

//DeleteObject Get object
func (s *StorageSys) DeleteObject(user, bucket, object string) error {
	//err := s.dagPool.DelFile(bucket, object)
	err := s.Db.Delete(fmt.Sprintf(objectPrefixTemplate, user, bucket, object))
	if err != nil {
		return err
	}
	return nil
}

//ListObject list user object
func (s *StorageSys) ListObject(user, bucket string) ([]ObjectInfo, error) {
	var objs []ObjectInfo
	objMap, err := s.Db.ReadAll(fmt.Sprintf(allObjectPrefixTemplate, user, bucket))
	if err != nil {
		return nil, err
	}
	for _, v := range objMap {
		var o ObjectInfo
		json.Unmarshal([]byte(v), &o)
		objs = append(objs, o)
	}
	return objs, nil
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
//todo use more param
func (s *StorageSys) ListObjectsV2(ctx context.Context, bucket, user string, prefix string, token string, delimiter string, keys int, owner bool, after string) (ListObjectsV2Info, error) {
	objects, err := s.ListObject(user, bucket)
	var o ListObjectsV2Info
	if err != nil {
		return o, err
	}
	count := 0
	for _, v := range objects {
		if v.Name != after {
			continue
		}
		if count > keys {
			break
		}
		count++
		o.ContinuationToken = token
		o.IsTruncated = true
		o.Objects = append(o.Objects, v)
	}
	return o, nil
}

//Init storage sys
func (s *StorageSys) Init() error {
	s.Db = uleveldb.DBClient
	batchNum, err := strconv.Atoi(os.Getenv(PoolBatchNum))
	if err != nil {
		//log.Errorf("get PoolBatchNum err %v,use default",err)
		batchNum = defaultPoolBatchNum
	}
	caskNum, err := strconv.Atoi(os.Getenv(PoolCaskNum))
	if err != nil {
		//log.Errorf("get PoolCaskNum err %v,use default",err)
		caskNum = defaultPoolCaskNum
	}
	var path string
	if os.Getenv(PoolStorePath) == "" {
		//log.Errorf("get PoolStorePath err %v,use default",err)
		path = defaultPoolStorePath
	} else {
		path = os.Getenv(PoolStorePath)
	}
	s.DagPool, err = pool.NewSimplePool(&config.SimplePoolConfig{
		StorePath: path,
		BatchNum:  batchNum,
		CaskNum:   caskNum,
	})
	if err != nil {
		return err
	}
	return nil
}
