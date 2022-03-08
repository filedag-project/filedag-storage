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
	userDB       *uleveldb.ULeveldb
	policyDB     *uleveldb.ULeveldb
	userPolicyDB *uleveldb.ULeveldb
	groupDB      *uleveldb.ULeveldb
}

func (I *iamLevelDBStore) init() {
	I.userDB = uleveldb.GlobalUserLevelDB
	I.policyDB = uleveldb.GlobalPolicyLevelDB
	I.userPolicyDB = uleveldb.GlobalUserPolicyLevelDB
	I.groupDB = uleveldb.GlobalGroupLevelDB
}
func (I *iamLevelDBStore) loadUser(ctx context.Context, user string, m *auth.Credentials) error {
	err := I.userDB.Get(user, m)
	if err != nil {
		return err
	}
	return nil
}

func (I *iamLevelDBStore) loadUsers(ctx context.Context) (map[string]auth.Credentials, error) {
	m := make(map[string]auth.Credentials)

	mc, err := I.userDB.ReadAll()
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
	err := I.userDB.Put(name, u.Credentials)
	if err != nil {
		return err
	}
	return nil
}

func (I *iamLevelDBStore) RemoveUserIdentity(ctx context.Context, name string) error {
	err := I.userDB.Delete(name)
	if err != nil {
		return err
	}
	return nil
}
func (I *iamLevelDBStore) createPolicy(ctx context.Context, policyName string, policyDocument policy.PolicyDocument) error {
	err := I.policyDB.Put(policyName, policyDocument)
	if err != nil {
		return err
	}
	return nil
}
func (I *iamLevelDBStore) createUserPolicy(ctx context.Context, userName, policyName string, policyDocument policy.PolicyDocument) error {
	err := I.userPolicyDB.Put(userName+policyName, policyDocument)
	if err != nil {
		return err
	}
	return nil
}

func (I *iamLevelDBStore) getUserPolicy(ctx context.Context, userName, policyName string, policyDocument policy.PolicyDocument) error {
	err := I.userPolicyDB.Get(userName+policyName, policyDocument)
	if err != nil {
		return err
	}
	return nil
}
func (I *iamLevelDBStore) removeUserPolicy(ctx context.Context, userName, policyName string) error {
	err := I.userPolicyDB.Delete(userName + policyName)
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
		userDB: uleveldb.GlobalUserLevelDB,
	}
}
