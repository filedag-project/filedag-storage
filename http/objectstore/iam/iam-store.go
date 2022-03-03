package iam

import (
	"context"
)

// IAMStoreAPI defines an interface for the IAM persistence layer
type IAMStoreAPI interface {
	loadUser(ctx context.Context, user string) error
	loadUsers(ctx context.Context) error
	loadGroup(ctx context.Context, group string) error
	loadGroups(ctx context.Context) error
	saveUserIdentity(ctx context.Context, name string, u UserIdentity) error
	saveGroupInfo(ctx context.Context, group string, gi GroupInfo) error
	deleteUserIdentity(ctx context.Context, name string) error
	deleteGroupInfo(ctx context.Context, name string) error
}

// IAMStoreSys contains IAMStorageAPI to add higher-level methods on the storage
// layer.
type IAMStoreSys struct {
	IAMStoreAPI
}
