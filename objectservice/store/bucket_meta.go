package store

import (
	"context"
	"github.com/filedag-project/filedag-storage/objectservice/pkg/policy"
)

// UpdateBucketPolicy Update bucket metadata .
// The configData data should not be modified after being sent here.
func (sys *BucketMetadataSys) UpdateBucketPolicy(ctx context.Context, bucket string, p *policy.Policy) error {
	lk := sys.NewNSLock(bucket)
	lkctx, err := lk.GetLock(ctx, globalOperationTimeout)
	if err != nil {
		return err
	}
	ctx = lkctx.Context()
	defer lk.Unlock(lkctx.Cancel)

	meta, err := sys.getBucketMeta(bucket)
	if err != nil {
		return err
	}

	meta.PolicyConfig = p
	return sys.setBucketMeta(bucket, &meta)
}

// DeleteBucketPolicy Delete bucket metadata .
// The configData data should not be modified after being sent here.
func (sys *BucketMetadataSys) DeleteBucketPolicy(ctx context.Context, bucket string) error {
	return sys.UpdateBucketPolicy(ctx, bucket, nil)
}

// GetPolicyConfig returns configured bucket policy
func (sys *BucketMetadataSys) GetPolicyConfig(ctx context.Context, bucket string) (*policy.Policy, error) {
	meta, err := sys.GetBucketMeta(ctx, bucket)
	if err != nil {
		switch err.(type) {
		case BucketNotFound:
			return nil, BucketPolicyNotFound{Bucket: bucket}
		}
		return nil, err
	}
	if meta.PolicyConfig == nil {
		return nil, BucketPolicyNotFound{Bucket: bucket}
	}
	return meta.PolicyConfig, nil
}
