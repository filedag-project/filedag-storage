package pool

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/node"
	"github.com/filedag-project/filedag-storage/dag/pool/config"
	"github.com/filedag-project/filedag-storage/dag/pool/dagpooluser"
	"github.com/filedag-project/filedag-storage/dag/pool/referencecount"
	"github.com/filedag-project/filedag-storage/dag/pool/userpolicy"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	"github.com/google/martian/log"
	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-blockservice"
	bserv "github.com/ipfs/go-blockservice"
	cid "github.com/ipfs/go-cid"
	offline "github.com/ipfs/go-ipfs-exchange-offline"
	format "github.com/ipfs/go-ipld-format"
	legacy "github.com/ipfs/go-ipld-legacy"
	"github.com/ipfs/go-merkledag"
	"io/ioutil"
	"strings"
	// blank import is used to register the IPLD raw codec
	_ "github.com/ipld/go-ipld-prime/codec/raw"
)

// DagPool is an IPFS Merkle DAG service.
// - the root is virtual (like a forest)
// - stores nodes' data in a BlockService
// TODO: should cache Nodes that are in memory, and be
//       able to free some of them when vm pressure is high
type DagPool struct {
	Blocks           []bserv.BlockService
	Iam              dagpooluser.IdentityUserSys
	refer            referencecount.IdentityRefe
	CidBuilder       cid.Builder
	ImporterBatchNum int
	TheNode          RecordSys
}

// NewDagPoolService constructs a new DAGService (using the default implementation).
// Note that the default implementation is also an ipld.LinkGetter.
func NewDagPoolService(path string) (*DagPool, error) {
	cidBuilder, err := merkledag.PrefixForCidVersion(0)
	if err != nil {
		return nil, err
	}
	var cfg config.PoolConfig
	file, err := ioutil.ReadFile(path)
	json.Unmarshal(file, &cfg)
	db, err := uleveldb.OpenDb(cfg.LeveldbPath)
	if err != nil {
		return nil, err
	}
	i, err := dagpooluser.NewIdentityUserSys(db)
	if err != nil {
		return nil, err
	}
	r, err := referencecount.NewIdentityRefe(db)
	var dn []blockservice.BlockService
	for _, nc := range cfg.NodesConfig {
		bs, err := node.NewDagNode(&node.Config{
			CaskNum: nc.CaskNum,
			Batch:   nc.Batch,
			Path:    nc.Path,
		})
		if err != nil {
			log.Errorf("new dagnode err:%v", err)
			return nil, err
		}
		dn = append(dn, blockservice.New(bs, offline.Exchange(bs)))
	}
	return &DagPool{Blocks: dn, Iam: i, refer: r, CidBuilder: cidBuilder, ImporterBatchNum: cfg.ImporterBatchNum, TheNode: NewRecordSys(db)}, nil
}

// CheckPolicy check user policy
func (d *DagPool) CheckPolicy(ctx context.Context, policy userpolicy.DagPoolPolicy) bool {
	s := strings.Split((ctx.Value("user")).(string), ",")
	if len(s) != 2 {
		return false
	}
	return d.Iam.CheckUserPolicy(s[0], s[1], policy)
}

// GetNode get the DagNode
func (d *DagPool) GetNode(ctx context.Context, c cid.Cid) bserv.BlockService {
	//todo mul node
	get, err := d.TheNode.Get(c.String())
	if err != nil {
		return nil
	}
	return d.Blocks[get]
}

// UseNode get the DagNode
func (d *DagPool) UseNode(ctx context.Context, c cid.Cid) bserv.BlockService {
	//todo mul node
	err := d.TheNode.Add(c.String(), 0)
	if err != nil {
		return nil
	}
	return d.Blocks[0]
}

// GetNodes get the DagNode
func (d *DagPool) GetNodes(ctx context.Context, cids []cid.Cid) map[bserv.BlockService][]cid.Cid {
	//todo mul node
	//
	m := make(map[bserv.BlockService][]cid.Cid)
	for _, c := range cids {
		get, err := d.TheNode.Get(c.String())
		if err != nil {
			return nil
		}
		m[d.Blocks[get]] = append(m[d.Blocks[get]], c)
	}
	return m
}

// UseNodes get the DagNode
func (d *DagPool) UseNodes(ctx context.Context, c []cid.Cid) bserv.BlockService {
	//todo mul node
	err := d.TheNode.Add(c[0].String(), 0)
	if err != nil {
		return nil
	}
	return d.Blocks[0]
}

// Add adds a node to the DagPool, storing the block in the BlockService
func (d *DagPool) Add(ctx context.Context, nd format.Node) error {
	if !d.CheckPolicy(ctx, userpolicy.OnlyWrite) {
		return userpolicy.AccessDenied
	}
	if d == nil { // FIXME remove this assertion. protect with constructor invariant
		return fmt.Errorf("DagPool is nil")
	}
	err := d.refer.AddReference(nd.Cid().String())
	if err != nil {
		return err
	}
	return d.UseNode(ctx, nd.Cid()).AddBlock(nd)
}

func (d *DagPool) AddMany(ctx context.Context, nds []format.Node) error {
	if !d.CheckPolicy(ctx, userpolicy.OnlyWrite) {
		return userpolicy.AccessDenied
	}
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
	return d.UseNodes(ctx, cids).AddBlocks(blks)
}

// Get retrieves a node from the DagPool, fetching the block in the BlockService
func (d *DagPool) Get(ctx context.Context, c cid.Cid) (format.Node, error) {
	if !d.CheckPolicy(ctx, userpolicy.OnlyRead) {
		return nil, userpolicy.AccessDenied
	}
	if d == nil {
		return nil, fmt.Errorf("DagPool is nil")
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	b, err := d.GetNode(ctx, c).GetBlock(ctx, c)
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
	if !d.CheckPolicy(ctx, userpolicy.OnlyWrite) {
		return userpolicy.AccessDenied
	}
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
	if !d.CheckPolicy(ctx, userpolicy.OnlyWrite) {
		return userpolicy.AccessDenied
	}
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
func (d *DagPool) GetMany(ctx context.Context, keys []cid.Cid) <-chan *format.NodeOption {
	if !d.CheckPolicy(ctx, userpolicy.OnlyRead) {
		return nil
	}
	m := d.GetNodes(ctx, keys)
	var a <-chan *format.NodeOption
	for _, b := range d.Blocks {
		a = getNodesFromBG(ctx, b, m[b])
	}
	return a
}

func dedupKeys(keys []cid.Cid) []cid.Cid {
	set := cid.NewSet()
	for _, c := range keys {
		set.Add(c)
	}
	if set.Len() == len(keys) {
		return keys
	}
	return set.Keys()
}

func getNodesFromBG(ctx context.Context, bs bserv.BlockService, keys []cid.Cid) <-chan *format.NodeOption {
	keys = dedupKeys(keys)

	out := make(chan *format.NodeOption, len(keys))

	blocks := bs.GetBlocks(ctx, keys)
	var count int

	go func() {
		defer close(out)
		for {
			select {
			case b, ok := <-blocks:
				if !ok {
					if count != len(keys) {
						out <- &format.NodeOption{Err: fmt.Errorf("failed to fetch all nodes")}
					}
					return
				}

				nd, err := legacy.DecodeNode(ctx, b)
				if err != nil {
					out <- &format.NodeOption{Err: err}
					return
				}

				out <- &format.NodeOption{Node: nd}
				count++

			case <-ctx.Done():
				out <- &format.NodeOption{Err: ctx.Err()}
				return
			}
		}
	}()
	return out
}

// GetLinks is the type of function passed to the EnumerateChildren function(s)
// for getting the children of an IPLD node.
type GetLinks func(context.Context, cid.Cid) ([]*format.Link, error)

var _ format.LinkGetter = &DagPool{}
var _ format.NodeGetter = &DagPool{}
var _ format.DAGService = &DagPool{}

//import (
//	"context"
//	"github.com/filedag-project/filedag-storage/dag/node"
//	storagekv "github.com/filedag-project/filedag-storage/kv"
//	blocks "github.com/ipfs/go-block-format"
//	"github.com/ipfs/go-cid"
//	blockstore "github.com/ipfs/go-ipfs-blockstore"
//	"golang.org/x/xerrors"
//	"strings"
//	"sync"
//)
//
//const lockFileName = "repo.lock"
//
//var _ blockstore.Blockstore = (*DagPool)(nil)
//
//type DagPool struct {
//	kv    storagekv.KVDB
//	batch int
//}
//
//func NewDagPool(cfg *Config) (*DagPool, error) {
//	//if cfg.Batch == 0 {
//	//	cfg.Batch = default_batch_num
//	//}
//	//if cfg.CaskNum == 0 {
//	//	cfg.CaskNum = default_cask_num
//	//}
//	mc, err := node.NewDagNode(node.CaskNumConf(cfg.CaskNum), node.PathConf(cfg.Path))
//	if err != nil {
//		return nil, err
//	}
//	return &DagPool{
//		batch: cfg.Batch,
//		kv:    mc,
//	}, nil
//}
//
//func (d *DagPool) DeleteBlock(cid cid.Cid) error {
//	return d.kv.Delete(cid.String())
//}
//
//func (d *DagPool) Has(cid cid.Cid) (bool, error) {
//	_, err := d.kv.Size(cid.String())
//	if err != nil {
//		if err == storagekv.ErrNotFound {
//			return false, nil
//		}
//		return false, err
//	}
//
//	return true, nil
//}
//
//func (d *DagPool) Get(cid cid.Cid) (blocks.Block, error) {
//	data, err := d.kv.Get(cid.String())
//	if err != nil {
//		if err == storagekv.ErrNotFound {
//			return nil, blockstore.ErrNotFound
//		}
//		return nil, err
//	}
//	b, err := blocks.NewBlockWithCid(data, cid)
//	if err == blocks.ErrWrongHash {
//		return nil, blockstore.ErrHashMismatch
//	}
//	return b, err
//}
//
//func (d *DagPool) GetSize(cid cid.Cid) (int, error) {
//	n, err := d.kv.Size(cid.String())
//	if err != nil && err == storagekv.ErrNotFound {
//		return -1, blockstore.ErrNotFound
//	}
//	return n, err
//}
//
//func (d *DagPool) Put(block blocks.Block) error {
//	return d.kv.Put(block.Cid().String(), block.RawData())
//}
//
//func (d *DagPool) PutMany(blos []blocks.Block) error {
//	var errlist []string
//	var wg sync.WaitGroup
//	batchChan := make(chan struct{}, d.batch)
//	wg.Add(len(blos))
//	for _, blo := range blos {
//		go func(d *DagPool, block blocks.Block) {
//			defer func() {
//				<-batchChan
//			}()
//			batchChan <- struct{}{}
//			err := d.kv.Put(blo.Cid().String(), blo.RawData())
//			if err != nil {
//				errlist = append(errlist, err.Error())
//			}
//		}(d, blo)
//	}
//	wg.Wait()
//	if len(errlist) > 0 {
//		return xerrors.New(strings.Join(errlist, "\n"))
//	}
//	return nil
//}
//
//func (d DagPool) AllKeysChan(ctx context.Context) (<-chan cid.Cid, error) {
//	panic("implement me")
//}
//
//func (d DagPool) HashOnRead(enabled bool) {
//	panic("implement me")
//}
//
//type Config struct {
//	Batch   int
//	Path    string
//	CaskNum int
//}
