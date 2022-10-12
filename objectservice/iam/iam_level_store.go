package iam

import (
	"context"
	"fmt"
	"github.com/filedag-project/filedag-storage/objectservice/iam/auth"
	"github.com/filedag-project/filedag-storage/objectservice/iam/policy"
	"github.com/filedag-project/filedag-storage/objectservice/uleveldb"
	"strings"
)

const (
	userKeyFormat       = "user/%s"
	policyKeyFormat     = "policy/%s"
	userPolicyKeyFormat = "user_policy/%s/%s"
	groupPrefix         = "group/"
)

func getUserKey(username string) string {
	return fmt.Sprintf(userKeyFormat, username)
}

func getPolicyKey(policyName string) string {
	return fmt.Sprintf(policyKeyFormat, policyName)
}

func getUserPolicyKey(username, policyName string) string {
	return fmt.Sprintf(userPolicyKeyFormat, username, policyName)
}

// iamLevelDBStore implements IAMStorageAPI
type iamLevelDBStore struct {
	levelDB *uleveldb.ULevelDB
}

func (I *iamLevelDBStore) loadUser(ctx context.Context, user string, m *auth.Credentials) error {
	err := I.levelDB.Get(getUserKey(user), m)
	if err != nil {
		return err
	}
	return nil
}

func (I *iamLevelDBStore) loadUsers(ctx context.Context) (map[string]auth.Credentials, error) {
	m := make(map[string]auth.Credentials)

	all, err := I.levelDB.ReadAllChan(ctx, getUserKey(""), "")
	if err != nil {
		return m, err
	}
	for entry := range all {
		cred := auth.Credentials{}
		if err = entry.UnmarshalValue(&cred); err != nil {
			continue
		}
		m[cred.AccessKey] = cred
	}
	return m, nil
}

func (I *iamLevelDBStore) saveUserIdentity(ctx context.Context, u UserIdentity) error {
	err := I.levelDB.Put(getUserKey(u.Credentials.AccessKey), u.Credentials)
	if err != nil {
		return err
	}
	return nil
}

func (I *iamLevelDBStore) removeUserIdentity(ctx context.Context, name string) error {
	err := I.levelDB.Delete(getUserKey(name))
	if err != nil {
		return err
	}
	return nil
}

func (I *iamLevelDBStore) savePolicy(ctx context.Context, policyName string, policyDocument policy.PolicyDocument) error {
	err := I.levelDB.Put(getPolicyKey(policyName), policyDocument)
	if err != nil {
		return err
	}
	return nil
}

func (I *iamLevelDBStore) saveUserPolicy(ctx context.Context, userName, policyName string, policyDocument policy.PolicyDocument) error {
	err := I.levelDB.Put(getUserPolicyKey(userName, policyName), policyDocument)
	if err != nil {
		return err
	}
	return nil
}

func (I *iamLevelDBStore) loadUserPolicy(ctx context.Context, userName, policyName string, policyDocument *policy.PolicyDocument) error {
	err := I.levelDB.Get(getUserPolicyKey(userName, policyName), policyDocument)
	if err != nil {
		return err
	}
	return nil
}

func (I *iamLevelDBStore) loadUserAllPolicies(ctx context.Context, userName string) ([]policy.Policy, []string, error) {
	var ps []policy.Policy
	var keys []string
	all, err := I.levelDB.ReadAllChan(ctx, getUserPolicyKey(userName, ""), "")
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
		k := strings.TrimPrefix(entry.Key, getUserPolicyKey(userName, ""))
		keys = append(keys, k)
	}
	return ps, keys, nil
}

func (I *iamLevelDBStore) removeUserPolicy(ctx context.Context, userName, policyName string) error {
	err := I.levelDB.Delete(getUserPolicyKey(userName, policyName))
	if err != nil {
		return err
	}
	return nil
}

func (I *iamLevelDBStore) removeUserAllPolicies(ctx context.Context, userName string) error {
	all, err := I.levelDB.ReadAllChan(ctx, getUserPolicyKey(userName, ""), "")
	if err != nil {
		return err
	}
	for entry := range all {
		if err = I.levelDB.Delete(entry.Key); err != nil {
			return err
		}
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
