package s3action

import (
	"github.com/filedag-project/filedag-storage/objectservice/iam/policy/condition"
	"github.com/filedag-project/filedag-storage/objectservice/iam/set"
)

// ActionSet - set of actions.
//https://docs.aws.amazon.com/service-authorization/latest/reference/list_amazons3.html#amazons3-actions-as-permissions
const (
	// AbortMultipartUploadAction - AbortMultipartUpload Rest API action.
	AbortMultipartUploadAction Action = "s3:AbortMultipartUpload"

	// CreateBucketAction - CreateBucket Rest API action.
	CreateBucketAction = "s3:CreateBucket"

	// DeleteBucketAction - DeleteBucket Rest API action.
	DeleteBucketAction = "s3:DeleteBucket"

	// ForceDeleteBucketAction - DeleteBucket Rest API action when x-FileDagStorage-force-delete flag
	// is specified.
	ForceDeleteBucketAction = "s3:ForceDeleteBucket"

	// DeleteBucketPolicyAction - DeleteBucketPolicy Rest API action.
	DeleteBucketPolicyAction = "s3:DeleteBucketPolicy"

	// DeleteObjectAction - DeleteObject Rest API action.
	DeleteObjectAction = "s3:DeleteObject"

	// GetBucketLocationAction - GetBucketLocation Rest API action.
	GetBucketLocationAction = "s3:GetBucketLocation"

	// GetBucketNotificationAction - GetBucketNotification Rest API action.
	GetBucketNotificationAction = "s3:GetBucketNotification"

	// GetBucketPolicyAction - GetBucketPolicy Rest API action.
	GetBucketPolicyAction = "s3:GetBucketPolicy"

	// GetObjectAction - GetObject Rest API action.
	GetObjectAction = "s3:GetObject"

	// HeadBucketAction - HeadBucket Rest API action. This action is unused in FileDagStorage.
	HeadBucketAction = "s3:HeadBucket"

	// ListAllMyBucketsAction - ListAllMyBuckets (List buckets) Rest API action.
	ListAllMyBucketsAction = "s3:ListAllMyBuckets"

	// ListBucketAction - ListBucket Rest API action.
	ListBucketAction = "s3:ListBucket"

	// GetBucketPolicyStatusAction - Retrieves the policy status for a bucket.
	GetBucketPolicyStatusAction = "s3:GetBucketPolicyStatus"

	// ListBucketVersionsAction - ListBucketVersions Rest API action.
	ListBucketVersionsAction = "s3:ListBucketVersions"

	// ListBucketMultipartUploadsAction - ListMultipartUploads Rest API action.
	ListBucketMultipartUploadsAction = "s3:ListBucketMultipartUploads"

	// ListenNotificationAction - ListenNotification Rest API action.
	// This is FileDagStorage extension.
	ListenNotificationAction = "s3:ListenNotification"

	// ListenBucketNotificationAction - ListenBucketNotification Rest API action.
	// This is FileDagStorage extension.
	ListenBucketNotificationAction = "s3:ListenBucketNotification"

	// ListMultipartUploadPartsAction - ListParts Rest API action.
	ListMultipartUploadPartsAction = "s3:ListMultipartUploadParts"

	// PutBucketLifecycleAction - PutBucketLifecycle Rest API action.
	PutBucketLifecycleAction = "s3:PutLifecycleConfiguration"

	// GetBucketLifecycleAction - GetBucketLifecycle Rest API action.
	GetBucketLifecycleAction = "s3:GetLifecycleConfiguration"

	// PutBucketNotificationAction - PutObjectNotification Rest API action.
	PutBucketNotificationAction = "s3:PutBucketNotification"

	// PutBucketPolicyAction - PutBucketPolicy Rest API action.
	PutBucketPolicyAction = "s3:PutBucketPolicy"

	// PutObjectAction - PutObject Rest API action.
	PutObjectAction = "s3:PutObject"

	// DeleteObjectVersionAction - DeleteObjectVersion Rest API action.
	DeleteObjectVersionAction = "s3:DeleteObjectVersion"

	// DeleteObjectVersionTaggingAction - DeleteObjectVersionTagging Rest API action.
	DeleteObjectVersionTaggingAction = "s3:DeleteObjectVersionTagging"

	// GetObjectVersionAction - GetObjectVersionAction Rest API action.
	GetObjectVersionAction = "s3:GetObjectVersion"

	// GetObjectVersionTaggingAction - GetObjectVersionTagging Rest API action.
	GetObjectVersionTaggingAction = "s3:GetObjectVersionTagging"

	// PutObjectVersionTaggingAction - PutObjectVersionTagging Rest API action.
	PutObjectVersionTaggingAction = "s3:PutObjectVersionTagging"

	// BypassGovernanceRetentionAction - bypass governance retention for PutObjectRetention, PutObject and DeleteObject Rest API action.
	BypassGovernanceRetentionAction = "s3:BypassGovernanceRetention"

	// PutObjectRetentionAction - PutObjectRetention Rest API action.
	PutObjectRetentionAction = "s3:PutObjectRetention"

	// GetObjectRetentionAction - GetObjectRetention, GetObject, HeadObject Rest API action.
	GetObjectRetentionAction = "s3:GetObjectRetention"

	// GetObjectLegalHoldAction - GetObjectLegalHold, GetObject Rest API action.
	GetObjectLegalHoldAction = "s3:GetObjectLegalHold"

	// PutObjectLegalHoldAction - PutObjectLegalHold, PutObject Rest API action.
	PutObjectLegalHoldAction = "s3:PutObjectLegalHold"

	// GetBucketObjectLockConfigurationAction - GetBucketObjectLockConfiguration Rest API action
	GetBucketObjectLockConfigurationAction = "s3:GetBucketObjectLockConfiguration"

	// PutBucketObjectLockConfigurationAction - PutBucketObjectLockConfiguration Rest API action
	PutBucketObjectLockConfigurationAction = "s3:PutBucketObjectLockConfiguration"

	// GetBucketTaggingAction - GetBucketTagging Rest API action
	GetBucketTaggingAction = "s3:GetBucketTagging"

	// PutBucketTaggingAction - PutBucketTagging Rest API action
	PutBucketTaggingAction = "s3:PutBucketTagging"

	// GetObjectTaggingAction - Get Object Tags API action
	GetObjectTaggingAction = "s3:GetObjectTagging"

	// PutObjectTaggingAction - Put Object Tags API action
	PutObjectTaggingAction = "s3:PutObjectTagging"

	// DeleteObjectTaggingAction - Delete Object Tags API action
	DeleteObjectTaggingAction = "s3:DeleteObjectTagging"

	// PutBucketEncryptionAction - PutBucketEncryption REST API action
	PutBucketEncryptionAction = "s3:PutEncryptionConfiguration"

	// GetBucketEncryptionAction - GetBucketEncryption REST API action
	GetBucketEncryptionAction = "s3:GetEncryptionConfiguration"

	// PutBucketVersioningAction - PutBucketVersioning REST API action
	PutBucketVersioningAction = "s3:PutBucketVersioning"

	// GetBucketVersioningAction - GetBucketVersioning REST API action
	GetBucketVersioningAction = "s3:GetBucketVersioning"
	// GetReplicationConfigurationAction  - GetReplicationConfiguration REST API action
	GetReplicationConfigurationAction = "s3:GetReplicationConfiguration"
	// PutReplicationConfigurationAction  - PutReplicationConfiguration REST API action
	PutReplicationConfigurationAction = "s3:PutReplicationConfiguration"

	// ReplicateObjectAction  - ReplicateObject REST API action
	ReplicateObjectAction = "s3:ReplicateObject"

	// ReplicateDeleteAction  - ReplicateDelete REST API action
	ReplicateDeleteAction = "s3:ReplicateDelete"

	// ReplicateTagsAction  - ReplicateTags REST API action
	ReplicateTagsAction = "s3:ReplicateTags"

	// GetObjectVersionForReplicationAction  - GetObjectVersionForReplication REST API action
	GetObjectVersionForReplicationAction = "s3:GetObjectVersionForReplication"

	// RestoreObjectAction - RestoreObject REST API action
	RestoreObjectAction  = "s3:RestoreObject"
	GetUserInfoAction    = "iam:GetUserInfo"
	RemoveUserAction     = "iam:RemoveUser"
	SetStatusAction      = "iam:SetStatusUser"
	ChangePassWordAction = "iam:ChangePassWordUser"
	// AllActions - all API actions
	AllActions    = "s3:*"
	AllIamActions = "iam:*"
)

// SupportedActions List of all supported actions.
var SupportedActions = map[Action]struct{}{
	AbortMultipartUploadAction:             {},
	CreateBucketAction:                     {},
	DeleteBucketAction:                     {},
	ForceDeleteBucketAction:                {},
	DeleteBucketPolicyAction:               {},
	DeleteObjectAction:                     {},
	GetBucketLocationAction:                {},
	GetBucketNotificationAction:            {},
	GetBucketPolicyAction:                  {},
	GetObjectAction:                        {},
	HeadBucketAction:                       {},
	ListAllMyBucketsAction:                 {},
	ListBucketAction:                       {},
	GetBucketPolicyStatusAction:            {},
	ListBucketVersionsAction:               {},
	ListBucketMultipartUploadsAction:       {},
	ListenNotificationAction:               {},
	ListenBucketNotificationAction:         {},
	ListMultipartUploadPartsAction:         {},
	PutBucketLifecycleAction:               {},
	GetBucketLifecycleAction:               {},
	PutBucketNotificationAction:            {},
	PutBucketPolicyAction:                  {},
	PutObjectAction:                        {},
	BypassGovernanceRetentionAction:        {},
	PutObjectRetentionAction:               {},
	GetObjectRetentionAction:               {},
	GetObjectLegalHoldAction:               {},
	PutObjectLegalHoldAction:               {},
	GetBucketObjectLockConfigurationAction: {},
	PutBucketObjectLockConfigurationAction: {},
	GetBucketTaggingAction:                 {},
	PutBucketTaggingAction:                 {},
	GetObjectVersionAction:                 {},
	GetObjectVersionTaggingAction:          {},
	DeleteObjectVersionAction:              {},
	DeleteObjectVersionTaggingAction:       {},
	PutObjectVersionTaggingAction:          {},
	GetObjectTaggingAction:                 {},
	PutObjectTaggingAction:                 {},
	DeleteObjectTaggingAction:              {},
	PutBucketEncryptionAction:              {},
	GetBucketEncryptionAction:              {},
	PutBucketVersioningAction:              {},
	GetBucketVersioningAction:              {},
	GetReplicationConfigurationAction:      {},
	PutReplicationConfigurationAction:      {},
	ReplicateObjectAction:                  {},
	ReplicateDeleteAction:                  {},
	ReplicateTagsAction:                    {},
	GetObjectVersionForReplicationAction:   {},
	AllActions:                             {},
	GetUserInfoAction:                      {},
	RemoveUserAction:                       {},
	SetStatusAction:                        {},
	ChangePassWordAction:                   {},
	AllIamActions:                          {},
}

// IsValid - checks if action is valid or not.
func (action Action) IsValid() bool {
	for supAction := range SupportedActions {
		if action.Match(supAction) {
			return true
		}
	}
	return false
}

// Match - matches action name with action patter.
func (action Action) Match(a Action) bool {
	return set.Match(string(action), string(a))
}

// List of all supported object actions.
var supportedObjectActions = map[Action]struct{}{
	AbortMultipartUploadAction:     {},
	DeleteObjectAction:             {},
	GetObjectAction:                {},
	ListMultipartUploadPartsAction: {},
	PutObjectAction:                {},
	//BypassGovernanceRetentionAction:      {},
	//PutObjectRetentionAction:             {},
	//GetObjectRetentionAction:             {},
	//PutObjectLegalHoldAction:             {},
	//GetObjectLegalHoldAction:             {},
	GetObjectTaggingAction:    {},
	PutObjectTaggingAction:    {},
	DeleteObjectTaggingAction: {},
	//GetObjectVersionAction:               {},
	//GetObjectVersionTaggingAction:        {},
	//DeleteObjectVersionAction:            {},
	//DeleteObjectVersionTaggingAction:     {},
	//PutObjectVersionTaggingAction:        {},
	//ReplicateObjectAction:                {},
	//ReplicateDeleteAction:                {},
	//ReplicateTagsAction:                  {},
	//GetObjectVersionForReplicationAction: {},
	//RestoreObjectAction:                  {},
}

// IsObjectAction - returns whether action is object type or not.
func (action Action) IsObjectAction() bool {
	_, ok := supportedObjectActions[action]
	return ok
}

func createActionConditionKeyMap() map[Action]condition.KeySet {
	commonKeys := []condition.Key{}
	for _, keyName := range condition.CommonKeys {
		commonKeys = append(commonKeys, keyName.ToKey())
	}

	return map[Action]condition.KeySet{
		AbortMultipartUploadAction: condition.NewKeySet(commonKeys...),

		CreateBucketAction: condition.NewKeySet(commonKeys...),

		DeleteObjectAction: condition.NewKeySet(commonKeys...),

		GetBucketLocationAction: condition.NewKeySet(commonKeys...),

		GetBucketPolicyStatusAction: condition.NewKeySet(commonKeys...),

		GetObjectAction: condition.NewKeySet(
			append([]condition.Key{
				condition.S3XAmzServerSideEncryption.ToKey(),
				condition.S3XAmzServerSideEncryptionCustomerAlgorithm.ToKey(),
			}, commonKeys...)...),

		HeadBucketAction: condition.NewKeySet(commonKeys...),

		ListAllMyBucketsAction: condition.NewKeySet(commonKeys...),

		ListBucketAction: condition.NewKeySet(
			append([]condition.Key{
				condition.S3Prefix.ToKey(),
				condition.S3Delimiter.ToKey(),
				condition.S3MaxKeys.ToKey(),
			}, commonKeys...)...),

		ListBucketVersionsAction: condition.NewKeySet(
			append([]condition.Key{
				condition.S3Prefix.ToKey(),
				condition.S3Delimiter.ToKey(),
				condition.S3MaxKeys.ToKey(),
			}, commonKeys...)...),

		ListBucketMultipartUploadsAction: condition.NewKeySet(commonKeys...),

		ListenNotificationAction: condition.NewKeySet(commonKeys...),

		ListenBucketNotificationAction: condition.NewKeySet(commonKeys...),

		ListMultipartUploadPartsAction: condition.NewKeySet(commonKeys...),

		PutObjectAction: condition.NewKeySet(
			append([]condition.Key{
				condition.S3XAmzCopySource.ToKey(),
				condition.S3XAmzServerSideEncryption.ToKey(),
				condition.S3XAmzServerSideEncryptionCustomerAlgorithm.ToKey(),
				condition.S3XAmzMetadataDirective.ToKey(),
				condition.S3XAmzStorageClass.ToKey(),
				condition.S3ObjectLockRetainUntilDate.ToKey(),
				condition.S3ObjectLockMode.ToKey(),
				condition.S3ObjectLockLegalHold.ToKey(),
				condition.S3RequestObjectTagKeys.ToKey(),
				condition.S3RequestObjectTag.ToKey(),
			}, commonKeys...)...),

		// https://docs.aws.amazon.com/AmazonS3/latest/dev/list_amazons3.html
		// LockLegalHold is not supported with PutObjectRetentionAction
		PutObjectRetentionAction: condition.NewKeySet(
			append([]condition.Key{
				condition.S3ObjectLockRemainingRetentionDays.ToKey(),
				condition.S3ObjectLockRetainUntilDate.ToKey(),
				condition.S3ObjectLockMode.ToKey(),
			}, commonKeys...)...),

		GetObjectRetentionAction: condition.NewKeySet(commonKeys...),
		PutObjectLegalHoldAction: condition.NewKeySet(
			append([]condition.Key{
				condition.S3ObjectLockLegalHold.ToKey(),
			}, commonKeys...)...),
		GetObjectLegalHoldAction: condition.NewKeySet(commonKeys...),

		// https://docs.aws.amazon.com/AmazonS3/latest/dev/list_amazons3.html
		BypassGovernanceRetentionAction: condition.NewKeySet(
			append([]condition.Key{
				condition.S3ObjectLockRemainingRetentionDays.ToKey(),
				condition.S3ObjectLockRetainUntilDate.ToKey(),
				condition.S3ObjectLockMode.ToKey(),
				condition.S3ObjectLockLegalHold.ToKey(),
			}, commonKeys...)...),

		GetBucketObjectLockConfigurationAction: condition.NewKeySet(commonKeys...),
		PutBucketObjectLockConfigurationAction: condition.NewKeySet(commonKeys...),
		GetBucketTaggingAction:                 condition.NewKeySet(commonKeys...),
		PutBucketTaggingAction: condition.NewKeySet(
			append([]condition.Key{
				condition.S3RequestObjectTagKeys.ToKey(),
				condition.S3RequestObjectTag.ToKey(),
			}, commonKeys...)...),
		PutObjectTaggingAction: condition.NewKeySet(
			append([]condition.Key{
				condition.S3RequestObjectTagKeys.ToKey(),
				condition.S3RequestObjectTag.ToKey(),
			}, commonKeys...)...),
		GetObjectTaggingAction: condition.NewKeySet(commonKeys...),
		DeleteObjectTaggingAction: condition.NewKeySet(
			append([]condition.Key{
				condition.S3RequestObjectTagKeys.ToKey(),
				condition.S3RequestObjectTag.ToKey(),
			}, commonKeys...)...),
		PutObjectVersionTaggingAction: condition.NewKeySet(
			append([]condition.Key{
				condition.S3VersionID.ToKey(),
				condition.S3RequestObjectTagKeys.ToKey(),
				condition.S3RequestObjectTag.ToKey(),
			}, commonKeys...)...),
		GetObjectVersionAction: condition.NewKeySet(
			append([]condition.Key{
				condition.S3VersionID.ToKey(),
			}, commonKeys...)...),
		GetObjectVersionTaggingAction: condition.NewKeySet(
			append([]condition.Key{
				condition.S3VersionID.ToKey(),
			}, commonKeys...)...),
		DeleteObjectVersionAction: condition.NewKeySet(
			append([]condition.Key{
				condition.S3VersionID.ToKey(),
			}, commonKeys...)...),
		DeleteObjectVersionTaggingAction: condition.NewKeySet(
			append([]condition.Key{
				condition.S3VersionID.ToKey(),
				condition.S3RequestObjectTagKeys.ToKey(),
				condition.S3RequestObjectTag.ToKey(),
			}, commonKeys...)...),
		GetReplicationConfigurationAction:    condition.NewKeySet(commonKeys...),
		PutReplicationConfigurationAction:    condition.NewKeySet(commonKeys...),
		ReplicateObjectAction:                condition.NewKeySet(commonKeys...),
		ReplicateDeleteAction:                condition.NewKeySet(commonKeys...),
		ReplicateTagsAction:                  condition.NewKeySet(commonKeys...),
		GetObjectVersionForReplicationAction: condition.NewKeySet(commonKeys...),
		RestoreObjectAction:                  condition.NewKeySet(commonKeys...),
	}
}

// ActionConditionKeyMap - holds mapping of supported condition key for an action.
var ActionConditionKeyMap = createActionConditionKeyMap()
