package user

type dagPoolPolicy string
type IdentityUser struct {
}

var (
	OnlyRead  dagPoolPolicy = "only-read"
	OnlyWrite dagPoolPolicy = "only-write"
	ReadWrite dagPoolPolicy = "readWrite"
)

type DagPoolUser struct {
	username string
	password string
	policy   dagPoolPolicy
	capacity uint64
}

// AddUser add user
func (identity *IdentityUser) AddUser(user *DagPoolUser) error {
	//Todo
	var err error
	return err
}

// RemoveUser remove user
func (identity *IdentityUser) RemoveUser(userName string) error {
	//Todo
	var err error
	return err
}

// QueryUser query user
func (identity *IdentityUser) QueryUser(userName string) (*DagPoolUser, error) {
	//Todo
	var err error
	return &DagPoolUser{}, err
}
