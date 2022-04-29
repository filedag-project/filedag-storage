package user

import (
	"github.com/filedag-project/filedag-storage/dag/pool/userpolicy"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
)

type IdentityUser struct {
	db uleveldb.ULevelDB
}

const dagPoolUser = "dagPoolUser/"

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type DagPoolUser struct {
	username string
	password string
	policy   userpolicy.DagPoolPolicy
	capacity uint64
}

// AddUser add user
func (i *IdentityUser) AddUser(user DagPoolUser) error {
	err := i.db.Put(dagPoolUser+user.username, user)
	if err != nil {
		return err
	}
	return nil
}

// RemoveUser remove user
func (i *IdentityUser) RemoveUser(username string) error {
	err := i.db.Delete(dagPoolUser + username)
	if err != nil {
		return err
	}
	return nil
}

// QueryUser query user
func (i *IdentityUser) QueryUser(username string) (DagPoolUser, error) {
	var u DagPoolUser
	err := i.db.Get(dagPoolUser+username, &u)
	if err != nil {
		return u, err
	}
	return u, nil
}

// UpdateUser Update user
func (i *IdentityUser) UpdateUser(u DagPoolUser) error {
	err := i.db.Put(dagPoolUser+u.username, u)
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
