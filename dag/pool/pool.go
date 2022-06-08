package pool

import (
	"context"
	"github.com/filedag-project/filedag-storage/dag/pool/dagpooluser"
	"github.com/filedag-project/filedag-storage/dag/pool/userpolicy"
	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	// blank import is used to register the IPLD raw codec
	_ "github.com/ipld/go-ipld-prime/codec/raw"
)

//DagPool define dagpool interface
type DagPool interface {
	Add(ctx context.Context, block blocks.Block, pin bool) error
	Get(ctx context.Context, c cid.Cid, pin bool) (blocks.Block, error)
	Remove(ctx context.Context, c cid.Cid, pin bool) error
	DataRepairHost(ctx context.Context, oldIp, newIp, oldPort, newPort string) error
	DataRepairDisk(ctx context.Context, ip, port string) error
	CheckUserPolicy(username, pass string, policy userpolicy.DagPoolPolicy) bool
	CheckDeal(user, pass string) bool
	AddUser(user dagpooluser.DagPoolUser) error
	RemoveUser(username string) error
	QueryUser(username string) (dagpooluser.DagPoolUser, error)
	UpdateUser(u dagpooluser.DagPoolUser) error
	Close() error
	UnPin(context.Context, cid.Cid) error
	Pin(context.Context, cid.Cid) error
	NeedPin(username string) bool
	//IsPinned(ctx context.Context, cid cid.Cid) bool
}
