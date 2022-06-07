package poolservice

import (
	"context"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/config"
	"github.com/filedag-project/filedag-storage/dag/node"
	"github.com/filedag-project/filedag-storage/dag/pool"
	"github.com/filedag-project/filedag-storage/dag/pool/blockpinner"
	"github.com/filedag-project/filedag-storage/dag/pool/dagpooluser"
	dnm "github.com/filedag-project/filedag-storage/dag/pool/datanodemanager"
	leveldbds "github.com/filedag-project/filedag-storage/dag/pool/leveldb-datastore"
	"github.com/filedag-project/filedag-storage/dag/pool/referencecount"
	"github.com/filedag-project/filedag-storage/dag/pool/userpolicy"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	format "github.com/ipfs/go-ipld-format"
	ipldlegacy "github.com/ipfs/go-ipld-legacy"
	logging "github.com/ipfs/go-log/v2"
	"github.com/ipfs/go-merkledag"
)

var log = logging.Logger("dag-pool")

var _ pool.DagPool = &dagPoolService{}

// dagPoolService is an IPFS Merkle DAG service.
type dagPoolService struct {
	dagNodes map[string]*node.DagNode
	iam      dagpooluser.IdentityUserSys
	refer    referencecount.IdentityRefe
	//Pining     datapin.PinService
	pinner     *blockpinner.Pinner
	cidBuilder cid.Builder
	NRSys      dnm.NodeRecordSys
	db         *uleveldb.ULevelDB
}

// NewDagPoolService constructs a new DAGService (using the default implementation).
// Note that the default implementation is also an ipld.LinkGetter.
func NewDagPoolService(cfg config.PoolConfig) (*dagPoolService, error) {
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
	ldstore, err := leveldbds.NewDatastore(cfg.DatastorePath, nil)
	if err != nil {
		panic(err)
	}
	pn, _ := blockpinner.New(context.TODO(), ldstore)
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
	return &dagPoolService{
		dagNodes:   dn,
		iam:        i,
		refer:      r,
		pinner:     pn,
		cidBuilder: cidBuilder,
		NRSys:      nrs,
		db:         db,
	}, nil
}

// Add adds a node to the dagPoolService, storing the block in the BlockService
func (d *dagPoolService) Add(ctx context.Context, block blocks.Block) error {
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
	return useNode.Put(ctx, block)
}

// Get retrieves a node from the dagPoolService, fetching the block in the BlockService
func (d *dagPoolService) Get(ctx context.Context, c cid.Cid) (blocks.Block, error) {
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
	b, err := getNode.Get(ctx, c)
	if err != nil {
		if format.IsNotFound(err) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get block for %s: %v", c, err)
	}

	return b, nil
}

func (d *dagPoolService) Remove(ctx context.Context, c cid.Cid) error {
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
		go getNode.DeleteBlock(ctx, c)
	}
	return nil
}

// DataRepairHost Data repair host
func (d *dagPoolService) DataRepairHost(ctx context.Context, oldIp, newIp, oldPort, newPort string) error {
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
func (d *dagPoolService) DataRepairDisk(ctx context.Context, ip, port string) error {
	if d == nil {
		return fmt.Errorf("Pool is nil")
	}
	dagNode, err := d.GetNodeUseIP(ctx, ip)
	if err != nil {
		return err
	}
	return dagNode.RepairDisk(ip, port)
}

func (d *dagPoolService) CheckDeal(user, pass string) bool {
	return d.iam.CheckDeal(user, pass)
}

func (d *dagPoolService) AddUser(user dagpooluser.DagPoolUser) error {
	return d.iam.AddUser(user)
}

func (d *dagPoolService) RemoveUser(username string) error {
	return d.iam.RemoveUser(username)
}

func (d *dagPoolService) QueryUser(username string) (dagpooluser.DagPoolUser, error) {
	return d.iam.QueryUser(username)
}

func (d *dagPoolService) UpdateUser(u dagpooluser.DagPoolUser) error {
	return d.iam.UpdateUser(u)
}

func (d *dagPoolService) CheckUserPolicy(username, pass string, policy userpolicy.DagPoolPolicy) bool {
	return d.iam.CheckUserPolicy(username, pass, policy)
}

func (d *dagPoolService) Close() error {
	return d.db.Close()
}

func (d *dagPoolService) GetLinks(ctx context.Context, ci cid.Cid) ([]*format.Link, error) {
	get, err := d.Get(ctx, ci)
	if err != nil {
		return nil, err
	}
	decodeNode, err := ipldlegacy.DecodeNode(ctx, get)
	if err != nil {
		return nil, err
	}
	return decodeNode.Links(), nil
}
