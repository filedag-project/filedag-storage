package iam

import (
	"context"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/auth"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/policy"
)

// iamStoreAPI defines an interface for the IAM persistence layer
type iamStoreAPI interface {
	saveUserIdentity(ctx context.Context, name string, u UserIdentity) error
	RemoveUserIdentity(ctx context.Context, name string) error
	loadUser(ctx context.Context, user string, m *auth.Credentials) error
	loadUsers(ctx context.Context) (map[string]auth.Credentials, error)
	loadGroup(ctx context.Context, group string, m *GroupInfo) error
	loadGroups(ctx context.Context) (map[string]GroupInfo, error)
	saveGroupInfo(ctx context.Context, group string, gi GroupInfo) error
	RemoveGroupInfo(ctx context.Context, name string) error
	createPolicy(ctx context.Context, policyName string, policyDocument policy.PolicyDocument) error
	createUserPolicy(ctx context.Context, userName, policyName string, policyDocument policy.PolicyDocument) error
	getUserPolicy(ctx context.Context, userName, policyName string, policyDocument policy.PolicyDocument) error
	removeUserPolicy(ctx context.Context, userName, policyName string) error
}

// iamStoreSys contains IAMStorageAPI to add higher-level methods on the storage
// layer.
type iamStoreSys struct {
	iamStoreAPI
}
