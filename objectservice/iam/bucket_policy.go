package iam

import (
	"github.com/filedag-project/filedag-storage/objectservice/iam/auth"
	"github.com/filedag-project/filedag-storage/objectservice/iam/policy"
	"github.com/filedag-project/filedag-storage/objectservice/store"
	"github.com/filedag-project/filedag-storage/objectservice/uleveldb"
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
	p, err := sys.bmSys.GetPolicyConfig(args.BucketName)
	if err != nil {
		return false
	} else {
		return p.IsAllowed(args)
	}
}

// SetDefaultPolicy sets bucket policy
func (sys *iPolicySys) SetDefaultPolicy(bucket, accessKey, region string) error {
	return sys.bmSys.SetBucketMeta(bucket, store.NewBucketMetadata(bucket, region, accessKey))
}

// GetPolicy returns stored bucket policy
func (sys *iPolicySys) GetPolicy(bucket string) (*policy.Policy, error) {
	return sys.bmSys.GetPolicyConfig(bucket)
}
