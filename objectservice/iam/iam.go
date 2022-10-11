package iam

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/filedag-project/filedag-storage/objectservice/iam/auth"
	"github.com/filedag-project/filedag-storage/objectservice/iam/policy"
	"github.com/filedag-project/filedag-storage/objectservice/iam/s3action"
	"github.com/filedag-project/filedag-storage/objectservice/uleveldb"
	logging "github.com/ipfs/go-log/v2"
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
	// Persistence layer for IAM subsystem
	store *iamStoreSys
}

// NewIdentityAMSys - new an IdentityAM config system
func NewIdentityAMSys(db *uleveldb.ULevelDB) *IdentityAMSys {
	sys := &IdentityAMSys{}
	sys.store = &iamStoreSys{newIAMLevelDBStore(db)}
	// TODO: Is it necessary?
	//err := sys.store.saveUserIdentity(context.Background(), auth.DefaultAccessKey, UserIdentity{Credentials: auth.GetDefaultActiveCred()})
	//if err != nil {
	//	log.Errorf("add first user fail%v", err)
	//}
	return sys
}

// IsAllowed - checks given policy args is allowed to continue the Rest API.
func (sys *IdentityAMSys) IsAllowed(ctx context.Context, args auth.Args) bool {

	// Policies don't apply to the owner.
	if args.IsOwner {
		return true
	}
	// If the credential is temporary, perform STS related checks.
	ok, parentUser, err := sys.IsTempUser(ctx, args.AccountName)
	if err != nil {
		return false
	}
	if ok {
		return sys.IsAllowedSTS(args, parentUser)
	}
	// Continue with the assumption of a regular user
	keys, err := sys.store.getUserPolices(ctx, args.AccountName)
	if err != nil {
		return false
	}
	if len(keys) == 0 {
		// No policy found.
		return false
	}
	var pod policy.PolicyDocument
	err = sys.GetUserPolicy(ctx, args.AccountName, "default", &pod)
	if err != nil {
		return false
	}
	for _, pn := range keys {
		if pn == "default" {
			continue
		}
		var po policy.PolicyDocument
		err1 := sys.GetUserPolicy(ctx, args.AccountName, pn, &po)
		if err1 != nil {
			return false
		}

		pod = pod.Merge(po)
	}
	// Policies were found, evaluate all of them.
	perm := policy.Policy{

		Version:    pod.Version,
		Statements: pod.Statement,
	}
	return perm.IsAllowed(args)
}

// IsAllowedSTS is meant for STS based temporary credentials,
// which implements claims validation and verification other than
// applying policies.
func (sys *IdentityAMSys) IsAllowedSTS(args auth.Args, parentUser string) bool {
	//todo check parentUser policy
	return true
}

// IsTempUser - returns if given key is a temporary user.
func (sys *IdentityAMSys) IsTempUser(ctx context.Context, name string) (bool, string, error) {
	cred, found := sys.GetUser(ctx, name)
	if !found {
		return false, "", errNoSuchUser
	}
	if cred.IsExpired() {
		err := sys.store.removeUserIdentity(ctx, name)
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
func (sys *IdentityAMSys) GetUserList(ctx context.Context, accressKey string) ([]*iam.User, error) {
	var u []*iam.User
	users, err := sys.store.loadUsers(ctx)
	if err != nil {
		return nil, err
	}
	for i := range users {
		cerd := users[i]
		if cerd.IsExpired() {
			continue
		}
		var a = iam.User{
			Arn:                 nil,
			CreateDate:          &cerd.CreateTime,
			PasswordLastUsed:    nil,
			Path:                nil,
			PermissionsBoundary: nil,
			Tags:                nil,
			UserId:              &cerd.AccessKey,
			UserName:            &cerd.AccessKey,
		}
		u = append(u, &a)
	}
	return u, err
}

//AddUser add user
func (sys *IdentityAMSys) AddUser(ctx context.Context, accessKey, secretKey string) error {
	m := make(map[string]interface{})
	credentials, err := auth.CreateNewCredentialsWithMetadata(accessKey, secretKey, m, auth.DefaultSecretKey)
	if err != nil {
		log.Errorf("Create NewCredentials WithMetadata err:%v,%v,%v", accessKey, secretKey, err)
		return err
	}
	p := policy.CreateUserPolicy(accessKey, []s3action.Action{s3action.CreateBucketAction, s3action.ListBucketAction, s3action.ListAllMyBucketsAction}, "*")
	err = sys.store.createUserPolicy(ctx, accessKey, "default", policy.PolicyDocument{
		Version:   p.Version,
		Statement: p.Statements,
	})
	if err != nil {
		return err
	}
	err = sys.store.saveUserIdentity(ctx, accessKey, UserIdentity{credentials})
	if err != nil {
		log.Errorf("save UserIdentity err:%v", err)
		return err
	}

	return nil
}

//AddSubUser add user
func (sys *IdentityAMSys) AddSubUser(ctx context.Context, accessKey, secretKey, parentUser string) error {
	m := make(map[string]interface{})
	credentials, err := auth.CreateNewCredentialsWithMetadata(accessKey, secretKey, m, auth.DefaultSecretKey)
	if err != nil {
		log.Errorf("Create NewCredentials WithMetadata err:%v,%v,%v", accessKey, secretKey, err)
		return err
	}
	credentials.ParentUser = parentUser
	err = sys.store.saveUserIdentity(ctx, accessKey, UserIdentity{credentials})
	if err != nil {
		log.Errorf("save UserIdentity err:%v", err)
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
	err := sys.store.createUserPolicy(ctx, userName, policyName, policyDocument)
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
	keys, err := sys.store.getUserPolices(ctx, userName)
	if err != nil {
		log.Errorf("get UserPolicy err:%v", err)
		return nil, err
	}
	return keys, nil
}

//UpdateUserPolicy update user policy
func (sys *IdentityAMSys) UpdateUserPolicy(ctx context.Context, username string, pname string, p *policy.Policy) error {
	return sys.PutUserPolicy(ctx, username, pname, policy.PolicyDocument{
		Version:   p.Version,
		Statement: p.Statements,
	})

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
func (sys *IdentityAMSys) GetUserInfo(ctx context.Context, accessKey string) (cred auth.Credentials, err error) {
	err = sys.store.loadUser(ctx, accessKey, &cred)
	return
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

//func (sys *IdentityAMSys) CreateGroup(ctx context.Context, groupName string, version int) error {
//	err := sys.store.CreateGroup(ctx, groupName, version)
//	if err != nil {
//		return err
//	}
//	return nil
//}
//func (sys *IdentityAMSys) GetGroup(ctx context.Context, groupName string) (GroupInfo, error) {
//	g, err := sys.store.GetGroup(ctx, groupName)
//	if err != nil {
//		return g, err
//	}
//	return g, nil
//}
//func (sys *IdentityAMSys) DeleteGroup(ctx context.Context, groupName string) error {
//	err := sys.store.DeleteGroup(ctx, groupName)
//	if err != nil {
//		return err
//	}
//	return nil
//}
//func (sys *IdentityAMSys) ListGroups(ctx context.Context, path string) ([]GroupInfo, error) {
//	m, err := sys.store.loadGroups(ctx)
//	var s []GroupInfo
//	for _, v := range m {
//		s = append(s, v)
//	}
//	if err != nil {
//		return s, err
//	}
//	return s, nil
//}
