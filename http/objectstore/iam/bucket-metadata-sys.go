package iam

import (
	"encoding/json"
	"errors"
	"github.com/filedag-project/filedag-storage/http/objectstore/api_errors"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/policy"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	"strings"
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
	PolicyConfig *policy.Policy
}

// newBucketMetadata creates bucketMetadata with the supplied name and Created to Now.
func newBucketMetadata(name string) bucketMetadata {
	return bucketMetadata{
		Name:    name,
		Created: time.Now().UTC(),
	}
}

// GetPolicyConfig returns configured bucket policy
func (sys *bucketMetadataSys) GetPolicyConfig(bucket, accessKey string) (*policy.Policy, error) {
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

// Head metadata for a bucket.
func (sys *bucketMetadataSys) Head(bucket, username string) bool {
	var meta bucketMetadata
	err := sys.db.Get(bucketPrefix+username+"-"+bucket, &meta)
	if err != nil {
		return false
	}
	return true
}

// Delete bucket.
func (sys *bucketMetadataSys) Delete(username, bucket string) error {
	err := sys.db.Delete(bucketPrefix + username + "-" + bucket)
	if err != nil {
		return err
	}
	return nil
}

// Get metadata for a bucket.
func (sys *bucketMetadataSys) Update(username, bucket string, meta *bucketMetadata) error {
	err := sys.db.Put(bucketPrefix+username+"-"+bucket, meta)
	if err != nil {
		return err
	}
	return nil
}

// getAllBucketOfUser metadata for all bucket.
func (sys *bucketMetadataSys) getAllBucketOfUser(username string) ([]bucketMetadata, error) {
	var m []bucketMetadata
	mb, err := sys.db.ReadAll(bucketPrefix + username + "-")
	if err != nil {
		return nil, err
	}
	for key, value := range mb {
		a := bucketMetadata{}
		err := json.Unmarshal([]byte(value), &a)
		if err != nil {
			continue
		}
		key = strings.Split(key, "-")[1]
		m = append(m, a)
	}
	return m, nil
}
