package pool

import (
	"context"
	"github.com/filedag-project/filedag-storage/dag/pool/dagpooluser"
	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	// blank import is used to register the IPLD raw codec
	_ "github.com/ipld/go-ipld-prime/codec/raw"
)

//DagPool define dagpool interface
type DagPool interface {
	Add(ctx context.Context, block blocks.Block, user string, password string) error
	Get(ctx context.Context, c cid.Cid, user string, password string) (blocks.Block, error)
	GetSize(ctx context.Context, c cid.Cid, user string, password string) (int, error)
	Remove(ctx context.Context, c cid.Cid, user string, password string) error
	DataRepairHost(ctx context.Context, oldIp, newIp, oldPort, newPort string) error
	DataRepairDisk(ctx context.Context, ip, port string) error
	AddUser(newUser dagpooluser.DagPoolUser, user string, password string) error
	RemoveUser(rmUser string, user string, password string) error
	QueryUser(qUser string, user string, password string) (*dagpooluser.DagPoolUser, error)
	UpdateUser(uUser dagpooluser.DagPoolUser, user string, password string) error
	Close() error
	UnPin(ctx context.Context, c cid.Cid, user string, password string) error
	Pin(ctx context.Context, c cid.Cid, user string, password string) error
	NeedPin(username string) bool
}
