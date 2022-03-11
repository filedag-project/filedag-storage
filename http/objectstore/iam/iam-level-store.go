package iam

import (
	"context"
	"encoding/json"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/auth"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/policy"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	"strings"
)

const (
	userPrefix       = "user/"
	policyPrefix     = "policy/"
	userPolicyPrefix = "user_policy/"
	groupPrefix      = "group/"
)

// iamLevelDBStore implements IAMStorageAPI
type iamLevelDBStore struct {
	levelDB *uleveldb.ULevelDB
}

func (I *iamLevelDBStore) loadUser(ctx context.Context, user string, m *auth.Credentials) error {
	err := I.levelDB.Get(userPrefix+user, m)
	if err != nil {
		return err
	}
	return nil
}

func (I *iamLevelDBStore) loadUsers(ctx context.Context) (map[string]auth.Credentials, error) {
	m := make(map[string]auth.Credentials)

	mc, err := I.levelDB.ReadAll(userPrefix)
	if err != nil {
		return m, err
	}
	for key, value := range mc {
		a := auth.Credentials{}
		err := json.Unmarshal([]byte(value), &a)
		if err != nil {
			continue
		}
		key = strings.Split(key, "/")[1]
		m[key] = a
	}
	return m, nil
}
func (I *iamLevelDBStore) saveUserIdentity(ctx context.Context, name string, u UserIdentity) error {
	err := I.levelDB.Put(userPrefix+name, u.Credentials)
	if err != nil {
		return err
	}
	return nil
}

func (I *iamLevelDBStore) RemoveUserIdentity(ctx context.Context, name string) error {
	err := I.levelDB.Delete(userPrefix + name)
	if err != nil {
		return err
	}
	return nil
}
func (I *iamLevelDBStore) createPolicy(ctx context.Context, policyName string, policyDocument policy.PolicyDocument) error {
	err := I.levelDB.Put(policyPrefix+policyName, policyDocument)
	if err != nil {
		return err
	}
	return nil
}
func (I *iamLevelDBStore) createUserPolicy(ctx context.Context, userName, policyName string, policyDocument policy.PolicyDocument) error {
	err := I.levelDB.Put(userPolicyPrefix+userName+policyName, policyDocument)
	if err != nil {
		return err
	}
	return nil
}

func (I *iamLevelDBStore) getUserPolicy(ctx context.Context, userName, policyName string, policyDocument policy.PolicyDocument) error {
	err := I.levelDB.Get(userPrefix+userName+policyName, policyDocument)
	if err != nil {
		return err
	}
	return nil
}
func (I *iamLevelDBStore) removeUserPolicy(ctx context.Context, userName, policyName string) error {
	err := I.levelDB.Delete(userPrefix + userName + policyName)
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
		levelDB: uleveldb.DBClient,
	}
}
