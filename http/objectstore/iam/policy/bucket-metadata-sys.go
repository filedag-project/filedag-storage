package policy

import (
	"errors"
	"github.com/filedag-project/filedag-storage/http/objectstore/berrors"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	"sync"
	"time"
)

var globalBucketMetadataSys = NewBucketMetadataSys()

// GenericBucketError - generic object layer error.
type GenericBucketError struct {
	Bucket string
	Err    error
}

// BucketPolicyNotFound - no bucket policy found.
type BucketPolicyNotFound GenericBucketError

func (e BucketPolicyNotFound) Error() string {
	return "No bucket policy configuration found for bucket: " + e.Bucket
}

// BucketMetadataSys captures all bucket metadata for a given cluster.
type BucketMetadataSys struct {
	sync.RWMutex
	metadataMap map[string]BucketMetadata
}

// NewBucketMetadataSys - creates new policy system.
func NewBucketMetadataSys() *BucketMetadataSys {
	return &BucketMetadataSys{
		metadataMap: make(map[string]BucketMetadata),
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
func (sys *BucketMetadataSys) GetPolicyConfig(bucket string) (*Policy, error) {
	meta, err := sys.GetConfig(bucket)
	if err != nil {
		if errors.Is(err, berrors.ErrConfigNotFound) {
			return nil, BucketPolicyNotFound{Bucket: bucket}
		}
		return nil, err
	}
	if meta.PolicyConfig == nil {
		return nil, BucketPolicyNotFound{Bucket: bucket}
	}
	return meta.PolicyConfig, nil
}

// GetConfig returns a specific configuration from the bucket metadata.
func (sys *BucketMetadataSys) GetConfig(bucket string) (BucketMetadata, error) {
	sys.RLock()
	meta, ok := sys.metadataMap[bucket]
	sys.RUnlock()
	if !ok {
		return BucketMetadata{}, BucketPolicyNotFound{Bucket: bucket}
	}

	return meta, nil
}

// Set - sets a new metadata in-memory.
func (sys *BucketMetadataSys) Set(bucket string, meta BucketMetadata) {
	sys.Lock()
	sys.metadataMap[bucket] = meta
	sys.Unlock()
}

// Get metadata for a bucket.
func (sys *BucketMetadataSys) Get(bucket string) (BucketMetadata, error) {

	sys.RLock()
	defer sys.RUnlock()

	meta, ok := sys.metadataMap[bucket]
	if !ok {
		return newBucketMetadata(bucket), berrors.ErrConfigNotFound
	}

	return meta, nil
}

func (sys *BucketMetadataSys) Update(bucket string, meta BucketMetadata) {
	db := uleveldb.OpenDb("./fds.db")
	db.Put(bucket, meta)
}
