package iam

import (
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/auth"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/policy"
	"github.com/filedag-project/filedag-storage/http/objectstore/store"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
)

// iPolicySys - policy subsystem.
type iPolicySys struct {
	bmSys *store.BucketMetadataSys
}

// newIPolicySys  - creates new policy system.
func newIPolicySys(db *uleveldb.ULevelDB) *iPolicySys {
	return &iPolicySys{
		bmSys: store.NewBucketMetadataSys(db),
	}
}

// isAllowed - checks given policy args is allowed to continue the Rest API.
func (sys *iPolicySys) isAllowed(args auth.Args) bool {
	p, err := sys.bmSys.GetPolicyConfig(args.BucketName, args.AccountName)
	if err != nil {
		return false
	} else {
		return p.IsAllowed(args)
	}
}

// SetPolicy sets bucket policy
func (sys *iPolicySys) SetPolicy(bucket, accessKey, region string) error {
	return sys.bmSys.SetBucketMeta(bucket, accessKey, store.NewBucketMetadata(bucket, region, accessKey))
}

// GetPolicy returns stored bucket policy
func (sys *iPolicySys) GetPolicy(bucket, accessKey string) (*policy.Policy, error) {
	return sys.bmSys.GetPolicyConfig(bucket, accessKey)
}
