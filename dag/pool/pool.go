package pool

import (
	"context"
	"github.com/filedag-project/filedag-storage/dag/pool/poolservice/dpuser"
	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	// blank import is used to register the IPLD raw codec
	_ "github.com/ipld/go-ipld-prime/codec/raw"
)

//go:generate go run github.com/golang/mock/mockgen -destination=mocks/mock_dagpool.go -package=mocks . DagPool

//DagPool is an interface that defines the basic operations of a dag pool
type DagPool interface {
	Add(ctx context.Context, block blocks.Block, user string, password string, pin bool) error
	Get(ctx context.Context, c cid.Cid, user string, password string) (blocks.Block, error)
	GetSize(ctx context.Context, c cid.Cid, user string, password string) (int, error)
	Remove(ctx context.Context, c cid.Cid, user string, password string, unpin bool) error
	AddUser(newUser dpuser.DagPoolUser, user string, password string) error
	RemoveUser(rmUser string, user string, password string) error
	QueryUser(qUser string, user string, password string) (*dpuser.DagPoolUser, error)
	UpdateUser(uUser dpuser.DagPoolUser, user string, password string) error
	Close() error
}
