package policy

import (
	"errors"
	"github.com/filedag-project/filedag-storage/http/objectstore/api_errors"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	"time"
)

var globalBucketMetadataSys = NewBucketMetadataSys()

const (
	bucketPrefix = "buckets/"
)

// bucketPolicyNotFound - no bucket policy found.
type bucketPolicyNotFound api_errors.GenericBucketError

func (e bucketPolicyNotFound) Error() string {
	return "No bucket policy configuration found for bucket: " + e.Bucket
}

// BucketMetadataSys captures all bucket metadata for a given cluster.
type BucketMetadataSys struct {
	db *uleveldb.ULeveldb
}

// NewBucketMetadataSys - creates new policy system.
func NewBucketMetadataSys() *BucketMetadataSys {
	return &BucketMetadataSys{
		db: uleveldb.NewLevelDB(),
	}
}

// BucketMetadata contains bucket metadata.
type BucketMetadata struct {
	Name         string
	Created      time.Time
	PolicyConfig *Policy
}

// newBucketMetadata creates BucketMetadata with the supplied name and Created to Now.
func newBucketMetadata(name string) BucketMetadata {
	return BucketMetadata{
		Name:    name,
		Created: time.Now().UTC(),
	}
}

// GetPolicyConfig returns configured bucket policy
func (sys *BucketMetadataSys) GetPolicyConfig(bucket, accessKey string) (*Policy, error) {
	meta, err := sys.GetConfig(bucket, accessKey)
	if err != nil {
		if errors.Is(err, api_errors.ErrConfigNotFound) {
			return nil, bucketPolicyNotFound{Bucket: bucket}
		}
		return nil, err
	}
	if meta.PolicyConfig == nil {
		return nil, bucketPolicyNotFound{Bucket: bucket}
	}
	return meta.PolicyConfig, nil
}

// GetConfig returns a specific configuration from the bucket metadata.
func (sys *BucketMetadataSys) GetConfig(bucket, accessKey string) (BucketMetadata, error) {
	var meta BucketMetadata
	err := sys.Get(bucket, accessKey, &meta)
	if err != nil {
		return BucketMetadata{}, err
	}
	return meta, nil
}

// Set - sets a new metadata in-db
func (sys *BucketMetadataSys) Set(bucket, username string, meta BucketMetadata) error {
	err := sys.db.Put(bucketPrefix+username+"-"+bucket, meta)
	if err != nil {
		return err
	}
	return nil
}

// Get metadata for a bucket.
func (sys *BucketMetadataSys) Get(bucket, username string, meta *BucketMetadata) error {
	err := sys.db.Get(bucketPrefix+username+"-"+bucket, meta)
	if err != nil {
		return err
	}
	return nil
}
