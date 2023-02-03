package store

import (
	"context"
	"github.com/filedag-project/filedag-storage/objectservice/lock"
	"github.com/filedag-project/filedag-storage/objectservice/objmetadb"
	"github.com/filedag-project/filedag-storage/objectservice/pkg/policy"
	"github.com/syndtr/goleveldb/leveldb"
	"time"
)

const (
	bucketPrefix         = "bkt/"
	userBucketInfoPrefix = "userbktinfo/"
	bucketInfoPrefix     = "allbktinfo/"
)

// BucketMetadataSys captures all bucket metadata for a given cluster.
type bucketMetadataSys struct {
	bucketMetaStore objmetadb.ObjStoreMetaDBAPI
	nsLock          *lock.NsLockMap
	emptyBucket     func(ctx context.Context, bucket string) (bool, error)
}

// NewBucketMetadataSys - creates new policy system.
func NewBucketMetadataSys(db objmetadb.ObjStoreMetaDBAPI) *bucketMetadataSys {
	return &bucketMetadataSys{
		bucketMetaStore: db,
		nsLock:          lock.NewNSLock(),
	}
}

// NewBucketMetadata creates BucketMetadata with the supplied name and Created to Now.
func NewBucketMetadata(name, region, accessKey string) *BucketMetadata {
	p := policy.CreateUserBucketPolicy(name, accessKey)
	return &BucketMetadata{
		Name:         name,
		Region:       region,
		Owner:        accessKey,
		Created:      time.Now().UTC(),
		PolicyConfig: p,
	}
}

// NewNSLock - initialize a new namespace RWLocker instance.
func (sys *bucketMetadataSys) NewNSLock(bucket string) lock.RWLocker {
	return sys.nsLock.NewNSLock("meta", bucket)
}

func (sys *bucketMetadataSys) SetEmptyBucket(emptyBucket func(ctx context.Context, bucket string) (bool, error)) {
	sys.emptyBucket = emptyBucket
}

// setBucketMeta - sets a new metadata in-db
func (sys *bucketMetadataSys) setBucketMeta(bucket string, meta *BucketMetadata) error {
	return sys.bucketMetaStore.Put(bucketPrefix+bucket, meta)
}

// CreateBucket - create a new Bucket
func (sys *bucketMetadataSys) CreateBucket(ctx context.Context, bucket, region, accessKey string) error {
	lk := sys.NewNSLock(bucket)
	lkctx, err := lk.GetLock(ctx, globalOperationTimeout)
	if err != nil {
		return err
	}
	ctx = lkctx.Context()
	defer lk.Unlock(lkctx.Cancel)
	meta := NewBucketMetadata(bucket, region, accessKey)
	err = sys.recordUserBucketInfo(ctx, bucket, accessKey, *meta)
	if err != nil {
		return err
	}
	err = sys.setBucketMeta(bucket, meta)
	if err != nil {
		sys.delUserBucketInfo(ctx, bucket, accessKey)
	}
	return err
}

func (sys *bucketMetadataSys) getBucketMeta(bucket string) (meta BucketMetadata, err error) {
	err = sys.bucketMetaStore.Get(bucketPrefix+bucket, &meta)
	if err == leveldb.ErrNotFound {
		err = BucketNotFound{Bucket: bucket, Err: err}
	}
	return meta, err
}

// GetBucketMeta metadata for a bucket.
func (sys *bucketMetadataSys) GetBucketMeta(ctx context.Context, bucket string) (meta BucketMetadata, err error) {
	lk := sys.NewNSLock(bucket)
	lkctx, err := lk.GetRLock(ctx, globalOperationTimeout)
	if err != nil {
		return BucketMetadata{}, err
	}
	ctx = lkctx.Context()
	defer lk.RUnlock(lkctx.Cancel)

	return sys.getBucketMeta(bucket)
}

// HasBucket  metadata for a bucket.
func (sys *bucketMetadataSys) HasBucket(ctx context.Context, bucket string) bool {
	_, err := sys.GetBucketMeta(ctx, bucket)
	return err == nil
}

// DeleteBucket bucket.
func (sys *bucketMetadataSys) DeleteBucket(ctx context.Context, bucket string, accessKey string) error {
	lk := sys.NewNSLock(bucket)
	lkctx, err := lk.GetLock(ctx, deleteOperationTimeout)
	if err != nil {
		return err
	}
	ctx = lkctx.Context()
	defer lk.Unlock(lkctx.Cancel)

	if _, err = sys.getBucketMeta(bucket); err != nil {
		return err
	}

	if empty, err := sys.emptyBucket(ctx, bucket); err != nil {
		return err
	} else if !empty {
		return ErrBucketNotEmpty
	}
	// todo deal del fail
	err = sys.delUserBucketInfo(ctx, bucket, accessKey)
	if err != nil {
		return err
	}
	return sys.bucketMetaStore.Delete(bucketPrefix + bucket)
}
