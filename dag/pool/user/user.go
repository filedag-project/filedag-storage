package user

type dagPoolPolicy string

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
