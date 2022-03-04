package iam

import (
	"context"
	"encoding/json"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/auth"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
)

// IAMLevelDBStore implements IAMStorageAPI
type IAMLevelDBStore struct {
	db *uleveldb.Uleveldb
}

func (I *IAMLevelDBStore) init() {
	I.db = uleveldb.GlobalLevelDB
}
func (I *IAMLevelDBStore) loadUser(ctx context.Context, user string, m *auth.Credentials) error {
	err := I.db.Get(user, m)
	if err != nil {
		return err
	}
	return nil
}

func (I *IAMLevelDBStore) loadUsers(ctx context.Context) (map[string]auth.Credentials, error) {
	m := make(map[string]auth.Credentials)

	mc, err := I.db.ReadAll()
	if err != nil {
		return m, err
	}
	for key, value := range mc {
		a := auth.Credentials{}
		err := json.Unmarshal(value, &a)
		if err != nil {
			continue
		}
		m[key] = a
	}
	return m, nil
}

func (I *IAMLevelDBStore) loadGroup(ctx context.Context, group string, m *GroupInfo) error {
	//TODO implement me
	panic("implement me")
}

func (I *IAMLevelDBStore) loadGroups(ctx context.Context) (map[string]GroupInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (I *IAMLevelDBStore) saveUserIdentity(ctx context.Context, name string, u UserIdentity) error {
	err := I.db.Put(name, u.Credentials)
	if err != nil {
		return err
	}
	return nil
}

func (I *IAMLevelDBStore) saveGroupInfo(ctx context.Context, group string, gi GroupInfo) error {
	//TODO implement me
	panic("implement me")
}

func (I *IAMLevelDBStore) deleteUserIdentity(ctx context.Context, name string) error {
	err := I.db.Delete(name)
	if err != nil {
		return err
	}
	return nil
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
