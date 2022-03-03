package iam

import (
	"context"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
)

// IAMLevelDBStore implements IAMStorageAPI
type IAMLevelDBStore struct {
	db *uleveldb.Uleveldb
}

func (I *IAMLevelDBStore) loadUser(ctx context.Context, user string) error {
	//TODO implement me
	panic("implement me")
}

func (I *IAMLevelDBStore) loadUsers(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (I *IAMLevelDBStore) loadGroup(ctx context.Context, group string) error {
	//TODO implement me
	panic("implement me")
}

func (I *IAMLevelDBStore) loadGroups(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (I *IAMLevelDBStore) saveUserIdentity(ctx context.Context, name string, u UserIdentity) error {
	//TODO implement me
	panic("implement me")
}

func (I *IAMLevelDBStore) saveGroupInfo(ctx context.Context, group string, gi GroupInfo) error {
	//TODO implement me
	panic("implement me")
}

func (I *IAMLevelDBStore) deleteUserIdentity(ctx context.Context, name string) error {
	//TODO implement me
	panic("implement me")
}

func (I *IAMLevelDBStore) deleteGroupInfo(ctx context.Context, name string) error {
	//TODO implement me
	panic("implement me")
}

func newIAMLevelDBStore() *IAMLevelDBStore {
	return &IAMLevelDBStore{
		db: uleveldb.GlobalLevelDB,
	}
}
