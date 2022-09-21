package store

import (
	"context"
)

func (sys *BucketMetadataSys) UpdateBucketTagging(ctx context.Context, bucket string, tags *Tags) error {
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

	meta.TaggingConfig = tags
	return sys.setBucketMeta(bucket, &meta)
}

func (sys *BucketMetadataSys) DeleteBucketTagging(ctx context.Context, bucket string) error {
	return sys.UpdateBucketPolicy(ctx, bucket, nil)
}

func (sys *BucketMetadataSys) GetTaggingConfig(ctx context.Context, bucket string) (*Tags, error) {
	meta, err := sys.GetBucketMeta(ctx, bucket)
	if err != nil {
		switch err.(type) {
		case BucketNotFound:
			return nil, BucketTaggingNotFound{Bucket: bucket}
		}
		return nil, err
	}
	if meta.TaggingConfig == nil {
		return nil, BucketTaggingNotFound{Bucket: bucket}
	}
	return meta.TaggingConfig, nil
}
