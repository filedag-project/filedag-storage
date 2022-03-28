package iam

import (
	"context"
	"encoding/json"
	"encoding/xml"
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

// bucketMetadata contains bucket metadata.
type bucketMetadata struct {
	Name          string
	Region        string
	Created       time.Time
	PolicyConfig  *policy.Policy
	taggingConfig *Tags
}

// newBucketMetadata creates bucketMetadata with the supplied name and Created to Now.
func newBucketMetadata(name, region string) bucketMetadata {
	return bucketMetadata{
		Name:    name,
		Region:  region,
		Created: time.Now().UTC(),
	}
}

// GetPolicyConfig returns configured bucket policy
func (sys *bucketMetadataSys) GetPolicyConfig(bucket, accessKey string) (*policy.Policy, error) {
	meta, err := sys.GetMeta(bucket, accessKey)
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

// GetMeta returns a specific configuration from the bucket metadata.
func (sys *bucketMetadataSys) GetMeta(bucket, accessKey string) (bucketMetadata, error) {
	var meta bucketMetadata
	err := sys.GetBucketMeta(bucket, accessKey, &meta)
	if err != nil {
		return bucketMetadata{}, err
	}
	return meta, nil
}

// SetBucketMeta - sets a new metadata in-db
func (sys *bucketMetadataSys) SetBucketMeta(bucket, username string, meta bucketMetadata) error {
	err := sys.db.Put(bucketPrefix+username+"-"+bucket, meta)
	if err != nil {
		return err
	}
	return nil
}

// GetBucketMeta metadata for a bucket.
func (sys *bucketMetadataSys) GetBucketMeta(bucket, username string, meta *bucketMetadata) error {
	err := sys.db.Get(bucketPrefix+username+"-"+bucket, meta)
	if err != nil {
		return err
	}
	return nil
}

// HeadBucketMeta metadata for a bucket.
func (sys *bucketMetadataSys) HeadBucketMeta(bucket, username string) bool {
	var meta bucketMetadata
	err := sys.db.Get(bucketPrefix+username+"-"+bucket, &meta)
	if err != nil {
		return false
	}
	return true
}

// DeleteBucketMeta bucket.
func (sys *bucketMetadataSys) DeleteBucketMeta(username, bucket string) error {
	err := sys.db.Delete(bucketPrefix + username + "-" + bucket)
	if err != nil {
		return err
	}
	return nil
}

// UpdateBucketMeta  metadata for a bucket.
func (sys *bucketMetadataSys) UpdateBucketMeta(username, bucket string, meta *bucketMetadata) error {
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
func (sys *IPolicySys) UpdateBucketMeta(ctx context.Context, user string, bucket string, tags *Tags) error {
	meta, err := sys.bmSys.GetMeta(bucket, user)
	if err != nil {
		return err
	}
	err = sys.bmSys.UpdateBucketMeta(user, bucket, &bucketMetadata{
		Name:          meta.Name,
		Region:        meta.Region,
		Created:       meta.Created,
		PolicyConfig:  meta.PolicyConfig,
		taggingConfig: tags,
	})
	if err != nil {
		return err
	}
	return nil
}
