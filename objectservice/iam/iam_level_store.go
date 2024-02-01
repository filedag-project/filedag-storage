package iam

import (
	"context"
	"fmt"
	"github.com/filedag-project/filedag-storage/objectservice/objmetadb"
	"github.com/filedag-project/filedag-storage/objectservice/pkg/policy"
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
	levelDB objmetadb.ObjStoreMetaDBAPI
}

func (i *iamLevelDBStore) loadUser(ctx context.Context, user string, m *UserIdentity) error {
	err := i.levelDB.Get(getUserKey(user), m)
	if err != nil {
		return err
	}
	return nil
}

func (i *iamLevelDBStore) loadUsers(ctx context.Context) (map[string]UserIdentity, error) {
	m := make(map[string]UserIdentity)

	all, err := i.levelDB.ReadAllChan(ctx, getUserKey(""), "")
	if err != nil {
		return m, err
	}
	for entry := range all {
		cred := UserIdentity{}
		if err = entry.UnmarshalValue(&cred); err != nil {
			continue
		}
		if cred.Credentials.AccessKey == "" {
			continue
		}
		m[cred.Credentials.AccessKey] = cred
	}
	return m, nil
}

func (i *iamLevelDBStore) saveUserIdentity(ctx context.Context, u UserIdentity) error {
	err := i.levelDB.Put(getUserKey(u.Credentials.AccessKey), u)
	if err != nil {
		return err
	}
	return nil
}

func (i *iamLevelDBStore) removeUserIdentity(ctx context.Context, name string) error {
	err := i.levelDB.Delete(getUserKey(name))
	if err != nil {
		return err
	}
	return nil
}

func (i *iamLevelDBStore) savePolicy(ctx context.Context, policyName string, policyDocument policy.PolicyDocument) error {
	err := i.levelDB.Put(getPolicyKey(policyName), policyDocument)
	if err != nil {
		return err
	}
	return nil
}

func (i *iamLevelDBStore) saveUserPolicy(ctx context.Context, userName, policyName string, policyDocument policy.PolicyDocument) error {
	err := i.levelDB.Put(getUserPolicyKey(userName, policyName), policyDocument)
	if err != nil {
		return err
	}
	return nil
}

func (i *iamLevelDBStore) loadUserPolicy(ctx context.Context, userName, policyName string, policyDocument *policy.PolicyDocument) error {
	err := i.levelDB.Get(getUserPolicyKey(userName, policyName), policyDocument)
	if err != nil {
		return err
	}
	return nil
}

func (i *iamLevelDBStore) loadUserAllPolicies(ctx context.Context, userName string) ([]policy.Policy, []string, error) {
	var ps []policy.Policy
	var keys []string
	all, err := i.levelDB.ReadAllChan(ctx, getUserPolicyKey(userName, ""), "")
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
		k := strings.TrimPrefix(entry.GetKey(), getUserPolicyKey(userName, ""))
		keys = append(keys, k)
	}
	return ps, keys, nil
}

func (i *iamLevelDBStore) removeUserPolicy(ctx context.Context, userName, policyName string) error {
	err := i.levelDB.Delete(getUserPolicyKey(userName, policyName))
	if err != nil {
		return err
	}
	return nil
}

func (i *iamLevelDBStore) removeUserAllPolicies(ctx context.Context, userName string) error {
	all, err := i.levelDB.ReadAllChan(ctx, getUserPolicyKey(userName, ""), "")
	if err != nil {
		return err
	}
	for entry := range all {
		if err = i.levelDB.Delete(entry.GetKey()); err != nil {
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

func newIAMLevelDBStore(db objmetadb.ObjStoreMetaDBAPI) *iamLevelDBStore {
	return &iamLevelDBStore{
		levelDB: db,
	}
}
