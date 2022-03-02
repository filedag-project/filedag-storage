package iam

import (
	"context"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/auth"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
)

// IAMLevelDBStore implements IAMStorageAPI
type IAMLevelDBStore struct {
	db *uleveldb.Uleveldb
}

func (I *IAMLevelDBStore) loadPolicyDoc(ctx context.Context, policy string, m map[string]PolicyDoc) error {
	//TODO implement me
	panic("implement me")
}

func (I *IAMLevelDBStore) loadPolicyDocs(ctx context.Context, m map[string]PolicyDoc) error {
	//TODO implement me
	panic("implement me")
}

func (I *IAMLevelDBStore) loadUser(ctx context.Context, user string, m map[string]auth.Credentials) error {
	//TODO implement me
	panic("implement me")
}

func (I *IAMLevelDBStore) loadUsers(ctx context.Context, m map[string]auth.Credentials) error {
	//TODO implement me
	panic("implement me")
}

func (I *IAMLevelDBStore) loadGroup(ctx context.Context, group string, m map[string]GroupInfo) error {
	//TODO implement me
	panic("implement me")
}

func (I *IAMLevelDBStore) loadGroups(ctx context.Context, m map[string]GroupInfo) error {
	//TODO implement me
	panic("implement me")
}

func (I *IAMLevelDBStore) loadMappedPolicy(ctx context.Context, name string, isGroup bool, m map[string]MappedPolicy) error {
	//TODO implement me
	panic("implement me")
}

func (I *IAMLevelDBStore) loadMappedPolicies(ctx context.Context, isGroup bool, m map[string]MappedPolicy) error {
	//TODO implement me
	panic("implement me")
}

func (I *IAMLevelDBStore) savePolicyDoc(ctx context.Context, policyName string, p PolicyDoc) error {
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

func (I *IAMLevelDBStore) deletePolicyDoc(ctx context.Context, policyName string) error {
	//TODO implement me
	panic("implement me")
}

func (I *IAMLevelDBStore) deleteMappedPolicy(ctx context.Context, name string, isGroup bool) error {
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
