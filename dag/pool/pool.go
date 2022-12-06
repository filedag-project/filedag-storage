package pool

import (
	"context"
	"github.com/filedag-project/filedag-storage/dag/config"
	"github.com/filedag-project/filedag-storage/dag/pool/poolservice/dpuser"
	"github.com/filedag-project/filedag-storage/dag/proto"
	"github.com/filedag-project/filedag-storage/dag/slotsmgr"
	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	// blank import is used to register the IPLD raw codec
	_ "github.com/ipld/go-ipld-prime/codec/raw"
)

//go:generate go run github.com/golang/mock/mockgen -destination=mocks/mock_dagpool.go -package=mocks . DagPool

// DagPool is an interface that defines the basic operations of a dag pool
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

// Cluster is an interface that defines the basic operations of a Cluster
type Cluster interface {
	AddDagNode(nodeConfig *config.DagNodeConfig) error
	GetDagNode(dagNodeName string) (*config.DagNodeConfig, error)
	RemoveDagNode(dagNodeName string) (*config.DagNodeConfig, error)
	MigrateSlots(fromDagNodeName, toDagNodeName string, pairs []slotsmgr.SlotPair) error
	BalanceSlots() error
	Status() (*proto.StatusReply, error)
	RepairDataNode(ctx context.Context, dagNodeName string, fromNodeIndex int, repairNodeIndex int) error
}
