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

//GlobalIAMSys iam system
var GlobalIAMSys iamSys

// iamSys - config system.
type iamSys struct {
	sync.Mutex
	// Persistence layer for IAM subsystem
	store *iamStoreSys
}

// Init - initializes config system
func (sys *iamSys) Init(ctx context.Context) {
	sys.Lock()
	defer sys.Unlock()
	sys.initStore()
}

// initStore initializes IAM stores
func (sys *iamSys) initStore() {
	sys.store = &iamStoreSys{newIAMLevelDBStore()}
}

// IsAllowed - checks given policy args is allowed to continue the Rest API.
func (sys *iamSys) IsAllowed(args policy.Args) bool {

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
func (sys *iamSys) GetUserList(ctx context.Context) []*iam.User {
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
func (sys *iamSys) AddUser(ctx context.Context, accessKey, secretKey string) error {
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

// GetUser - get user credentials
func (sys *iamSys) GetUser(ctx context.Context, accessKey string) (cred auth.Credentials, ok bool) {
	m := auth.Credentials{}
	err := sys.store.loadUser(ctx, accessKey, &m)
	if err != nil {
		return m, false
	}

	return m, cred.IsValid()
}

// RemoveUser Remove User
func (sys *iamSys) RemoveUser(ctx context.Context, accessKey string) error {
	err := sys.store.RemoveUserIdentity(ctx, accessKey)
	if err != nil {
		log.Errorf("Remove UserIdentity err:%v", err)
		return err
	}
	return nil
}

// CreatePolicy Create Policy
func (sys *iamSys) CreatePolicy(ctx context.Context, policyName string, policyDocument policy.PolicyDocument) error {
	err := sys.store.createPolicy(ctx, policyName, policyDocument)
	if err != nil {
		log.Errorf("create Policy err:%v", err)
		return err
	}
	return nil
}

// PutUserPolicy Create Policy
func (sys *iamSys) PutUserPolicy(ctx context.Context, userName, policyName string, policyDocument policy.PolicyDocument) error {
	err := sys.store.createUserPolicy(ctx, userName, policyName, policyDocument)
	if err != nil {
		log.Errorf("create UserPolicy err:%v", err)
		return err
	}
	return nil
}

// GetUserPolicy Get User Policy
func (sys *iamSys) GetUserPolicy(ctx context.Context, userName, policyName string, policyDocument policy.PolicyDocument) error {
	err := sys.store.getUserPolicy(ctx, userName, policyName, policyDocument)
	if err != nil {
		log.Errorf("get UserPolicy err:%v", err)
		return err
	}
	return nil
}

// RemoveUserPolicy remove User Policy
func (sys *iamSys) RemoveUserPolicy(ctx context.Context, userName, policyName string) error {
	err := sys.store.removeUserPolicy(ctx, userName, policyName)
	if err != nil {
		log.Errorf("remove UserPolicy err:%v", err)
		return err
	}
	return nil
}
