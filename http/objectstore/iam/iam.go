package iam

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/auth"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/policy"
	"github.com/filedag-project/filedag-storage/http/objectstore/utils"
	logging "github.com/ipfs/go-log/v2"
	"sync"
)

const (
	statusEnabled  = "enabled"
	statusDisabled = "disabled"
)

// error returned to IAM subsystem when user doesn't exist.
var errNoSuchUser = errors.New("specified user does not exist")
var errUserIsExpired = errors.New("specified user is expired")

var log = logging.Logger("iam")

// IdentityAMSys - config system.
type IdentityAMSys struct {
	sync.Mutex
	// Persistence layer for IAM subsystem
	store *iamStoreSys
}

// Init - initializes IdentityAM config system
func (sys *IdentityAMSys) Init() {
	sys.Lock()
	defer sys.Unlock()
	sys.initStore()
}

// initStore initializes IAM stores
func (sys *IdentityAMSys) initStore() {
	sys.store = &iamStoreSys{newIAMLevelDBStore()}

}

// IsAllowed - checks given policy args is allowed to continue the Rest API.
func (sys *IdentityAMSys) IsAllowed(args auth.Args) bool {

	// Policies don't apply to the owner.
	if args.IsOwner {
		return true
	}
	// If the credential is temporary, perform STS related checks.
	ok, parentUser, err := sys.IsTempUser(args.AccountName)
	if err != nil {
		return false
	}
	if ok {
		return sys.IsAllowedSTS(args, parentUser)
	}
	// Continue with the assumption of a regular user
	ps, _, err := sys.store.getUserPolices(context.Background(), args.AccountName)
	if err != nil {
		return false
	}
	var pol, pmer policy.Policy
	for _, p := range ps {
		pmer = pol.Merge(p)
	}
	// Policies were found, evaluate all of them.
	return pmer.IsAllowed(args)
}

// IsAllowedSTS is meant for STS based temporary credentials,
// which implements claims validation and verification other than
// applying policies.
func (sys *IdentityAMSys) IsAllowedSTS(args auth.Args, parentUser string) bool {
	//todo check parentUser policy
	return true
}

// IsTempUser - returns if given key is a temporary user.
func (sys *IdentityAMSys) IsTempUser(name string) (bool, string, error) {
	cred, found := sys.GetUser(context.Background(), name)
	if !found {
		return false, "", errNoSuchUser
	}
	if cred.IsExpired() {
		err := sys.store.removeUserIdentity(context.Background(), name)
		if err != nil {
			return false, "", err
		}
		return false, "", errUserIsExpired
	}
	if cred.IsTemp() {
		return true, cred.ParentUser, nil
	}

	return false, "", nil
}

// GetUserList all user
func (sys *IdentityAMSys) GetUserList(ctx context.Context) []*iam.User {
	var u []*iam.User
	users, err := sys.store.loadUsers(ctx)
	if err != nil {
		return nil
	}
	for _, cerd := range users {
		if cerd.IsExpired() {
			continue
		}
		var a = &iam.User{
			Arn:                 nil,
			CreateDate:          utils.Time(cerd.CreateTime),
			PasswordLastUsed:    nil,
			Path:                nil,
			PermissionsBoundary: nil,
			Tags:                nil,
			UserId:              utils.String(cerd.AccessKey),
			UserName:            utils.String(cerd.AccessKey),
		}
		u = append(u, a)
	}
	return u
}

//AddUser add user
func (sys *IdentityAMSys) AddUser(ctx context.Context, accessKey, secretKey string) error {
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
	err = sys.store.createUserPolicy(ctx, accessKey, policy.DefaultPolicies[1].Name, policy.PolicyDocument{
		Version: policy.DefaultPolicies[1].Definition.Version,
		Statement: []policy.Statement{{
			SID:       policy.DefaultPolicies[1].Definition.Statements[0].SID,
			Effect:    policy.DefaultPolicies[1].Definition.Statements[0].Effect,
			Principal: policy.NewPrincipal(accessKey),
			Actions:   policy.DefaultPolicies[1].Definition.Statements[0].Actions,
			Resources: policy.DefaultPolicies[1].Definition.Statements[0].Resources},
		},
	})
	if err != nil {
		return err
	}
	return nil
}

//UpdateUser Update User
func (sys *IdentityAMSys) UpdateUser(ctx context.Context, cred auth.Credentials) error {
	err := sys.store.saveUserIdentity(ctx, cred.AccessKey, UserIdentity{Credentials: cred})
	if err != nil {
		return err
	}
	return nil
}

// GetUser - get user credentials
func (sys *IdentityAMSys) GetUser(ctx context.Context, accessKey string) (cred auth.Credentials, ok bool) {
	m := auth.Credentials{}
	err := sys.store.loadUser(ctx, accessKey, &m)
	if err != nil {
		return m, false
	}

	return m, m.IsValid()
}

// RemoveUser Remove User
func (sys *IdentityAMSys) RemoveUser(ctx context.Context, accessKey string) error {
	err := sys.store.removeUserIdentity(ctx, accessKey)
	if err != nil {
		log.Errorf("Remove UserIdentity err:%v", err)
		return err
	}
	return nil
}

// CreatePolicy Create Policy
func (sys *IdentityAMSys) CreatePolicy(ctx context.Context, policyName string, policyDocument policy.PolicyDocument) error {
	err := sys.store.createPolicy(ctx, policyName, policyDocument)
	if err != nil {
		log.Errorf("create Policy err:%v", err)
		return err
	}
	return nil
}

// PutUserPolicy Create Policy
func (sys *IdentityAMSys) PutUserPolicy(ctx context.Context, userName, policyName string, policyDocument policy.PolicyDocument) error {
	var pd policy.PolicyDocument
	err := sys.GetUserPolicy(ctx, userName, policyName, &pd)
	pd.Merge(policyDocument)
	err = sys.store.createUserPolicy(ctx, userName, policyName, pd)
	if err != nil {
		log.Errorf("create UserPolicy err:%v", err)
		return err
	}
	return nil
}

// GetUserPolicy Get User Policy
func (sys *IdentityAMSys) GetUserPolicy(ctx context.Context, userName, policyName string, policyDocument *policy.PolicyDocument) error {
	err := sys.store.getUserPolicy(ctx, userName, policyName, policyDocument)
	if err != nil {
		log.Errorf("get UserPolicy err:%v", err)
		return err
	}
	return nil
}

// GetUserPolices Get User all Policy
func (sys *IdentityAMSys) GetUserPolices(ctx context.Context, userName string) ([]string, error) {
	_, keys, err := sys.store.getUserPolices(ctx, userName)
	if err != nil {
		log.Errorf("get UserPolicy err:%v", err)
		return nil, err
	}
	return keys, nil
}

// RemoveUserPolicy remove User Policy
func (sys *IdentityAMSys) RemoveUserPolicy(ctx context.Context, userName, policyName string) error {
	err := sys.store.removeUserPolicy(ctx, userName, policyName)
	if err != nil {
		log.Errorf("remove UserPolicy err:%v", err)
		return err
	}
	return nil
}

// GetUserInfo  - get user info
func (sys *IdentityAMSys) GetUserInfo(ctx context.Context, accessKey string) (cred auth.Credentials, ok bool) {
	m := auth.Credentials{}
	err := sys.store.loadUser(ctx, accessKey, &m)
	if err != nil {
		return m, false
	}

	return m, m.IsValid()
}

// SetTempUser - set temporary user credentials, these credentials have an
// expiry. The permissions for these STS credentials is determined in one of the
// following ways:
func (sys *IdentityAMSys) SetTempUser(ctx context.Context, accessKey string, cred auth.Credentials, policyName string) error {
	err := sys.store.SetTempUser(ctx, accessKey, cred, policyName)
	if err != nil {
		return err
	}
	return nil
}

func (sys *IdentityAMSys) CreateGroup(ctx context.Context, groupName string, version int) error {
	err := sys.store.CreateGroup(ctx, groupName, version)
	if err != nil {
		return err
	}
	return nil
}
func (sys *IdentityAMSys) GetGroup(ctx context.Context, groupName string) (GroupInfo, error) {
	g, err := sys.store.GetGroup(ctx, groupName)
	if err != nil {
		return g, err
	}
	return g, nil
}
func (sys *IdentityAMSys) DeleteGroup(ctx context.Context, groupName string) error {
	err := sys.store.DeleteGroup(ctx, groupName)
	if err != nil {
		return err
	}
	return nil
}
