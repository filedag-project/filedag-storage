package iam

import (
	"context"
	"errors"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/auth"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/policy"
)

// errInvalidArgument means that input argument is invalid.
var errInvalidArgument = errors.New("Invalid arguments specified")

// iamStoreAPI defines an interface for the IAM persistence layer
type iamStoreAPI interface {
	saveUserIdentity(ctx context.Context, name string, u UserIdentity) error
	removeUserIdentity(ctx context.Context, name string) error
	loadUser(ctx context.Context, user string, m *auth.Credentials) error
	loadUsers(ctx context.Context) (map[string]auth.Credentials, error)
	loadGroup(ctx context.Context, group string, m *GroupInfo) error
	loadGroups(ctx context.Context) (map[string]GroupInfo, error)
	saveGroupInfo(ctx context.Context, group string, gi GroupInfo) error
	removeGroupInfo(ctx context.Context, name string) error
	createPolicy(ctx context.Context, policyName string, policyDocument policy.PolicyDocument) error
	createUserPolicy(ctx context.Context, userName, policyName string, policyDocument policy.PolicyDocument) error
	getUserPolicy(ctx context.Context, userName, policyName string, policyDocument *policy.PolicyDocument) error
	getUserPolices(ctx context.Context, userName string) ([]policy.Policy, []string, error)
	removeUserPolicy(ctx context.Context, userName, policyName string) error
}

// iamStoreSys contains IAMStorageAPI to add higher-level methods on the storage
// layer.
type iamStoreSys struct {
	iamStoreAPI
}

// SetTempUser - saves temporary (STS) credential to storage and cache. If a
// policy name is given, it is associated with the parent user specified in the
// credential.
func (store *iamStoreSys) SetTempUser(ctx context.Context, accessKey string, cred auth.Credentials, policyName string) error {
	if accessKey == "" || !cred.IsTemp() || cred.IsExpired() || cred.ParentUser == "" {
		return errInvalidArgument
	}
	if policyName != "" {
		//todo policy
	}

	u := newUserIdentity(cred)
	err := store.saveUserIdentity(ctx, accessKey, u)
	if err != nil {
		return err
	}
	//todo policy name
	err = store.createUserPolicy(ctx, accessKey, policy.DefaultPolicies[1].Name, policy.PolicyDocument{
		Version:   policy.DefaultPolicies[1].Definition.Version,
		Statement: policy.DefaultPolicies[1].Definition.Statements,
	})
	if err != nil {
		return err
	}
	return nil
}
func (store *iamStoreSys) CreateGroup(ctx context.Context, groupName string, version int) error {
	var g = GroupInfo{
		Name:    groupName,
		Version: version,
		Status:  "on",
		Members: nil,
	}
	err := store.saveGroupInfo(ctx, groupName, g)
	if err != nil {
		return err
	}
	return nil
}
func (store *iamStoreSys) GetGroup(ctx context.Context, groupName string) (GroupInfo, error) {
	var g GroupInfo
	err := store.loadGroup(ctx, groupName, &g)
	if err != nil {
		return g, err
	}
	return g, nil
}
func (store *iamStoreSys) DeleteGroup(ctx context.Context, groupName string) error {
	err := store.removeGroupInfo(ctx, groupName)
	if err != nil {
		return err
	}
	return nil
}
