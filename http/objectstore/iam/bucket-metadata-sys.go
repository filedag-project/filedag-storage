package iam

import (
	"errors"
	"github.com/filedag-project/filedag-storage/http/objectstore/api_errors"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	"time"
)

const (
	bucketPrefix = "buckets/"
)

// bucketPolicyNotFound - no bucket policy found.
type bucketPolicyNotFound api_errors.GenericBucketError

func (e bucketPolicyNotFound) Error() string {
	return "No bucket policy configuration found for bucket: " + e.Bucket
}

// bucketMetadataSys captures all bucket metadata for a given cluster.
type bucketMetadataSys struct {
	db *uleveldb.ULevelDB
}

// newBucketMetadataSys - creates new policy system.
func newBucketMetadataSys() *bucketMetadataSys {
	return &bucketMetadataSys{
		db: uleveldb.DBClient,
	}
}

// bucketMetadata contains bucket metadata.
type bucketMetadata struct {
	Name         string
	Created      time.Time
	PolicyConfig *Policy
}

// newBucketMetadata creates bucketMetadata with the supplied name and Created to Now.
func newBucketMetadata(name string) bucketMetadata {
	return bucketMetadata{
		Name:    name,
		Created: time.Now().UTC(),
	}
}

// GetPolicyConfig returns configured bucket policy
func (sys *bucketMetadataSys) GetPolicyConfig(bucket, accessKey string) (*Policy, error) {
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
func (sys *bucketMetadataSys) GetConfig(bucket, accessKey string) (bucketMetadata, error) {
	var meta bucketMetadata
	err := sys.Get(bucket, accessKey, &meta)
	if err != nil {
		return bucketMetadata{}, err
	}
	return meta, nil
}

// Set - sets a new metadata in-db
func (sys *bucketMetadataSys) Set(bucket, username string, meta bucketMetadata) error {
	err := sys.db.Put(bucketPrefix+username+"-"+bucket, meta)
	if err != nil {
		return err
	}
	return nil
}

// Get metadata for a bucket.
func (sys *bucketMetadataSys) Get(bucket, username string, meta *bucketMetadata) error {
	err := sys.db.Get(bucketPrefix+username+"-"+bucket, meta)
	if err != nil {
		return err
	}
	return nil
}
