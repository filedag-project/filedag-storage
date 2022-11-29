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
	"github.com/filedag-project/filedag-storage/dag/pool/poolservice/reference"
	"github.com/filedag-project/filedag-storage/objectservice/uleveldb"
	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	format "github.com/ipfs/go-ipld-format"
	logging "github.com/ipfs/go-log/v2"
	"golang.org/x/xerrors"
	"time"
)

var log = logging.Logger("dag-pool")

var _ pool.DagPool = &dagPoolService{}

// dagPoolService is an IPFS Merkle DAG service.
type dagPoolService struct {
	dagNodes map[string]*dagnode.DagNode
	iam      *dpuser.IdentityUserSys
	nrSys    *dnm.NodeRecordSys
	db       *uleveldb.ULevelDB

	refCounter *reference.RefCounter
	cacheSet   *reference.CacheSet

	gcControl *GcControl
	gcPeriod  time.Duration
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
	cacheSet := reference.NewCacheSet(db)
	refCounter := reference.NewRefCounter(db, cacheSet)
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
		dagNodes:   dn,
		iam:        i,
		nrSys:      nrs,
		db:         db,
		refCounter: refCounter,
		cacheSet:   cacheSet,
		gcControl:  NewGcControl(),
		gcPeriod:   cfg.GcPeriod,
	}, nil
}

// Add adds a node to the dagPoolService, storing the block in the BlockService
func (d *dagPoolService) Add(ctx context.Context, block blocks.Block, user string, password string, pin bool) error {
	if !d.iam.CheckUserPolicy(user, password, upolicy.WriteOnly) {
		return upolicy.AccessDenied
	}

	key := block.Cid().String()
	addBlock := func() error {
		selNode, err := d.choseDagNode(ctx, block.Cid())
		if err != nil {
			return err
		}
		return selNode.Put(ctx, block)
	}

	if pin {
		d.InterruptGC()
		return d.refCounter.IncrOrCreate(key, addBlock)
	}

	if has, _ := d.Has(key); !has {
		if err := addBlock(); err != nil {
			return err
		}
		return d.cacheSet.Add(key)
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

	key := c.String()
	if has, err := d.Has(key); err != nil {
		return nil, err
	} else if !has {
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
func (d *dagPoolService) Remove(ctx context.Context, c cid.Cid, user string, password string, unpin bool) error {
	if !d.iam.CheckUserPolicy(user, password, upolicy.WriteOnly) {
		return upolicy.AccessDenied
	}

	if unpin {
		return d.refCounter.Decr(c.String())
	}
	return nil
}

//GetSize get the block size
func (d *dagPoolService) GetSize(ctx context.Context, c cid.Cid, user string, password string) (int, error) {
	if !d.iam.CheckUserPolicy(user, password, upolicy.ReadOnly) {
		return 0, upolicy.AccessDenied
	}

	key := c.String()
	if has, err := d.Has(key); err != nil {
		return 0, err
	} else if !has {
		return 0, format.ErrNotFound{Cid: c}
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
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

func (d *dagPoolService) Has(key string) (bool, error) {
	// is pinned?
	if has, err := d.refCounter.Has(key); err != nil {
		return false, err
	} else if has {
		return true, nil
	}
	// is cached?
	return d.cacheSet.Has(key)
}
