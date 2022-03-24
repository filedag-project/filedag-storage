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
func (sys *IPolicySys) GetAllBucketOfUser(accessKey string) ([]bucketMetadata, error) {
	return sys.bmSys.getAllBucketOfUser(accessKey)
}

// Get returns stored bucket policy
func (sys *IPolicySys) Get(bucket, accessKey string) (*policy.Policy, error) {
	return sys.bmSys.GetPolicyConfig(bucket, accessKey)
}

// GetLocation returns stored bucket GetLocation
func (sys *IPolicySys) GetLocation(bucket, accessKey string) (bucketMetadata, error) {
	return sys.bmSys.GetConfig(bucket, accessKey)
}

// Head returns stored bucket policy
func (sys *IPolicySys) Head(bucket, accessKey string) bool {
	return sys.bmSys.Head(bucket, accessKey)
}

// Set returns stored bucket policy
func (sys *IPolicySys) Set(bucket, accessKey, region string) error {
	return sys.bmSys.Set(bucket, accessKey, newBucketMetadata(bucket, region))
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

// UpdatePolicy Update bucket metadata for the specified config file.
// The configData data should not be modified after being sent here.
func (sys *IPolicySys) UpdatePolicy(ctx context.Context, accessKey, bucket string, p *policy.Policy) error {
	pConfig, err := sys.bmSys.GetPolicyConfig(bucket, accessKey)
	if err != nil {
		return err
	}
	p.Merge(*pConfig)
	err = sys.bmSys.Update(accessKey, bucket, &bucketMetadata{
		Name:         bucket,
		PolicyConfig: p,
	})
	if err != nil {
		return err
	}
	return nil
}

//Delete the bucket
func (sys *IPolicySys) Delete(ctx context.Context, accessKey, bucket string) error {
	err := sys.bmSys.Delete(accessKey, bucket)
	if err != nil {
		return err
	}
	return nil
}
