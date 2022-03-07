package iam

import (
	"context"
	"encoding/json"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/auth"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/policy"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
)

// iamLevelDBStore implements IAMStorageAPI
type iamLevelDBStore struct {
	db *uleveldb.Uleveldb
}

func (I *iamLevelDBStore) init() {
	I.db = uleveldb.GlobalLevelDB
}
func (I *iamLevelDBStore) loadUser(ctx context.Context, user string, m *auth.Credentials) error {
	err := I.db.Get(user, m)
	if err != nil {
		return err
	}
	return nil
}

func (I *iamLevelDBStore) loadUsers(ctx context.Context) (map[string]auth.Credentials, error) {
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
func (I *iamLevelDBStore) saveUserIdentity(ctx context.Context, name string, u UserIdentity) error {
	err := I.db.Put(name, u.Credentials)
	if err != nil {
		return err
	}
	return nil
}

func (I *iamLevelDBStore) RemoveUserIdentity(ctx context.Context, name string) error {
	err := I.db.Delete(name)
	if err != nil {
		return err
	}
	return nil
}
func (I *iamLevelDBStore) createPolicy(ctx context.Context, policyName string, policyDocument policy.PolicyDocument) error {
	err := I.db.Put(policyName, policyDocument)
	if err != nil {
		return err
	}
	return nil
}
func (I *iamLevelDBStore) createUserPolicy(ctx context.Context, userName, policyName string, policyDocument policy.PolicyDocument) error {
	err := I.db.Put(userName+policyName, policyDocument)
	if err != nil {
		return err
	}
	return nil
}

func (I *iamLevelDBStore) getUserPolicy(ctx context.Context, userName, policyName string, policyDocument policy.PolicyDocument) error {
	err := I.db.Get(userName+policyName, policyDocument)
	if err != nil {
		return err
	}
	return nil
}
func (I *iamLevelDBStore) removeUserPolicy(ctx context.Context, userName, policyName string) error {
	err := I.db.Delete(userName + policyName)
	if err != nil {
		return err
	}
	return nil
}
func (I *iamLevelDBStore) loadGroup(ctx context.Context, group string, m *GroupInfo) error {
	//TODO implement me
	panic("implement me")
}

func (I *iamLevelDBStore) loadGroups(ctx context.Context) (map[string]GroupInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (I *iamLevelDBStore) saveGroupInfo(ctx context.Context, group string, gi GroupInfo) error {
	//TODO implement me
	panic("implement me")
}

func (I *iamLevelDBStore) RemoveGroupInfo(ctx context.Context, name string) error {
	//TODO implement me
	panic("implement me")
}

func newIAMLevelDBStore() *iamLevelDBStore {
	return &iamLevelDBStore{
		db: uleveldb.GlobalLevelDB,
	}
}
