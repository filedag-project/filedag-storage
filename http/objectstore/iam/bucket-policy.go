package iam

import (
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/auth"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/policy"
	"github.com/filedag-project/filedag-storage/http/objectstore/store"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
)

// IPolicySys - policy subsystem.
type IPolicySys struct {
	bmSys *store.BucketMetadataSys
}

// NewIPolicySys  - creates new policy system.
func NewIPolicySys(db *uleveldb.ULevelDB) *IPolicySys {
	return &IPolicySys{
		bmSys: store.NewBucketMetadataSys(db),
	}
}

// IsAllowed - checks given policy args is allowed to continue the Rest API.
func (sys *IPolicySys) IsAllowed(args auth.Args) bool {
	p, err := sys.bmSys.GetPolicyConfig(args.BucketName, args.AccountName)
	if err != nil {
		return false
	} else {
		return p.IsAllowed(args)
	}
}

// SetPolicy returns stored bucket policy
func (sys *IPolicySys) SetPolicy(bucket, accessKey, region string) error {
	return sys.bmSys.SetBucketMeta(bucket, accessKey, store.NewBucketMetadata(bucket, region, accessKey))
}

// GetPolicy returns stored bucket policy
func (sys *IPolicySys) GetPolicy(bucket, accessKey string) (*policy.Policy, error) {
	return sys.bmSys.GetPolicyConfig(bucket, accessKey)
}
