package store

import (
	"context"
	"errors"
	"github.com/filedag-project/filedag-storage/objectservice/iam/policy"
	"github.com/syndtr/goleveldb/leveldb"
	"golang.org/x/xerrors"
)

// UpdateBucketPolicy Update bucket metadata .
// The configData data should not be modified after being sent here.
func (sys *BucketMetadataSys) UpdateBucketPolicy(ctx context.Context, bucket string, p *policy.Policy) error {
	meta, err := sys.GetBucketMeta(bucket)
	if err != nil {
		return err
	}
	meta.PolicyConfig = p
	if p == nil {
		return xerrors.Errorf("policy is nil")
	}
	err = sys.UpdateBucket(bucket, &meta)
	if err != nil {
		return err
	}
	return nil
}

// DeleteBucketPolicy Delete bucket metadata .
// The configData data should not be modified after being sent here.
func (sys *BucketMetadataSys) DeleteBucketPolicy(ctx context.Context, bucket string) error {
	meta, err := sys.GetBucketMeta(bucket)
	if err != nil {
		return err
	}
	err = sys.UpdateBucket(bucket, &BucketMetadata{
		Name:          bucket,
		Region:        meta.Region,
		Created:       meta.Created,
		PolicyConfig:  nil,
		TaggingConfig: meta.TaggingConfig,
	})
	if err != nil {
		return err
	}
	return nil
}

// GetPolicyConfig returns configured bucket policy
func (sys *BucketMetadataSys) GetPolicyConfig(bucket string) (*policy.Policy, error) {
	meta, err := sys.GetBucketMeta(bucket)
	if err != nil {
		if errors.Is(err, leveldb.ErrNotFound) {
			return nil, bucketPolicyNotFound{Bucket: bucket}
		}
		return nil, err
	}
	if meta.PolicyConfig == nil {
		return nil, bucketPolicyNotFound{Bucket: bucket}
	}
	return meta.PolicyConfig, nil
}
