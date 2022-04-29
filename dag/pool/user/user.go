package user

import "github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"

type dagPoolPolicy string
type IdentityUser struct {
	db uleveldb.ULevelDB
}

const dagPoolUser = "dagPoolUser/"

var (
	OnlyRead  dagPoolPolicy = "only-read"
	OnlyWrite dagPoolPolicy = "only-write"
	ReadWrite dagPoolPolicy = "read-write"
)

type DagPoolUser struct {
	username string
	password string
	policy   dagPoolPolicy
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
