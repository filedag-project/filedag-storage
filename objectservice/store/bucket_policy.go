package store

import (
	"context"
	"github.com/filedag-project/filedag-storage/objectservice/iam/auth"
	"github.com/filedag-project/filedag-storage/objectservice/iam/policy"
	"github.com/filedag-project/filedag-storage/objectservice/uleveldb"
)

// BucketPolicySys - policy subsystem.
type BucketPolicySys struct {
	BmSys *BucketMetadataSys
}

// NewIPolicySys  - creates new policy system.
func NewIPolicySys(db *uleveldb.ULevelDB) *BucketPolicySys {
	return &BucketPolicySys{
		BmSys: NewBucketMetadataSys(db),
	}
}

// IsAllowed - checks given policy args is allowed to continue the Rest API.
func (sys *BucketPolicySys) IsAllowed(ctx context.Context, args auth.Args) bool {
	p, err := sys.BmSys.GetPolicyConfig(ctx, args.BucketName)
	if err == nil {
		return p.IsAllowed(args)
	}
	if _, ok := err.(BucketPolicyNotFound); !ok {
		log.Errorw("can't find bucket policy", "bucket", args.BucketName)
	}
	return false
}

// GetPolicy returns stored bucket policy
func (sys *BucketPolicySys) GetPolicy(ctx context.Context, bucket string) (*policy.Policy, error) {
	return sys.BmSys.GetPolicyConfig(ctx, bucket)
}
