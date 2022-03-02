package iam

import (
	"context"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/auth"
)

// IAMStoreSys defines an interface for the IAM persistence layer
type IAMStoreSys interface {
	loadPolicyDoc(ctx context.Context, policy string, m map[string]PolicyDoc) error
	loadPolicyDocs(ctx context.Context, m map[string]PolicyDoc) error
	loadUser(ctx context.Context, user string, m map[string]auth.Credentials) error
	loadUsers(ctx context.Context, m map[string]auth.Credentials) error
	loadGroup(ctx context.Context, group string, m map[string]GroupInfo) error
	loadGroups(ctx context.Context, m map[string]GroupInfo) error
	loadMappedPolicy(ctx context.Context, name string, isGroup bool, m map[string]MappedPolicy) error
	loadMappedPolicies(ctx context.Context, isGroup bool, m map[string]MappedPolicy) error
	savePolicyDoc(ctx context.Context, policyName string, p PolicyDoc) error
	saveUserIdentity(ctx context.Context, name string, u UserIdentity) error
	saveGroupInfo(ctx context.Context, group string, gi GroupInfo) error
	deletePolicyDoc(ctx context.Context, policyName string) error
	deleteMappedPolicy(ctx context.Context, name string, isGroup bool) error
	deleteUserIdentity(ctx context.Context, name string) error
	deleteGroupInfo(ctx context.Context, name string) error
}
