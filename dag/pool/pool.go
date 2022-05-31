package pool

import (
	"context"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/config"
	"github.com/filedag-project/filedag-storage/dag/node"
	"github.com/filedag-project/filedag-storage/dag/pool/dagpooluser"
	dnm "github.com/filedag-project/filedag-storage/dag/pool/datanodemanager"
	"github.com/filedag-project/filedag-storage/dag/pool/referencecount"
	"github.com/filedag-project/filedag-storage/dag/pool/userpolicy"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	blocks "github.com/ipfs/go-block-format"
	bserv "github.com/ipfs/go-blockservice"
	cid "github.com/ipfs/go-cid"
	format "github.com/ipfs/go-ipld-format"
	logging "github.com/ipfs/go-log/v2"
	"github.com/ipfs/go-merkledag"
	// blank import is used to register the IPLD raw codec
	_ "github.com/ipld/go-ipld-prime/codec/raw"
)

var log = logging.Logger("dag-pool")

type DagPool interface {
	Add(ctx context.Context, block blocks.Block) error
	Get(ctx context.Context, c cid.Cid) (blocks.Block, error)
	Remove(ctx context.Context, c cid.Cid) error
	DataRepairHost(ctx context.Context, oldIp, newIp, oldPort, newPort string) error
	DataRepairDisk(ctx context.Context, ip, port string) error
	CheckUserPolicy(username, pass string, policy userpolicy.DagPoolPolicy) bool
	CheckDeal(user, pass string) bool
	AddUser(user dagpooluser.DagPoolUser) error
	RemoveUser(username string) error
	QueryUser(username string) (dagpooluser.DagPoolUser, error)
	UpdateUser(u dagpooluser.DagPoolUser) error
}

// Pool is an IPFS Merkle DAG service.
type Pool struct {
	DagNodes         map[string]*node.DagNode
	Iam              dagpooluser.IdentityUserSys
	refer            referencecount.IdentityRefe
	CidBuilder       cid.Builder
	ImporterBatchNum int
	NRSys            dnm.NodeRecordSys
}

// NewDagPoolService constructs a new DAGService (using the default implementation).
// Note that the default implementation is also an ipld.LinkGetter.
func NewDagPoolService(cfg config.PoolConfig) (*Pool, error) {
	cidBuilder, err := merkledag.PrefixForCidVersion(0)
	if err != nil {
		return nil, err
	}
	db, err := uleveldb.OpenDb(cfg.LeveldbPath)
	if err != nil {
		return nil, err
	}
	i, err := dagpooluser.NewIdentityUserSys(db, cfg.DefaultUser, cfg.DefaultPass)
	if err != nil {
		return nil, err
	}
	r, err := referencecount.NewIdentityRefe(db)
	dn := make(map[string]*node.DagNode)
	var nrs = dnm.NewRecordSys(db)
	for num, c := range cfg.DagNodeConfig {
		bs, err := node.NewDagNode(c)
		if err != nil {
			log.Errorf("new dagnode err:%v", err)
			return nil, err
		}
		name := "the" + fmt.Sprintf("%v", num)
		err = nrs.HandleDagNode(bs.Nodes, name)
		if err != nil {
			return nil, err
		}
		dn[name] = bs
	}
	return &Pool{DagNodes: dn, Iam: i, refer: r, CidBuilder: cidBuilder, ImporterBatchNum: cfg.ImporterBatchNum, NRSys: nrs}, nil
}

// Add adds a node to the Pool, storing the block in the BlockService
func (d *Pool) Add(ctx context.Context, block blocks.Block) error {
	if d == nil { // FIXME remove this assertion. protect with constructor invariant
		return fmt.Errorf("Pool is nil")
	}
	err := d.refer.AddReference(block.Cid().String())
	if err != nil {
		return err
	}
	reference, err := d.refer.QueryReference(block.Cid().String())
	if err != nil {
		return err
	}
	if reference > 1 {
		return nil
	}
	useNode, err := d.UseNode(ctx, block.Cid())
	if err != nil {
		return err
	}
	return useNode.Put(block)
}

// Get retrieves a node from the Pool, fetching the block in the BlockService
func (d *Pool) Get(ctx context.Context, c cid.Cid) (blocks.Block, error) {
	if d == nil {
		return nil, fmt.Errorf("Pool is nil")
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	reference, err := d.refer.QueryReference(c.String())
	if err != nil {
		return nil, err
	}
	if reference <= 0 {
		return nil, fmt.Errorf("block does not exist : %v", err)
	}
	getNode, err := d.GetNode(ctx, c)
	if err != nil {
		return nil, err
	}
	b, err := getNode.Get(c)
	if err != nil {
		if err == bserv.ErrNotFound {
			return nil, format.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get block for %s: %v", c, err)
	}

	return b, nil
}

func (d *Pool) Remove(ctx context.Context, c cid.Cid) error {
	if d == nil { // FIXME remove this assertion. protect with constructor invariant
		return fmt.Errorf("Pool is nil")
	}
	err := d.refer.RemoveReference(c.String())
	if err != nil {
		return err
	}
	reference, err := d.refer.QueryReference(c.String())
	if err != nil {
		return err
	}
	if reference == 0 {
		getNode, err := d.GetNode(ctx, c)
		if err != nil {
			return err
		}
		go getNode.DeleteBlock(c)
	}
	return nil
}

// DataRepairHost Data repair host
func (d *Pool) DataRepairHost(ctx context.Context, oldIp, newIp, oldPort, newPort string) error {
	if d == nil {
		return fmt.Errorf("Pool is nil")
	}
	dagNode, err := d.GetNodeUseIP(ctx, oldIp)
	if err != nil {
		return err
	}
	return dagNode.RepairHost(oldIp, newIp, oldPort, newPort)
}

// DataRepairDisk Data repair disk
func (d *Pool) DataRepairDisk(ctx context.Context, ip, port string) error {
	if d == nil {
		return fmt.Errorf("Pool is nil")
	}
	dagNode, err := d.GetNodeUseIP(ctx, ip)
	if err != nil {
		return err
	}
	return dagNode.RepairDisk(ip, port)
}

func (d *Pool) CheckDeal(user, pass string) bool {
	return d.Iam.CheckDeal(user, pass)
}

func (d *Pool) AddUser(user dagpooluser.DagPoolUser) error {
	return d.Iam.AddUser(user)
}

func (d *Pool) RemoveUser(username string) error {
	return d.Iam.RemoveUser(username)
}

func (d *Pool) QueryUser(username string) (dagpooluser.DagPoolUser, error) {
	return d.Iam.QueryUser(username)
}

func (d *Pool) UpdateUser(u dagpooluser.DagPoolUser) error {
	return d.Iam.UpdateUser(u)
}

func (d *Pool) CheckUserPolicy(username, pass string, policy userpolicy.DagPoolPolicy) bool {
	return d.Iam.CheckUserPolicy(username, pass, policy)
}

var _ DagPool = &Pool{}
