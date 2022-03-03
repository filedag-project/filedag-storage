package iam

import (
	"context"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/auth"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/policy"
	"sync"
)

const (
	statusEnabled  = "enabled"
	statusDisabled = "disabled"
)

var GlobalIAMSys IAMSys

// IAMSys - config system.
type IAMSys struct {
	sync.Mutex
	// Persistence layer for IAM subsystem
	store *IAMStoreSys
}

// Init - initializes config system
func (sys *IAMSys) Init(ctx context.Context) {
	sys.Lock()
	defer sys.Unlock()
	sys.initStore()

}

// initStore initializes IAM stores
func (sys *IAMSys) initStore() {
	sys.store = &IAMStoreSys{newIAMLevelDBStore()}
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

// IsAllowed - checks given policy args is allowed to continue the Rest API.
func (sys *IAMSys) GetUserList(ctx context.Context) []*iam.User {
	var u []*iam.User
	users, err := sys.store.loadUsers(ctx)
	if err != nil {
		return nil
	}
	for user, cerd := range users {
		u = append(u, &iam.User{
			Arn:                 &cerd.SessionToken,
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
