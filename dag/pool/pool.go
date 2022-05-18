package pool

import (
	"context"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/config"
	"github.com/filedag-project/filedag-storage/dag/node"
	"github.com/filedag-project/filedag-storage/dag/pool/dagpooluser"
	"github.com/filedag-project/filedag-storage/dag/pool/referencecount"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	blocks "github.com/ipfs/go-block-format"
	bserv "github.com/ipfs/go-blockservice"
	cid "github.com/ipfs/go-cid"
	format "github.com/ipfs/go-ipld-format"
	legacy "github.com/ipfs/go-ipld-legacy"
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
	CidBuilder       cid.Builder
	ImporterBatchNum int
	NRSys            NodeRecordSys
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
	i, err := dagpooluser.NewIdentityUserSys(db)
	if err != nil {
		return nil, err
	}
	r, err := referencecount.NewIdentityRefe(db)
	dn := make(map[string]*node.DagNode)
	var nrs = NewRecordSys(db)
	for num, c := range cfg.DagNodeConfig {
		bs, err := node.NewDagNode(c)
		if err != nil {
			log.Errorf("new dagnode err:%v", err)
			return nil, err
		}
		name := "the" + fmt.Sprintf("%v", num)
		err = nrs.HandleDagNode(c.Nodes, name)
		if err != nil {
			return nil, err
		}
		dn[name] = bs
	}
	return &DagPool{DagNodes: dn, Iam: i, refer: r, CidBuilder: cidBuilder, ImporterBatchNum: cfg.ImporterBatchNum, NRSys: nrs}, nil
}

// Add adds a node to the DagPool, storing the block in the BlockService
func (d *DagPool) Add(ctx context.Context, nd format.Node) error {
	if d == nil { // FIXME remove this assertion. protect with constructor invariant
		return fmt.Errorf("DagPool is nil")
	}
	err := d.refer.AddReference(nd.Cid().String())
	if err != nil {
		return err
	}
	return d.UseNode(ctx, nd.Cid()).Put(nd)
}

func (d *DagPool) AddMany(ctx context.Context, nds []format.Node) error {
	blks := make([]blocks.Block, len(nds))
	var cids []cid.Cid
	for i, nd := range nds {
		blks[i] = nd
		err := d.refer.AddReference(nd.Cid().String())
		if err != nil {
			return err
		}
		cids = append(cids, nd.Cid())
	}
	return d.UseNodes(ctx, cids).PutMany(blks)
}

// Get retrieves a node from the DagPool, fetching the block in the BlockService
func (d *DagPool) Get(ctx context.Context, c cid.Cid) (format.Node, error) {
	if d == nil {
		return nil, fmt.Errorf("DagPool is nil")
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	b, err := d.GetNode(ctx, c).Get(c)
	if err != nil {
		if err == bserv.ErrNotFound {
			return nil, format.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get block for %s: %v", c, err)
	}

	return legacy.DecodeNode(ctx, b)
}

// GetLinks return the links for the node, the node doesn't necessarily have
// to exist locally.
func (d *DagPool) GetLinks(ctx context.Context, c cid.Cid) ([]*format.Link, error) {
	if c.Type() == cid.Raw {
		return nil, nil
	}
	node, err := d.Get(ctx, c)
	if err != nil {
		return nil, err
	}
	return node.Links(), nil
}

func (d *DagPool) Remove(ctx context.Context, c cid.Cid) error {
	if d == nil { // FIXME remove this assertion. protect with constructor invariant
		return fmt.Errorf("DagPool is nil")
	}
	err := d.refer.RemoveReference(c.String())
	if err != nil {
		return err
	}
	return d.GetNode(ctx, c).DeleteBlock(c)
}

// RemoveMany removes multiple nodes from the DAG. It will likely be faster than
// removing them individually.
//
// This operation is not atomic. If it returns an error, some nodes may or may
// not have been removed.
func (d *DagPool) RemoveMany(ctx context.Context, cids []cid.Cid) error {
	// TODO(#4608): make this batch all the way down.
	for _, c := range cids {
		if err := d.GetNode(ctx, c).DeleteBlock(c); err != nil {
			return err
		}
		err := d.refer.RemoveReference(c.String())
		if err != nil {
			return err
		}
	}
	return nil
}

// GetMany gets many nodes from the DAG at once.
//
// This method may not return all requested nodes (and may or may not return an
// error indicating that it failed to do so. It is up to the caller to verify
// that it received all nodes.
//func (d *DagPool) GetMany(ctx context.Context, keys []cid.Cid) <-chan *format.NodeOption {
//	m := d.GetNodes(ctx, keys)
//	var a <-chan *format.NodeOption
//	for _, b := range d.DagNodes {
//		a = getNodesFromBG(ctx, b, m[b])
//	}
//	return a
//}

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
	if d == nil { // FIXME remove this assertion. protect with constructor invariant
		return fmt.Errorf("DagPool is nil")
	}
	dagNode, err := d.GetNodeUseIP(ctx, ip)
	if err != nil {
		return err
	}
	return dagNode.RepairDisk(ip, port)
}

// GetLinks is the type of function passed to the EnumerateChildren function(s)
// for getting the children of an IPLD node.
//type GetLinks func(context.Context, cid.Cid) ([]*format.Link, error)
//
//var _ format.LinkGetter = &DagPool{}
//var _ format.NodeGetter = &DagPool{}
//var _ format.DAGService = &DagPool{}
