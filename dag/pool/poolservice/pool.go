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
	logging "github.com/ipfs/go-log/v2"
	"github.com/ipfs/go-merkledag"
	"golang.org/x/xerrors"
)

var log = logging.Logger("dag-pool")

var _ pool.DagPool = &Pool{}

type Pool struct {
	DagNodes   map[string]*node.DagNode
	iam        *dagpooluser.IdentityUserSys
	refer      referencecount.IdentityRefe
	CidBuilder cid.Builder
	NRSys      dnm.NodeRecordSys
	db         *uleveldb.ULevelDB
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
	i, err := dagpooluser.NewIdentityUserSys(db, cfg.RootUser, cfg.RootPassword)
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
	return &Pool{
		DagNodes:   dn,
		iam:        i,
		refer:      r,
		CidBuilder: cidBuilder,
		NRSys:      nrs,
		db:         db,
	}, nil
}

// Add adds a node to the Pool, storing the block in the BlockService
func (d *Pool) Add(ctx context.Context, block blocks.Block, user string, password string) error {
	if !d.iam.CheckUserPolicy(user, password, userpolicy.OnlyWrite) {
		return userpolicy.AccessDenied
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

// Get retrieves a node from the Pool, fetching the block in the BlockService
func (d *Pool) Get(ctx context.Context, c cid.Cid, user string, password string) (blocks.Block, error) {
	if !d.iam.CheckUserPolicy(user, password, userpolicy.OnlyRead) {
		return nil, userpolicy.AccessDenied
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	reference, err := d.refer.QueryReference(c.String())
	if err != nil {
		return nil, err
	}
	if reference <= 0 {
		return nil, format.ErrNotFound{Cid: c}
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

func (d *Pool) GetSize(ctx context.Context, c cid.Cid, user string, password string) (int, error) {
	if !d.iam.CheckUserPolicy(user, password, userpolicy.OnlyRead) {
		return 0, userpolicy.AccessDenied
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	reference, err := d.refer.QueryReference(c.String())
	if err != nil {
		return 0, err
	}
	if reference <= 0 {
		return 0, format.ErrNotFound{Cid: c}
	}
	getNode, err := d.GetNode(ctx, c)
	if err != nil {
		return 0, err
	}
	return getNode.GetSize(ctx, c)
}

func (d *Pool) Remove(ctx context.Context, c cid.Cid, user string, password string) error {
	if !d.iam.CheckUserPolicy(user, password, userpolicy.OnlyWrite) {
		return userpolicy.AccessDenied
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

func (d *Pool) AddUser(newUser dagpooluser.DagPoolUser, user string, password string) error {
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

func (d *Pool) RemoveUser(rmUser string, user string, password string) error {
	if !d.iam.CheckAdmin(user, password) {
		return userpolicy.AccessDenied
	}
	if d.iam.IsAdmin(rmUser) {
		return xerrors.New("refuse to remove the admin user")
	}
	return d.iam.RemoveUser(rmUser)
}

func (d *Pool) QueryUser(qUser string, user string, password string) (*dagpooluser.DagPoolUser, error) {
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

func (d *Pool) UpdateUser(uUser dagpooluser.DagPoolUser, user string, password string) error {
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

func (d *Pool) Close() error {
	return d.db.Close()
}
