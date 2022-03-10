package iam

import (
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/auth"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/policy"
)

// IPolicySys - policy subsystem.
type IPolicySys struct {
	bmSys *bucketMetadataSys
}

// Get returns stored bucket policy
func (sys *IPolicySys) GetAllBucketOfUser(accessKey string) ([]bucketMetadata, error) {
	return sys.bmSys.getAllBucketOfUser(accessKey)
}

// Get returns stored bucket policy
func (sys *IPolicySys) Get(bucket, accessKey string) (*policy.Policy, error) {
	return sys.bmSys.GetPolicyConfig(bucket, accessKey)
}

// Set returns stored bucket policy
func (sys *IPolicySys) Set(bucket, accessKey string) error {
	return sys.bmSys.Set(bucket, accessKey, newBucketMetadata(bucket))
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
