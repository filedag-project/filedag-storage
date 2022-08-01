package store

import (
	"encoding/json"
	"encoding/xml"
	"github.com/filedag-project/filedag-storage/http/objectstore/api_errors"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/policy"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/policy/condition"
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

// BucketMetadataSys captures all bucket metadata for a given cluster.
type BucketMetadataSys struct {
	db *uleveldb.ULevelDB
}

// NewBucketMetadataSys - creates new policy system.
func NewBucketMetadataSys(db *uleveldb.ULevelDB) *BucketMetadataSys {
	return &BucketMetadataSys{
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

// NewBucketMetadata creates BucketMetadata with the supplied name and Created to Now.
func NewBucketMetadata(name, region, accessKey string) BucketMetadata {
	equalsFunc, err := condition.NewStringEqualsFunc("", condition.S3Prefix.ToKey(), name)
	if err != nil {
		return BucketMetadata{}
	}
	var p = policy.Policy{
		ID:      policy.DefaultPolicies[0].Definition.ID,
		Version: policy.DefaultPolicies[0].Definition.Version,
		Statements: []policy.Statement{{
			SID:        policy.DefaultPolicies[0].Definition.Statements[0].SID,
			Effect:     policy.DefaultPolicies[0].Definition.Statements[0].Effect,
			Principal:  policy.NewPrincipal(accessKey),
			Actions:    policy.DefaultPolicies[0].Definition.Statements[0].Actions,
			Resources:  policy.NewResourceSet(policy.NewResource(name, "*")),
			Conditions: condition.NewConFunctions(equalsFunc),
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

// SetBucketMeta - sets a new metadata in-db
func (sys *BucketMetadataSys) SetBucketMeta(bucket, username string, meta BucketMetadata) error {
	err := sys.db.Put(bucketPrefix+"-"+bucket, meta)
	if err != nil {
		return err
	}
	return nil
}

// GetBucketMeta metadata for a bucket.
func (sys *BucketMetadataSys) GetBucketMeta(bucket, username string) (meta *BucketMetadata, err error) {
	err = sys.db.Get(bucketPrefix+"-"+bucket, meta)
	return meta, err
}

// HasBucket  metadata for a bucket.
func (sys *BucketMetadataSys) HasBucket(bucket, username string) bool {
	var meta BucketMetadata
	err := sys.db.Get(bucketPrefix+"-"+bucket, &meta)
	if err != nil {
		return false
	}
	return true
}

// DeleteBucket bucket.
func (sys *BucketMetadataSys) DeleteBucket(username, bucket string) error {
	err := sys.db.Delete(bucketPrefix + "-" + bucket)
	if err != nil {
		return err
	}
	return nil
}

// UpdateBucket  metadata for a bucket.
func (sys *BucketMetadataSys) UpdateBucket(username, bucket string, meta *BucketMetadata) error {
	err := sys.db.Put(bucketPrefix+"-"+bucket, meta)
	if err != nil {
		return err
	}
	return nil
}

// GetAllBucketOfUser metadata for all bucket.
func (sys *BucketMetadataSys) GetAllBucketOfUser(username string) ([]BucketMetadata, error) {
	var m []BucketMetadata
	mb, err := sys.db.ReadAll(bucketPrefix + "-")
	if err != nil {
		return nil, err
	}
	for _, value := range mb {
		a := BucketMetadata{}
		err := json.Unmarshal([]byte(value), &a)
		if err != nil {
			continue
		}
		if a.Owner != username {
			continue
		}
		m = append(m, a)
	}
	return m, nil
}
