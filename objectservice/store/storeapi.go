package store

import (
	"context"
	"github.com/filedag-project/filedag-storage/objectservice/datatypes"
	"github.com/filedag-project/filedag-storage/objectservice/lock"
	"github.com/filedag-project/filedag-storage/objectservice/pkg/policy"
	"io"
)

var _ ObjectStoreSystemAPI = &storageSys{}

// ObjectStoreSystemAPI object store system API
type ObjectStoreSystemAPI interface {
	StoreStats(ctx context.Context, bucketMetadataMap map[string]BucketMetadata) (DataUsageInfo, error)
	SetNewBucketNSLock(newBucketNSLock func(bucket string) lock.RWLocker)
	SetHasBucket(hasBucket func(ctx context.Context, bucket string) bool)
	StoreObject(ctx context.Context, bucket string, object string, reader io.ReadCloser, size int64, meta map[string]string, fileFolder bool) (ObjectInfo, error)
	GetObject(ctx context.Context, bucket string, object string) (ObjectInfo, io.ReadCloser, error)
	GetObjectInfo(ctx context.Context, bucket string, object string) (meta ObjectInfo, err error)
	DeleteObject(ctx context.Context, bucket string, object string) error
	CleanObjectsInBucket(ctx context.Context, bucket string) error
	GetAllObjectsInBucketInfo(ctx context.Context, bucket string) (bi BucketInfo, err error)
	ListObjects(ctx context.Context, bucket string, prefix string, marker string, delimiter string, maxKeys int) (loi ListObjectsInfo, err error)
	EmptyBucket(ctx context.Context, bucket string) (bool, error)
	ListObjectsV2(ctx context.Context, bucket string, prefix string, continuationToken string, delimiter string, maxKeys int, owner bool, startAfter string) (ListObjectsV2Info, error)
	NewMultipartUpload(ctx context.Context, bucket string, object string, meta map[string]string) (MultipartInfo, error)
	GetMultipartInfo(ctx context.Context, bucket string, object string, uploadID string) (MultipartInfo, error)
	PutObjectPart(ctx context.Context, bucket string, object string, uploadID string, partID int, reader io.ReadCloser, size int64, meta map[string]string) (pi objectPartInfo, err error)
	CompleteMultiPartUpload(ctx context.Context, bucket string, object string, uploadID string, parts []datatypes.CompletePart) (oi ObjectInfo, err error)
	AbortMultipartUpload(ctx context.Context, bucket string, object string, uploadID string) error
	ListObjectParts(ctx context.Context, bucket string, object string, uploadID string, partNumberMarker int, maxParts int) (result ListPartsInfo, err error)
	ListMultipartUploads(ctx context.Context, bucket string, prefix string, keyMarker string, uploadIDMarker string, delimiter string, maxUploads int) (result ListMultipartsInfo, err error)
}

var _ BucketMetadataSysAPI = &bucketMetadataSys{}

//BucketMetadataSysAPI BucketMetadata Sys
type BucketMetadataSysAPI interface {
	NewNSLock(bucket string) lock.RWLocker
	SetEmptyBucket(emptyBucket func(ctx context.Context, bucket string) (bool, error))
	CreateBucket(ctx context.Context, bucket string, region string, accessKey string) error
	GetBucketMeta(ctx context.Context, bucket string) (meta BucketMetadata, err error)
	HasBucket(ctx context.Context, bucket string) bool
	DeleteBucket(ctx context.Context, bucket string, accessKey string) error
	GetAllBucketsOfUser(ctx context.Context, username string) ([]BucketMetadata, error)
	GetAllBucketInfo(ctx context.Context) (allBucketInfo, error)
	UpdateBucketPolicy(ctx context.Context, bucket string, p *policy.Policy) error
	DeleteBucketPolicy(ctx context.Context, bucket string) error
	GetPolicyConfig(ctx context.Context, bucket string) (*policy.Policy, error)
	UpdateBucketTagging(ctx context.Context, bucket string, tags *Tags) error
	DeleteBucketTagging(ctx context.Context, bucket string) error
	GetTaggingConfig(ctx context.Context, bucket string) (*Tags, error)
}
