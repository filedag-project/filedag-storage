package poolservice

import (
	"context"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/config"
	"github.com/filedag-project/filedag-storage/dag/node/dagnode"
	"github.com/filedag-project/filedag-storage/dag/pool"
	"github.com/filedag-project/filedag-storage/dag/pool/poolservice/dnm"
	"github.com/filedag-project/filedag-storage/dag/pool/poolservice/dpuser"
	"github.com/filedag-project/filedag-storage/dag/pool/poolservice/dpuser/upolicy"
	"github.com/filedag-project/filedag-storage/dag/pool/poolservice/refSys"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	format "github.com/ipfs/go-ipld-format"
	ipldlegacy "github.com/ipfs/go-ipld-legacy"
	logging "github.com/ipfs/go-log/v2"
	"golang.org/x/xerrors"
)

var log = logging.Logger("dag-pool")

var _ pool.DagPool = &dagPoolService{}

// dagPoolService is an IPFS Merkle DAG service.
type dagPoolService struct {
	dagNodes map[string]*dagnode.DagNode
	iam      *dpuser.IdentityUserSys
	refer    *refSys.ReferSys
	nrSys    *dnm.NodeRecordSys
	gc       *gc
	db       *uleveldb.ULevelDB
}

func (d *dagPoolService) NeedPin(username string) bool {
	//todo more check
	return d.iam.IsAdmin(username)
}

// NewDagPoolService constructs a new DAGPool (using the default implementation).
func NewDagPoolService(cfg config.PoolConfig) (*dagPoolService, error) {
	db, err := uleveldb.OpenDb(cfg.LeveldbPath)
	if err != nil {
		return nil, err
	}
	i, err := dpuser.NewIdentityUserSys(db, cfg.RootUser, cfg.RootPassword)
	if err != nil {
		return nil, err
	}
	r := refSys.NewReferSys(db, cfg.CacheExpireTime)
	dn := make(map[string]*dagnode.DagNode)
	var nrs = dnm.NewRecordSys(db)
	for num, c := range cfg.DagNodeConfig {
		bs, err := dagnode.NewDagNode(c)
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
		dagNodes: dn,
		iam:      i,
		refer:    r,
		nrSys:    nrs,
		gc: &gc{
			stopCacheCh: make(chan struct{}),
			stopStoreCh: make(chan struct{}),
			normalCh:    make(chan struct{}),
			gcPeriod:    cfg.GcPeriod,
		},
		db: db,
	}, nil
}

// Add adds a node to the dagPoolService, storing the block in the BlockService
func (d *dagPoolService) Add(ctx context.Context, block blocks.Block, user string, password string) error {
	if !d.iam.CheckUserPolicy(user, password, upolicy.WriteOnly) {
		return upolicy.AccessDenied
	}
	d.gc.Stop()
	if !d.refer.HasReference(block.Cid().String()) {
		useNode, err := d.choseDagNode(ctx, block.Cid())
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
	if !d.iam.CheckUserPolicy(user, password, upolicy.ReadOnly) {
		return nil, upolicy.AccessDenied
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	if !d.refer.HasReference(c.String()) {
		return nil, format.ErrNotFound{Cid: c}
	}
	getNode, err := d.getDagNodeInfo(ctx, c)
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
	if !d.iam.CheckUserPolicy(user, password, upolicy.WriteOnly) {
		return upolicy.AccessDenied
	}
	d.gc.Stop()
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

//GetSize get the block size
func (d *dagPoolService) GetSize(ctx context.Context, c cid.Cid, user string, password string) (int, error) {
	if !d.iam.CheckUserPolicy(user, password, upolicy.ReadOnly) {
		return 0, upolicy.AccessDenied
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	has := d.refer.HasReference(c.String())
	if !has {
		return 0, format.ErrNotFound{Cid: c}
	}
	getNode, err := d.getDagNodeInfo(ctx, c)
	if err != nil {
		return 0, err
	}
	return getNode.GetSize(ctx, c)
}

//AddUser add a user
func (d *dagPoolService) AddUser(newUser dpuser.DagPoolUser, user string, password string) error {
	if !d.iam.CheckAdmin(user, password) {
		return upolicy.AccessDenied
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
		return upolicy.AccessDenied
	}
	if d.iam.IsAdmin(rmUser) {
		return xerrors.New("refuse to remove the admin user")
	}
	return d.iam.RemoveUser(rmUser)
}

//QueryUser query the user
func (d *dagPoolService) QueryUser(qUser string, user string, password string) (*dpuser.DagPoolUser, error) {
	if !d.iam.CheckUser(user, password) {
		return nil, upolicy.AccessDenied
	}
	if d.iam.IsAdmin(user) {
		return d.iam.QueryUser(qUser)
	}
	// only query self config
	if qUser != user {
		return nil, upolicy.AccessDenied
	}
	return d.iam.QueryUser(qUser)
}

//UpdateUser update the user
func (d *dagPoolService) UpdateUser(uUser dpuser.DagPoolUser, user string, password string) error {
	if !d.iam.CheckAdmin(user, password) {
		return upolicy.AccessDenied
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
		return nil, format.ErrNotFound{Cid: ci}
	}
	getNode, err := d.getDagNodeInfo(ctx, ci)
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
