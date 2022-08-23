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
	db       *uleveldb.ULevelDB
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
	r, err := refSys.NewReferSys(db)
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
		db:       db,
	}, nil
}

// Add adds a node to the dagPoolService, storing the block in the BlockService
func (d *dagPoolService) Add(ctx context.Context, block blocks.Block, user string, password string) error {
	if !d.iam.CheckUserPolicy(user, password, upolicy.OnlyWrite) {
		return upolicy.AccessDenied
	}
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
	err := d.refer.AddReference(block.Cid().String())
	if err != nil {
		return err
	}
	return nil

}

// Get retrieves a node from the dagPoolService, fetching the block in the BlockService
func (d *dagPoolService) Get(ctx context.Context, c cid.Cid, user string, password string) (blocks.Block, error) {
	if !d.iam.CheckUserPolicy(user, password, upolicy.OnlyRead) {
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

//GetSize get the block size
func (d *dagPoolService) GetSize(ctx context.Context, c cid.Cid, user string, password string) (int, error) {
	if !d.iam.CheckUserPolicy(user, password, upolicy.OnlyRead) {
		return 0, upolicy.AccessDenied
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	if !d.refer.HasReference(c.String()) {
		return 0, format.ErrNotFound{Cid: c}
	}
	getNode, err := d.getDagNodeInfo(ctx, c)
	if err != nil {
		return 0, err
	}
	return getNode.GetSize(ctx, c)
}

//Remove remove block from DAGPool
func (d *dagPoolService) Remove(ctx context.Context, c cid.Cid, user string, password string) error {
	if !d.iam.CheckUserPolicy(user, password, upolicy.OnlyWrite) {
		return upolicy.AccessDenied
	}
	err := d.refer.RemoveReference(c.String())
	if err != nil {
		return err
	}
	//reference, err := d.refer.QueryReference(c.String())
	//if err != nil {
	//	return err
	//}
	//if reference == 0 {
	//	getNode, err := d.GetDagNodeInfo(ctx, c)
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

//Close the dagPoolService
func (d *dagPoolService) Close() error {
	for _, nd := range d.dagNodes {
		nd.Close()
	}
	return d.db.Close()
}
