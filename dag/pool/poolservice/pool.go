package poolservice

import (
	"context"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/config"
	"github.com/filedag-project/filedag-storage/dag/node"
	"github.com/filedag-project/filedag-storage/dag/pool"
	"github.com/filedag-project/filedag-storage/dag/pool/dagpooluser"
	dnm "github.com/filedag-project/filedag-storage/dag/pool/datanodemanager"
	"github.com/filedag-project/filedag-storage/dag/pool/referencecount"
	"github.com/filedag-project/filedag-storage/dag/pool/userpolicy"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	format "github.com/ipfs/go-ipld-format"
	ipldlegacy "github.com/ipfs/go-ipld-legacy"
	logging "github.com/ipfs/go-log/v2"
	"github.com/ipfs/go-merkledag"
	"golang.org/x/xerrors"
)

var log = logging.Logger("dag-pool")

var _ pool.DagPool = &dagPoolService{}

// dagPoolService is an IPFS Merkle DAG service.
type dagPoolService struct {
	dagNodes   map[string]*node.DagNode
	iam        *dagpooluser.IdentityUserSys
	refer      *referencecount.ReferSys
	cidBuilder cid.Builder
	nrSys      *dnm.NodeRecordSys
	gc         *gc
	db         *uleveldb.ULevelDB
}

func (d *dagPoolService) NeedPin(username string) bool {
	//todo more check
	return d.iam.IsAdmin(username)
}

// NewDagPoolService constructs a new DAGPool (using the default implementation).
func NewDagPoolService(cfg config.PoolConfig) (*dagPoolService, error) {
	cidBuilder, err := merkledag.PrefixForCidVersion(0)
	if err != nil {
		return nil, err
	}
	db, err := uleveldb.OpenDb(cfg.LeveldbPath)
	if err != nil {
		return nil, err
	}
	i, err := dagpooluser.NewIdentityUserSys(db, cfg.RootUser, cfg.RootPassword)
	if err != nil {
		return nil, err
	}
	r := referencecount.NewIdentityRefe(db)
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
		cidBuilder: cidBuilder,
		nrSys:      nrs,
		gc: &gc{
			stopCh: make(chan struct{}),
		},
		db: db,
	}, nil
}

// Add adds a node to the dagPoolService, storing the block in the BlockService
func (d *dagPoolService) Add(ctx context.Context, block blocks.Block, user string, password string) error {
	if !d.iam.CheckUserPolicy(user, password, userpolicy.OnlyWrite) {
		return userpolicy.AccessDenied
	}
	d.gc.Stop()
	if !d.refer.HasReference(block.Cid().String()) {
		useNode, err := d.UseNode(ctx, block.Cid())
		if err != nil {
			return err
		}
		err = useNode.Put(ctx, block)
		if err != nil {
			return err
		}
	}
	err := d.refer.AddReference(block.Cid().String(), d.NeedPin(user))
	if err != nil {
		return err
	}
	return nil
}

// Get retrieves a node from the dagPoolService, fetching the block in the BlockService
func (d *dagPoolService) Get(ctx context.Context, c cid.Cid, user string, password string) (blocks.Block, error) {
	if !d.iam.CheckUserPolicy(user, password, userpolicy.OnlyRead) {
		return nil, userpolicy.AccessDenied
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	if !d.refer.HasReference(c.String()) {
		return nil, fmt.Errorf("block:%v does not exist", c.String())
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

//Remove remove block from DAGPool
func (d *dagPoolService) Remove(ctx context.Context, c cid.Cid, user string, password string) error {
	if !d.iam.CheckUserPolicy(user, password, userpolicy.OnlyWrite) {
		return userpolicy.AccessDenied
	}
	if d.refer.HasReference(c.String()) {
		err := d.refer.RemoveReference(c.String(), d.NeedPin(user))
		if err != nil {
			return err
		}
	}
	//if reference == 0 {
	//	getNode, err := d.GetNode(ctx, c)
	//	if err != nil {
	//		return err
	//	}
	//	go getNode.DeleteBlock(ctx, c)
	//}
	return nil
}

// DataRepairHost Data repair host
func (d *dagPoolService) DataRepairHost(ctx context.Context, oldIp, newIp, oldPort, newPort string) error {
	if d == nil {
		return fmt.Errorf("Pool is nil")
	}
	dagNode, err := d.getNodeUseIP(ctx, oldIp)
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
	dagNode, err := d.getNodeUseIP(ctx, ip)
	if err != nil {
		return err
	}
	return dagNode.RepairDisk(ip, port)
}

//AddUser add a user
func (d *dagPoolService) AddUser(newUser dagpooluser.DagPoolUser, user string, password string) error {
	if !d.iam.CheckAdmin(user, password) {
		return userpolicy.AccessDenied
	}
	if d.iam.IsAdmin(newUser.Username) {
		return xerrors.New("the user already exists")
	}
	if _, err := d.iam.QueryUser(newUser.Username); err == nil {
		return xerrors.New("the user already exists")
	}
	return d.iam.AddUser(newUser)
}

//RemoveUser remove the user
func (d *dagPoolService) RemoveUser(rmUser string, user string, password string) error {
	if !d.iam.CheckAdmin(user, password) {
		return userpolicy.AccessDenied
	}
	if d.iam.IsAdmin(rmUser) {
		return xerrors.New("refuse to remove the admin user")
	}
	return d.iam.RemoveUser(rmUser)
}

//QueryUser query the user
func (d *dagPoolService) QueryUser(qUser string, user string, password string) (*dagpooluser.DagPoolUser, error) {
	if !d.iam.CheckUser(user, password) {
		return nil, userpolicy.AccessDenied
	}
	if d.iam.IsAdmin(user) {
		return d.iam.QueryUser(qUser)
	}
	// only query self config
	if qUser != user {
		return nil, userpolicy.AccessDenied
	}
	return d.iam.QueryUser(qUser)
}

//UpdateUser update the user
func (d *dagPoolService) UpdateUser(uUser dagpooluser.DagPoolUser, user string, password string) error {
	if !d.iam.CheckAdmin(user, password) {
		return userpolicy.AccessDenied
	}
	if d.iam.IsAdmin(uUser.Username) {
		return xerrors.New("refuse to update the admin user")
	}
	u, err := d.iam.QueryUser(uUser.Username)
	if err != nil {
		return xerrors.New("not found the user")
	}
	if uUser.Password != "" {
		u.Password = uUser.Password
	}
	if uUser.Policy != "" {
		u.Policy = uUser.Policy
	}
	if uUser.Capacity != 0 {
		u.Capacity = uUser.Capacity
	}
	return d.iam.UpdateUser(*u)
}

//func (d *dagPoolService) CheckUserPolicy(username, pass string, policy userpolicy.DagPoolPolicy) bool {
//	return d.iam.CheckUserPolicy(username, pass, policy)
//}

//Close the dagPoolService
func (d *dagPoolService) Close() error {
	return d.db.Close()
}

//GetLinks get links from DAGPool
func (d *dagPoolService) GetLinks(ctx context.Context, ci cid.Cid) ([]*format.Link, error) {
	if d == nil {
		return nil, fmt.Errorf("Pool is nil")
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	if !d.refer.HasReference(ci.String()) {
		return nil, fmt.Errorf("block : %v does not exist ", ci.String())
	}
	getNode, err := d.GetNode(ctx, ci)
	if err != nil {
		return nil, err
	}
	b, err := getNode.Get(ctx, ci)
	if err != nil {
		if format.IsNotFound(err) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get block for %s: %v", ci, err)
	}
	decodeNode, err := ipldlegacy.DecodeNode(ctx, b)
	if err != nil {
		return nil, err
	}
	return decodeNode.Links(), nil
}
