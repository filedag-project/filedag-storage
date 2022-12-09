package store

import (
	"context"
	"github.com/filedag-project/filedag-storage/objectservice/pkg/auth"
	"github.com/filedag-project/filedag-storage/objectservice/pkg/policy"
	"github.com/filedag-project/filedag-storage/objectservice/pkg/s3action"
	"github.com/filedag-project/filedag-storage/objectservice/uleveldb"
)

// BucketPolicySys - policy subsystem.
type BucketPolicySys struct {
	BmSys *BucketMetadataSys
}

// Read only object actions.
var readOnlyObjectActions = s3action.Action("s3:GetObject")

// Write only object actions.
var writeOnlyObjectActions = s3action.NewActionSet("s3:AbortMultipartUpload", "s3:DeleteObject", "s3:ListMultipartUploadParts", "s3:PutObject")

// Common bucket actions for both read and write policies.
var commonBucketActions = s3action.Action("s3:GetBucketLocation")

// Read only bucket actions.
var readOnlyBucketActions = s3action.Action("s3:ListBucket")

// Write only bucket actions.
var writeOnlyBucketActions = s3action.Action("s3:ListBucketMultipartUploads")

// NewIPolicySys  - creates new policy system.
func NewIPolicySys(db *uleveldb.ULevelDB) *BucketPolicySys {
	return &BucketPolicySys{
		BmSys: NewBucketMetadataSys(db),
	}
}

// IsAllowed - checks given policy args is allowed to continue the Rest API.
func (sys *BucketPolicySys) IsAllowed(ctx context.Context, args auth.Args) bool {
	meta, err := sys.BmSys.GetBucketMeta(ctx, args.BucketName)
	if err == nil {
		return meta.PolicyConfig.IsAllowed(args)
	}
	if _, ok := err.(BucketPolicyNotFound); !ok {
		log.Debugw("can't find bucket policy", "bucket", args.BucketName)
	}
	return false
}

// BucketPolicy - Bucket level policy.
type BucketPolicy string

// Different types of Policies currently supported for buckets.
const (
	bucketPolicyNone      BucketPolicy = "private"
	bucketPolicyReadOnly               = "download"
	bucketPolicyReadWrite              = "public"
	bucketPolicyWriteOnly              = "upload"
)

// GetPolicyName - Returns policy of given bucket name, prefix in given statements.
func GetPolicyName(statements []policy.Statement, bucketName string, prefix string) BucketPolicy {
	bucketResource := policy.NewResource(bucketName, "")
	objectResource := policy.NewResource(bucketName, "*")

	bucketCommonFound := false
	bucketReadOnly := false
	bucketWriteOnly := false
	matchedResource := ""
	objReadOnly := false
	objWriteOnly := false

	for _, s := range statements {
		matchedObjResources := policy.NewResourceSet()
		if s.Resources.Contains(objectResource) {
			matchedObjResources.Add(objectResource)
		}
		if !matchedObjResources.IsEmpty() {
			readOnly, writeOnly := getObjectPolicy(s)
			for resource := range matchedObjResources {
				if len(matchedResource) < len(resource.String()) {
					objReadOnly = readOnly
					objWriteOnly = writeOnly
					matchedResource = resource.String()
				} else if len(matchedResource) == len(resource.String()) {
					objReadOnly = objReadOnly || readOnly
					objWriteOnly = objWriteOnly || writeOnly
					matchedResource = resource.String()
				}
			}
		}
		if s.Resources.Contains(bucketResource) {
			commonFound, readOnly, writeOnly := getBucketPolicy(s, prefix)
			bucketCommonFound = bucketCommonFound || commonFound
			bucketReadOnly = bucketReadOnly || readOnly
			bucketWriteOnly = bucketWriteOnly || writeOnly
		}
	}

	bucketPolicy := bucketPolicyNone
	if bucketCommonFound {
		if bucketReadOnly && bucketWriteOnly && objReadOnly && objWriteOnly {
			bucketPolicy = bucketPolicyReadWrite
		} else if bucketReadOnly && objReadOnly {
			bucketPolicy = bucketPolicyReadOnly
		} else if bucketWriteOnly && objWriteOnly {
			bucketPolicy = bucketPolicyWriteOnly
		}
	}

	return bucketPolicy
}

// Returns policy of given object statement.
func getObjectPolicy(statement policy.Statement) (readOnly bool, writeOnly bool) {
	if statement.Effect == "Allow" &&
		statement.Principal.AWS.Contains("*") &&
		statement.Conditions == nil {
		if statement.Actions.Contains(readOnlyObjectActions) {
			readOnly = true
		}
		for _, a := range writeOnlyObjectActions.ToSlice() {
			if !statement.Actions.Contains(a) {
				return readOnly, writeOnly
			}
		}
		writeOnly = true
	}

	return readOnly, writeOnly
}

// Returns policy of given bucket statement.
func getBucketPolicy(statement policy.Statement, prefix string) (commonFound, readOnly, writeOnly bool) {
	if !(statement.Effect == "Allow" && statement.Principal.AWS.Contains("*")) {
		return commonFound, readOnly, writeOnly
	}

	if statement.Actions.Contains(commonBucketActions) &&
		statement.Conditions == nil {
		commonFound = true
	}

	if statement.Actions.Contains(writeOnlyBucketActions) &&
		statement.Conditions == nil {
		writeOnly = true
	}

	if statement.Actions.Contains(readOnlyBucketActions) {
		readOnly = true
	}

	return commonFound, readOnly, writeOnly
}

//GetSelfPolicy get default policy
func GetSelfPolicy(accessKey, bucket string) (policy.Statement, error) {
	var sta = policy.Statement{
		SID:        "",
		Effect:     "Allow",
		Principal:  policy.NewPrincipal(accessKey),
		Actions:    s3action.NewActionSet(s3action.AllActions),
		Resources:  policy.NewResourceSet(policy.NewResource(bucket, "*")),
		Conditions: nil,
	}
	return sta, nil
}
