package iam

import (
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
func newBucketMetadataSys(db *uleveldb.ULevelDB) *bucketMetadataSys {
	return &bucketMetadataSys{
		db: db,
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
	Name          string
	Region        string
	Owner         string
	Created       time.Time
	PolicyConfig  *policy.Policy
	TaggingConfig *Tags
}

// newBucketMetadata creates BucketMetadata with the supplied name and Created to Now.
func newBucketMetadata(name, region, accessKey string) BucketMetadata {
	var p = policy.Policy{
		ID:      policy.DefaultPolicies[0].Definition.ID,
		Version: policy.DefaultPolicies[0].Definition.Version,
		Statements: []policy.Statement{{
			SID:       policy.DefaultPolicies[0].Definition.Statements[0].SID,
			Effect:    policy.DefaultPolicies[0].Definition.Statements[0].Effect,
			Principal: policy.NewPrincipal("*"),
			Actions:   policy.DefaultPolicies[0].Definition.Statements[0].Actions,
			Resources: policy.NewResourceSet(policy.NewResource(name, "*")),
		}},
	}
	return BucketMetadata{
		Name:         name,
		Region:       region,
		Owner:        accessKey,
		Created:      time.Now().UTC(),
		PolicyConfig: &p,
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
func (sys *bucketMetadataSys) GetMeta(bucket, accessKey string) (BucketMetadata, error) {
	var meta BucketMetadata
	err := sys.GetBucketMeta(bucket, accessKey, &meta)
	if err != nil {
		return BucketMetadata{}, err
	}
	return meta, nil
}

// SetBucketMeta - sets a new metadata in-db
func (sys *bucketMetadataSys) SetBucketMeta(bucket, username string, meta BucketMetadata) error {
	err := sys.db.Put(bucketPrefix+"-"+bucket, meta)
	if err != nil {
		return err
	}
	return nil
}

// GetBucketMeta metadata for a bucket.
func (sys *bucketMetadataSys) GetBucketMeta(bucket, username string, meta *BucketMetadata) error {
	err := sys.db.Get(bucketPrefix+"-"+bucket, meta)
	if err != nil {
		return err
	}
	return nil
}

// HasBucketMeta metadata for a bucket.
func (sys *bucketMetadataSys) HasBucketMeta(bucket, username string) bool {
	var meta BucketMetadata
	err := sys.db.Get(bucketPrefix+"-"+bucket, &meta)
	if err != nil {
		return false
	}
	return true
}

// DeleteBucketMeta bucket.
func (sys *bucketMetadataSys) DeleteBucketMeta(username, bucket string) error {
	err := sys.db.Delete(bucketPrefix + "-" + bucket)
	if err != nil {
		return err
	}
	return nil
}

// UpdateBucketMeta  metadata for a bucket.
func (sys *bucketMetadataSys) UpdateBucketMeta(username, bucket string, meta *BucketMetadata) error {
	err := sys.db.Put(bucketPrefix+"-"+bucket, meta)
	if err != nil {
		return err
	}
	return nil
}

// getAllBucketOfUser metadata for all bucket.
func (sys *bucketMetadataSys) getAllBucketOfUser(username string) ([]BucketMetadata, error) {
	var m []BucketMetadata
	mb, err := sys.db.ReadAll(bucketPrefix + "-")
	if err != nil {
		return nil, err
	}
	for key, value := range mb {
		a := BucketMetadata{}
		err := json.Unmarshal([]byte(value), &a)
		if err != nil {
			continue
		}
		key = strings.Split(key, "-")[1]
		m = append(m, a)
	}
	return m, nil
}
