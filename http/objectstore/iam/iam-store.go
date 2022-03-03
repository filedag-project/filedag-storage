package iam

import (
	"context"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/auth"
)

// IAMStoreAPI defines an interface for the IAM persistence layer
type IAMStoreAPI interface {
	loadUser(ctx context.Context, user string, m map[string]auth.Credentials) error
	loadUsers(ctx context.Context, m map[string]auth.Credentials) error
	loadGroup(ctx context.Context, group string, m map[string]GroupInfo) error
	loadGroups(ctx context.Context, m map[string]GroupInfo) error
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
