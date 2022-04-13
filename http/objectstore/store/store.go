package store

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/pool"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

//StorageSys store sys
type StorageSys struct {
	db      *uleveldb.ULevelDB
	dagPool pool.DAGPool
}

const (
	PoolStorePath = "POOL_STORE_PATH"
	PoolBatchNum  = "POOL_BATCH_NUM"
	PoolCaskNum   = "POOL_CASK_NUM"
)
const objectPrefixTemplate = "object-%s-%s-%s/"
const allObjectPrefixTemplate = "object-%s-%s-"

//StoreObject store object
func (s *StorageSys) StoreObject(ctx context.Context, user, bucket, object string, reader io.ReadCloser) (ObjectInfo, error) {
	if strings.HasPrefix(object, "/") {
		object = object[1:]
	}
	cid, err := s.dagPool.Add(ctx, reader)
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
	err = s.db.Put(fmt.Sprintf(objectPrefixTemplate, user, bucket, object), meta)
	if err != nil {
		return ObjectInfo{}, err
	}
	return meta, nil
}

//GetObject Get object
func (s *StorageSys) GetObject(ctx context.Context, user, bucket, object string) (ObjectInfo, io.ReadCloser, error) {
	meta := ObjectInfo{}
	err := s.db.Get(fmt.Sprintf(objectPrefixTemplate, user, bucket, object), &meta)
	if err != nil {
		return ObjectInfo{}, nil, err
	}
	reader, err := s.dagPool.Get(ctx, meta.ETag)
	return meta, reader, nil
}

// HasObject has Object ?
func (s *StorageSys) HasObject(ctx context.Context, user, bucket, object string) (ObjectInfo, bool) {
	meta := ObjectInfo{}
	err := s.db.Get(fmt.Sprintf(objectPrefixTemplate, user, bucket, object), &meta)
	if err != nil {
		return ObjectInfo{}, false
	}
	return meta, true
}

//DeleteObject Get object
func (s *StorageSys) DeleteObject(user, bucket, object string) error {
	//err := s.dagPool.DelFile(bucket, object)
	err := s.db.Delete(fmt.Sprintf(objectPrefixTemplate, user, bucket, object))
	if err != nil {
		return err
	}
	return nil
}

//ListObject list user object
func (s *StorageSys) ListObject(user, bucket string) ([]ObjectInfo, error) {
	var objs []ObjectInfo
	objMap, err := s.db.ReadAll(fmt.Sprintf(allObjectPrefixTemplate, user, bucket))
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

//Init storage sys
func (s *StorageSys) Init() error {
	s.db = uleveldb.DBClient
	batchNum, err := strconv.Atoi(os.Getenv(PoolBatchNum))
	if err != nil {
		return err
	}
	caskNum, err := strconv.Atoi(os.Getenv(PoolCaskNum))
	if err != nil {
		return err
	}
	s.dagPool, err = pool.NewSimplePool(&pool.SimplePoolConfig{
		StorePath: os.Getenv(PoolStorePath),
		BatchNum:  batchNum,
		CaskNum:   caskNum,
	})
	if err != nil {
		return err
	}
	return nil
}
