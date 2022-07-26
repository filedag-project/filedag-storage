package dpuser

import (
	"github.com/filedag-project/filedag-storage/dag/pool/poolservice/dpuser/upolicy"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
)

//IdentityUserSys identity user sys
type IdentityUserSys struct {
	DB           *uleveldb.ULevelDB
	rootUser     string
	rootPassword string
}

const dagPoolUser = "dagPoolUser/"

//DagPoolUser DagPool User
type DagPoolUser struct {
	Username string
	Password string
	Policy   upolicy.DagPoolPolicy
	Capacity uint64
}

//CheckAdmin check user admin policy
func (i *IdentityUserSys) CheckAdmin(user, pass string) bool {
	return i.rootUser == user && i.rootPassword == pass
}

//IsAdmin check user if admin user
func (i *IdentityUserSys) IsAdmin(user string) bool {
	return i.rootUser == user
}

//CheckUser check user if correct
func (i *IdentityUserSys) CheckUser(user, pass string) bool {
	if i.CheckAdmin(user, pass) {
		return true
	}
	queryUser, err := i.QueryUser(user)
	if err != nil {
		return false
	}
	if queryUser.Password != pass {
		return false
	}
	return true
}

// AddUser add user
func (i *IdentityUserSys) AddUser(user DagPoolUser) error {
	err := i.DB.Put(dagPoolUser+user.Username, user)
	if err != nil {
		return err
	}
	return nil
}

// RemoveUser remove user
func (i *IdentityUserSys) RemoveUser(username string) error {
	err := i.DB.Delete(dagPoolUser + username)
	if err != nil {
		return err
	}
	return nil
}

// QueryUser query user
func (i *IdentityUserSys) QueryUser(username string) (*DagPoolUser, error) {
	var u DagPoolUser
	err := i.DB.Get(dagPoolUser+username, &u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// UpdateUser Update user
func (i *IdentityUserSys) UpdateUser(u DagPoolUser) error {
	err := i.DB.Put(dagPoolUser+u.Username, u)
	if err != nil {
		return err
	}
	return nil
}

//CheckUserPolicy check user policy
func (i *IdentityUserSys) CheckUserPolicy(username, pass string, policy userpolicy.DagPoolPolicy) bool {
	if i.CheckAdmin(username, pass) {
		return true
	}
	user, err := i.QueryUser(username)
	if err != nil {
		return false
	}
	if user.Password != pass {
		return false
	}
	if !user.Policy.Allow(policy) {
		return false
	}
	return true
}

//NewIdentityUserSys new identity user sys
func NewIdentityUserSys(db *uleveldb.ULevelDB, rootUser, rootPassword string) (*IdentityUserSys, error) {
	return &IdentityUserSys{
		DB:           db,
		rootUser:     rootUser,
		rootPassword: rootPassword,
	}, nil
}
