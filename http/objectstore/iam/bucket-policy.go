package iam

import (
	"context"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/auth"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/policy"
)

// IPolicySys - policy subsystem.
type IPolicySys struct {
	bmSys *bucketMetadataSys
}

// GetAllBucketOfUser returns stored bucket policy
func (sys *IPolicySys) GetAllBucketOfUser(accessKey string) ([]BucketMetadata, error) {
	return sys.bmSys.getAllBucketOfUser(accessKey)
}

// Get returns stored bucket policy
func (sys *IPolicySys) Get(bucket, accessKey string) (*policy.Policy, error) {
	return sys.bmSys.GetPolicyConfig(bucket, accessKey)
}

// GetMeta returns stored bucket GetMeta
func (sys *IPolicySys) GetMeta(bucket, accessKey string) (BucketMetadata, error) {
	return sys.bmSys.GetMeta(bucket, accessKey)
}

// Head returns stored bucket policy
func (sys *IPolicySys) Head(bucket, accessKey string) bool {
	return sys.bmSys.HasBucketMeta(bucket, accessKey)
}

// Set returns stored bucket policy
func (sys *IPolicySys) Set(bucket, accessKey, region string) error {
	return sys.bmSys.SetBucketMeta(bucket, accessKey, newBucketMetadata(bucket, region, accessKey))
}

// IsAllowed - checks given policy args is allowed to continue the Rest API.
func (sys *IPolicySys) IsAllowed(args auth.Args) bool {
	p, err := sys.Get(args.BucketName, args.AccountName)
	if err != nil {
		return false
	} else {
		return p.IsAllowed(args)
	}
}

// Init  - creates new policy system.
func (sys *IPolicySys) Init() {
	sys.bmSys = newBucketMetadataSys()
}

// UpdatePolicy Update bucket metadata .
// The configData data should not be modified after being sent here.
func (sys *IPolicySys) UpdatePolicy(ctx context.Context, accessKey, bucket string, p *policy.Policy) error {
	meta, err := sys.bmSys.GetMeta(bucket, accessKey)
	if err != nil {
		return err
	}
	if meta.PolicyConfig != nil {
		*p = p.Merge(*meta.PolicyConfig)
	}
	err = sys.bmSys.UpdateBucketMeta(accessKey, bucket, &BucketMetadata{
		Name:         bucket,
		PolicyConfig: p,
	})
	if err != nil {
		return err
	}
	return nil
}

// DeletePolicy Delete bucket metadata .
// The configData data should not be modified after being sent here.
func (sys *IPolicySys) DeletePolicy(ctx context.Context, accessKey, bucket string, p *policy.Policy) error {
	meta, err := sys.bmSys.GetMeta(bucket, accessKey)
	if err != nil {
		return err
	}
	err = sys.bmSys.UpdateBucketMeta(accessKey, bucket, &BucketMetadata{
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
func (sys *IPolicySys) UpdateBucketMeta(ctx context.Context, user string, bucket string, tags *Tags) error {
	meta, err := sys.bmSys.GetMeta(bucket, user)
	if err != nil {
		return err
	}
	err = sys.bmSys.UpdateBucketMeta(user, bucket, &BucketMetadata{
		Name:          meta.Name,
		Region:        meta.Region,
		Created:       meta.Created,
		PolicyConfig:  meta.PolicyConfig,
		TaggingConfig: tags,
	})
	if err != nil {
		return err
	}
	return nil
}

//Delete the bucket
func (sys *IPolicySys) Delete(ctx context.Context, accessKey, bucket string) error {
	err := sys.bmSys.DeleteBucketMeta(accessKey, bucket)
	if err != nil {
		return err
	}
	return nil
}
