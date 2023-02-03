package store

import (
	"context"
	"github.com/syndtr/goleveldb/leveldb"
)

type userBucketInfo struct {
	AccessKey string
	Bucket    map[string]BucketMetadata
	Count     uint64
}
type allBucketInfo struct {
	Bucket map[string]BucketMetadata
	Count  uint64
}

// GetAllBucketsOfUser metadata for all bucket.
func (sys *bucketMetadataSys) GetAllBucketsOfUser(ctx context.Context, username string) ([]BucketMetadata, error) {
	var m []BucketMetadata
	info, err := sys.getUserBucketInfo(username)
	if err != nil {
		return nil, err
	}
	for _, meta := range info.Bucket {
		m = append(m, meta)
	}
	return m, nil
}

// GetAllBucketInfo - GetAllBucketInfo in-db
func (sys *bucketMetadataSys) GetAllBucketInfo(ctx context.Context) (allBucketInfo, error) {
	var allUbi allBucketInfo
	err := sys.bucketMetaStore.Get(bucketInfoPrefix, &allUbi)
	if err != nil {
		if err != leveldb.ErrNotFound {
			log.Debugf("bucketMetaStore Get err:%v", err)
			return allUbi, err
		}
	}
	if allUbi.Bucket == nil {
		allUbi.Bucket = make(map[string]BucketMetadata)
	}
	return allUbi, nil
}

// recordUserBucketInfo - recordUserBucketInfo in-db
func (sys *bucketMetadataSys) recordUserBucketInfo(ctx context.Context, bucket, accessKey string, meta BucketMetadata) error {
	ubi, err := sys.getUserBucketInfo(accessKey)
	if err != nil {
		log.Debugf("getUserBucketInfo err:%v", err)
		return err
	}
	ubi.Bucket[bucket] = meta
	ubi.Count++
	ubi.AccessKey = accessKey
	allUbi, err := sys.GetAllBucketInfo(ctx)
	if err != nil {
		log.Debugf("getAllBucketInfo err:%v", err)
		return err
	}
	allUbi.Bucket[bucket] = meta
	allUbi.Count++
	err = sys.bucketMetaStore.Put(bucketInfoPrefix, allUbi)
	if err != nil {
		log.Debugf("bucketMetaStore Put err:%v", err)
		return err
	}
	err = sys.bucketMetaStore.Put(userBucketInfoPrefix+accessKey, ubi)
	if err != nil {
		log.Debugf("bucketMetaStore Put err:%v", err)
		return err
	}
	return nil
}

// getUserBucketInfo - getUserBucketInfo in-db
func (sys *bucketMetadataSys) getUserBucketInfo(accessKey string) (userBucketInfo, error) {
	var ubi userBucketInfo
	err := sys.bucketMetaStore.Get(userBucketInfoPrefix+accessKey, &ubi)
	if err != nil {
		if err != leveldb.ErrNotFound {
			log.Debugf("bucketMetaStore Get err:%v", err)
			return ubi, err
		}
	}
	if ubi.Bucket == nil {
		ubi.Bucket = make(map[string]BucketMetadata)
	}
	return ubi, nil
}

// recordUserBucketInfo - recordUserBucketInfo in-db
func (sys *bucketMetadataSys) delUserBucketInfo(ctx context.Context, bucket, accessKey string) error {
	ubi, err := sys.getUserBucketInfo(accessKey)
	if err != nil {
		return err
	}
	delete(ubi.Bucket, bucket)
	ubi.Count--
	aubi, err := sys.GetAllBucketInfo(ctx)
	if err != nil {
		return err
	}
	delete(aubi.Bucket, bucket)
	aubi.Count--
	err = sys.bucketMetaStore.Put(bucketInfoPrefix, aubi)
	if err != nil {
		return err
	}
	return sys.bucketMetaStore.Put(userBucketInfoPrefix+accessKey, ubi)
}
