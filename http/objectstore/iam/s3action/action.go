package s3action

import "github.com/filedag-project/filedag-storage/http/objectstore/iam/set"

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

	// AllActions - all API actions
	AllActions = "s3:*"
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
