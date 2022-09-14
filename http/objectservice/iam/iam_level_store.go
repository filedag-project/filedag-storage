package iam

import (
	"context"
	"fmt"
	"github.com/filedag-project/filedag-storage/http/objectservice/iam/auth"
	"github.com/filedag-project/filedag-storage/http/objectservice/iam/policy"
	"github.com/filedag-project/filedag-storage/http/objectservice/uleveldb"
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

	all, err := I.levelDB.ReadAllChan(ctx, userPrefix, "")
	if err != nil {
		return m, err
	}
	for entry := range all {
		cred := auth.Credentials{}
		if err = entry.UnmarshalValue(&cred); err != nil {
			continue
		}
		strs := strings.Split(entry.Key, "/")
		if len(strs) < 2 {
			return nil, fmt.Errorf("invalid key[%s], missing user name", entry.Key)
		}
		key := strs[1]
		m[key] = cred
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

func (I *iamLevelDBStore) removeUserIdentity(ctx context.Context, name string) error {
	err := I.levelDB.Delete(userPrefix + name)
	if err != nil {
		return err
	}
	return nil
}
func (I *iamLevelDBStore) createPolicy(ctx context.Context, policyName string, policyDocument policy.PolicyDocument) error {
	err := I.levelDB.Put(policyPrefix+"-"+policyName, policyDocument)
	if err != nil {
		return err
	}
	return nil
}
func (I *iamLevelDBStore) createUserPolicy(ctx context.Context, userName, policyName string, policyDocument policy.PolicyDocument) error {
	err := I.levelDB.Put(userPolicyPrefix+userName+"-"+policyName, policyDocument)
	if err != nil {
		return err
	}
	return nil
}

func (I *iamLevelDBStore) getUserPolicy(ctx context.Context, userName, policyName string, policyDocument *policy.PolicyDocument) error {
	err := I.levelDB.Get(userPolicyPrefix+userName+"-"+policyName, policyDocument)
	if err != nil {
		return err
	}
	return nil
}
func (I *iamLevelDBStore) getUserPolices(ctx context.Context, userName string) ([]policy.Policy, []string, error) {
	var ps []policy.Policy
	var keys []string
	all, err := I.levelDB.ReadAllChan(ctx, userPolicyPrefix+userName, "")
	if err != nil {
		return nil, nil, err
	}
	for entry := range all {
		var p policy.PolicyDocument
		if err = entry.UnmarshalValue(&p); err != nil {
			return nil, nil, err
		}
		ps = append(ps, policy.Policy{
			ID:         "",
			Version:    p.Version,
			Statements: p.Statement,
		})
		k := strings.TrimPrefix(entry.Key, userPolicyPrefix+userName+"-")
		keys = append(keys, k)
	}
	return ps, keys, nil
}
func (I *iamLevelDBStore) removeUserPolicy(ctx context.Context, userName, policyName string) error {
	err := I.levelDB.Delete(userPolicyPrefix + userName + "-" + policyName)
	if err != nil {
		return err
	}
	return nil
}

//func (I *iamLevelDBStore) loadGroup(ctx context.Context, group string, m *GroupInfo) error {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (I *iamLevelDBStore) loadGroups(ctx context.Context) (map[string]GroupInfo, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (I *iamLevelDBStore) saveGroupInfo(ctx context.Context, group string, gi GroupInfo) error {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (I *iamLevelDBStore) removeGroupInfo(ctx context.Context, name string) error {
//	//TODO implement me
//	panic("implement me")
//}

func newIAMLevelDBStore(db *uleveldb.ULevelDB) *iamLevelDBStore {
	return &iamLevelDBStore{
		levelDB: db,
	}
}
