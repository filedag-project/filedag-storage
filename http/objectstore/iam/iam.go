package iam

import (
	"context"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/auth"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/policy"
	logging "github.com/ipfs/go-log/v2"
	"sync"
)

const (
	statusEnabled  = "enabled"
	statusDisabled = "disabled"
)

var log = logging.Logger("iam")
var GlobalIAMSys IAMSys

// IAMSys - config system.
type IAMSys struct {
	sync.Mutex
	// Persistence layer for IAM subsystem
	store *iAMStoreSys
}

// Init - initializes config system
func (sys *IAMSys) Init(ctx context.Context) {
	sys.Lock()
	defer sys.Unlock()
	sys.initStore()

}

// initStore initializes IAM stores
func (sys *IAMSys) initStore() {
	sys.store = &iAMStoreSys{newIAMLevelDBStore()}
}

// IsAllowed - checks given policy args is allowed to continue the Rest API.
func (sys *IAMSys) IsAllowed(args policy.Args) bool {

	// Policies don't apply to the owner.
	if args.IsOwner {
		return true
	}
	m := &auth.Credentials{}
	err := sys.store.loadUser(context.Background(), args.AccountName, m)
	if err != nil {
		return false
	}
	return true
}

// GetUserList all user
func (sys *IAMSys) GetUserList(ctx context.Context) []*iam.User {
	var u []*iam.User
	users, err := sys.store.loadUsers(ctx)
	if err != nil {
		return nil
	}
	for user, cerd := range users {
		u = append(u, &iam.User{
			Arn:                 nil,
			CreateDate:          &cerd.Expiration,
			PasswordLastUsed:    nil,
			Path:                nil,
			PermissionsBoundary: nil,
			Tags:                nil,
			UserId:              &user,
			UserName:            &user,
		})
	}
	return u
}

//AddUser add user
func (sys *IAMSys) AddUser(ctx context.Context, accessKey, secretKey string) error {
	m := make(map[string]interface{})
	credentials, err := auth.CreateNewCredentialsWithMetadata(accessKey, secretKey, m, auth.DefaultSecretKey)
	if err != nil {
		log.Errorf("Create NewCredentials WithMetadata err:%v,%v,%v", accessKey, secretKey, err)
		return err
	}
	err = sys.store.saveUserIdentity(ctx, accessKey, UserIdentity{credentials})
	if err != nil {
		log.Errorf("save UserIdentity err:%v", err)
		return err
	}
	return nil
}

// RemoveUser Remove User
func (sys *IAMSys) RemoveUser(ctx context.Context, accessKey string) error {
	err := sys.store.RemoveUserIdentity(ctx, accessKey)
	if err != nil {
		log.Errorf("Remove UserIdentity err:%v", err)
		return err
	}
	return nil
}

// CreatePolicy Create Policy
func (sys *IAMSys) CreatePolicy(ctx context.Context, policyName string, policyDocument policy.PolicyDocument) error {
	err := sys.store.createPolicy(ctx, policyName, policyDocument)
	if err != nil {
		log.Errorf("create Policy err:%v", err)
		return err
	}
	return nil
}

// PutUserPolicy Create Policy
func (sys *IAMSys) PutUserPolicy(ctx context.Context, userName, policyName string, policyDocument policy.PolicyDocument) error {
	err := sys.store.createUserPolicy(ctx, userName, policyName, policyDocument)
	if err != nil {
		log.Errorf("create UserPolicy err:%v", err)
		return err
	}
	return nil
}

// GetUserPolicy Get User Policy
func (sys *IAMSys) GetUserPolicy(ctx context.Context, userName, policyName string, policyDocument policy.PolicyDocument) error {
	err := sys.store.getUserPolicy(ctx, userName, policyName, policyDocument)
	if err != nil {
		log.Errorf("get UserPolicy err:%v", err)
		return err
	}
	return nil
}

// RemoveUserPolicy remove User Policy
func (sys *IAMSys) RemoveUserPolicy(ctx context.Context, userName, policyName string) error {
	err := sys.store.removeUserPolicy(ctx, userName, policyName)
	if err != nil {
		log.Errorf("remove UserPolicy err:%v", err)
		return err
	}
	return nil
}
