package dagpooluser

import (
	"github.com/filedag-project/filedag-storage/dag/pool/userpolicy"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
)

type IdentityUserSys struct {
	DB *uleveldb.ULevelDB
}

const dagPoolUser = "dagPoolUser/"

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type DagPoolUser struct {
	Username string
	Password string
	Policy   userpolicy.DagPoolPolicy
	Capacity uint64
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
func (i *IdentityUserSys) QueryUser(username string) (DagPoolUser, error) {
	var u DagPoolUser
	err := i.DB.Get(dagPoolUser+username, &u)
	if err != nil {
		return u, err
	}
	return u, nil
}

// UpdateUser Update user
func (i *IdentityUserSys) UpdateUser(u DagPoolUser) error {
	err := i.DB.Put(dagPoolUser+u.Username, u)
	if err != nil {
		return err
	}
	return nil
}
func (i *IdentityUserSys) CheckUserPolicy(username, pass string, policy userpolicy.DagPoolPolicy) bool {
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
func NewIdentityUserSys(db *uleveldb.ULevelDB) (IdentityUserSys, error) {
	return IdentityUserSys{db}, nil
}
