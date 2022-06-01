package pool

import (
	"context"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/config"
	"github.com/filedag-project/filedag-storage/dag/node"
	"github.com/filedag-project/filedag-storage/dag/pool/dagpooluser"
	dnm "github.com/filedag-project/filedag-storage/dag/pool/datanodemanager"
	"github.com/filedag-project/filedag-storage/dag/pool/datapin"
	"github.com/filedag-project/filedag-storage/dag/pool/referencecount"
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

// DagPool is an IPFS Merkle DAG service.
// - the root is virtual (like a forest)
// - stores nodes' data in a BlockService
// TODO: should cache Nodes that are in memory, and be
//       able to free some of them when vm pressure is high
type DagPool struct {
	DagNodes         map[string]*node.DagNode
	Iam              dagpooluser.IdentityUserSys
	refer            referencecount.IdentityRefe
	Pin              datapin.PinService
	CidBuilder       cid.Builder
	ImporterBatchNum int
	NRSys            dnm.NodeRecordSys
}

// NewDagPoolService constructs a new DAGService (using the default implementation).
// Note that the default implementation is also an ipld.LinkGetter.
func NewDagPoolService(cfg config.PoolConfig) (*DagPool, error) {
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
	return &DagPool{DagNodes: dn, Iam: i, refer: r, CidBuilder: cidBuilder, ImporterBatchNum: cfg.ImporterBatchNum, NRSys: nrs}, nil
}

// Add adds a node to the DagPool, storing the block in the BlockService
func (d *DagPool) Add(ctx context.Context, block blocks.Block) error {
	if d == nil { // FIXME remove this assertion. protect with constructor invariant
		return fmt.Errorf("DagPool is nil")
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

// Get retrieves a node from the DagPool, fetching the block in the BlockService
func (d *DagPool) Get(ctx context.Context, c cid.Cid) (blocks.Block, error) {
	if d == nil {
		return nil, fmt.Errorf("DagPool is nil")
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

func (d *DagPool) Remove(ctx context.Context, c cid.Cid) error {
	if d == nil { // FIXME remove this assertion. protect with constructor invariant
		return fmt.Errorf("DagPool is nil")
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
func (d *DagPool) DataRepairHost(ctx context.Context, oldIp, newIp, oldPort, newPort string) error {
	if d == nil {
		return fmt.Errorf("DagPool is nil")
	}
	dagNode, err := d.GetNodeUseIP(ctx, oldIp)
	if err != nil {
		return err
	}
	return dagNode.RepairHost(oldIp, newIp, oldPort, newPort)
}

// DataRepairDisk Data repair disk
func (d *DagPool) DataRepairDisk(ctx context.Context, ip, port string) error {
	if d == nil {
		return fmt.Errorf("DagPool is nil")
	}
	dagNode, err := d.GetNodeUseIP(ctx, ip)
	if err != nil {
		return err
	}
	return dagNode.RepairDisk(ip, port)
}
