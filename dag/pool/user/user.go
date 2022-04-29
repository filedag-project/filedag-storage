package user

import (
	"github.com/filedag-project/filedag-storage/dag/pool/userpolicy"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
)

type IdentityUser struct {
	DB *uleveldb.ULevelDB
}

const dagPoolUser = "dagPoolUser/"

type DagPoolUser struct {
	username string
	password string
	policy   userpolicy.DagPoolPolicy
	capacity uint64
}

// AddUser add user
func (i *IdentityUser) AddUser(user DagPoolUser) error {
	err := i.DB.Put(dagPoolUser+user.username, user)
	if err != nil {
		return err
	}
	return nil
}

// RemoveUser remove user
func (i *IdentityUser) RemoveUser(username string) error {
	err := i.DB.Delete(dagPoolUser + username)
	if err != nil {
		return err
	}
	return nil
}

// QueryUser query user
func (i *IdentityUser) QueryUser(username string) (DagPoolUser, error) {
	var u DagPoolUser
	err := i.DB.Get(dagPoolUser+username, &u)
	if err != nil {
		return u, err
	}
	return u, nil
}

// UpdateUser Update user
func (i *IdentityUser) UpdateUser(u DagPoolUser) error {
	err := i.DB.Put(dagPoolUser+u.username, u)
	if err != nil {
		return err
	}
	return nil
}
func (i *IdentityUser) CheckUserPolicy(username, pass string, policy userpolicy.DagPoolPolicy) bool {
	user, err := i.QueryUser(username)
	if err != nil {
		return false
	}
	if user.password != pass {
		return false
	}
	if !user.policy.Allow(policy) {
		return false
	}
	return true
}
func NewIdentityUser() (IdentityUser, error) {
	db, err := uleveldb.OpenDb("./")
	if err != nil {
		return IdentityUser{}, err
	}
	return IdentityUser{db}, nil
}
