package store

import (
	"context"
	"encoding/xml"
	"github.com/filedag-project/filedag-storage/objectservice/lock"
	"github.com/filedag-project/filedag-storage/objectservice/objmetadb"
	"github.com/filedag-project/filedag-storage/objectservice/pkg/policy"
	"github.com/syndtr/goleveldb/leveldb"
	"time"
)

const (
	bucketPrefix = "bkt/"
)

// BucketPolicyNotFound - no bucket policy found.
type BucketPolicyNotFound struct {
	Bucket string
	Err    error
}

func (e BucketPolicyNotFound) Error() string {
	return "No bucket policy configuration found for bucket: " + e.Bucket
}

// BucketNotFound - no bucket found.
type BucketNotFound struct {
	Bucket string
	Err    error
}

func (e BucketNotFound) Error() string {
	return "Not found for bucket: " + e.Bucket
}

type BucketTaggingNotFound struct {
	Bucket string
	Err    error
}

func (e BucketTaggingNotFound) Error() string {
	return "No bucket tagging configuration found for bucket: " + e.Bucket
}

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

// Tags is list of tags of XML request/response as per
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetBucketTagging.html#API_GetBucketTagging_RequestBody
type Tags tagging
type tagging struct {
	XMLName xml.Name `xml:"Tagging"`
	TagSet  *TagSet  `xml:"TagSet"`
}

// TagSet represents list of unique tags.
type TagSet struct {
	TagMap   map[string]string
	IsObject bool
}

// BucketMetadata contains bucket metadata.
type BucketMetadata struct {
	Name    string
	Region  string
	Owner   string
	Created time.Time

	PolicyConfig  *policy.Policy
	TaggingConfig *Tags
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

	return sys.setBucketMeta(bucket, NewBucketMetadata(bucket, region, accessKey))
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
func (sys *bucketMetadataSys) DeleteBucket(ctx context.Context, bucket string) error {
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

	return sys.bucketMetaStore.Delete(bucketPrefix + bucket)
}

// GetAllBucketsOfUser metadata for all bucket.
func (sys *bucketMetadataSys) GetAllBucketsOfUser(ctx context.Context, username string) ([]BucketMetadata, error) {
	var m []BucketMetadata
	all, err := sys.bucketMetaStore.ReadAllChan(ctx, bucketPrefix, "")
	if err != nil {
		return nil, err
	}
	for entry := range all {
		data := BucketMetadata{}
		if err = entry.UnmarshalValue(&data); err != nil {
			continue
		}
		if data.Owner != username {
			continue
		}
		m = append(m, data)
	}
	return m, nil
}
