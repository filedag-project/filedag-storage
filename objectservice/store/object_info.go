package store

import (
	"context"
	"github.com/syndtr/goleveldb/leveldb"
	"golang.org/x/xerrors"
)

const objInBktInfoPrefix = "objinbktinfo/"

func (s *storageSys) recordObjectInfo(ctx context.Context, info ObjectInfo) error {
	var oldSize uint64 = 0
	if s.hasObjectInfo(ctx, info.Bucket, info.Name) {
		objectInfo, err := s.getObjectInfo(ctx, info.Bucket, info.Name)
		if err != nil {
			return err
		}
		oldSize = uint64(objectInfo.Size)
	}
	bucketInfo, err := s.GetAllObjectsInBucketInfo(ctx, info.Bucket)
	if err != nil {
		return err
	}
	bucketInfo.Objects++
	bucketInfo.Name = info.Bucket
	bucketInfo.Size = bucketInfo.Size + uint64(info.Size) - oldSize
	err = s.Db.Put(objInBktInfoPrefix, bucketInfo)
	if err != nil {
		return err
	}
	return nil
}
func (s *storageSys) reduceObjectInfo(ctx context.Context, info ObjectInfo) error {
	bucketInfo, err := s.GetAllObjectsInBucketInfo(ctx, info.Bucket)
	if err != nil {
		return err
	}
	bucketInfo.Size = bucketInfo.Size - uint64(info.Size)
	bucketInfo.Objects--
	err = s.Db.Put(objInBktInfoPrefix, bucketInfo)
	if err != nil {
		return err
	}
	return nil
}
func (s *storageSys) hasObjectInfo(ctx context.Context, bucket, obj string) bool {
	_, err := s.getObjectInfo(ctx, bucket, obj)
	if xerrors.Is(err, ErrObjectNotFound) {
		return false
	}
	return true
}

//GetAllObjectsInBucketInfo Get BucketInfo
func (s *storageSys) GetAllObjectsInBucketInfo(ctx context.Context, bucket string) (bi BucketInfo, err error) {
	err = s.Db.Get(objInBktInfoPrefix, &bi)
	if err != nil {
		if err != leveldb.ErrNotFound {
			log.Debugf("bucketMetaStore Get err:%v", err)
			return bi, err
		}
	}
	return bi, nil
}
