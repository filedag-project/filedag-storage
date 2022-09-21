package iam

import (
	"context"
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
func (sys *iPolicySys) isAllowed(ctx context.Context, args auth.Args) bool {
	p, err := sys.bmSys.GetPolicyConfig(ctx, args.BucketName)
	if err != nil {
		return false
	} else {
		return p.IsAllowed(args)
	}
}

// GetPolicy returns stored bucket policy
func (sys *iPolicySys) GetPolicy(ctx context.Context, bucket string) (*policy.Policy, error) {
	return sys.bmSys.GetPolicyConfig(ctx, bucket)
}
