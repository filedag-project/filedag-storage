package node

import (
	"bytes"
	"context"
	"crypto/sha256"
	"flag"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/config"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	"github.com/filedag-project/filedag-storage/kv"
	"github.com/filedag-project/filedag-storage/proto"
	"github.com/google/martian/log"
	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	"google.golang.org/grpc"
	"strings"
	"sync"
)

const lockFileName = "repo.lock"

var _ blockstore.Blockstore = (*DagNode)(nil)

type DagNode struct {
	nodes                    []proto.MutCaskClient
	db                       *uleveldb.ULevelDB
	dataBlocks, parityBlocks int
}

type SliceNode struct {
	sync.Mutex
	cfg            *config.CaskConfig
	caskMap        *CaskMap
	createCaskChan chan *createCaskRequst
	close          func()
	closeChan      chan struct{}
}

func NewDagNode(cfg config.NodeConfig) (*DagNode, error) {
	var s []proto.MutCaskClient
	for i, c := range cfg.Nodes {
		sc, err := InitSliceConn(fmt.Sprint(i), c.Ip, c.Port)
		if err != nil {
			return nil, err
		}
		s = append(s, sc)
	}
	db, _ := uleveldb.OpenDb(cfg.LevelDbPath)
	return &DagNode{s, db, cfg.DataBlocks, cfg.ParityBlocks}, nil
}

func InitSliceConn(index, ip, port string) (c proto.MutCaskClient, err error) {
	addr := flag.String("addr"+fmt.Sprint(index), fmt.Sprintf("%s:%s", ip, port), "the address to connect to")
	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Errorf("did not connect: %v", err)
	}
	//defer conn.Close()
	// init client
	c = proto.NewMutCaskClient(conn)
	return c, nil
}

//func NewSliceNode(opts ...config.Option) (*SliceNode, error) {
//	m := &SliceNode{
//		cfg:            config.DefaultConfig(),
//		createCaskChan: make(chan *createCaskRequst),
//		closeChan:      make(chan struct{}),
//	}
//	for _, opt := range opts {
//		opt(m.cfg)
//	}
//	repoPath := m.cfg.Path
//	if repoPath == "" {
//		return nil, ErrPathUndefined
//	}
//	repo, err := os.Stat(repoPath)
//	if err == nil && !repo.IsDir() {
//		return nil, ErrPath
//	}
//	if err != nil {
//		if !os.IsNotExist(err) {
//			return nil, err
//		}
//		if err := os.Mkdir(repoPath, 0755); err != nil {
//			return nil, err
//		}
//	}
//	// try to get the repo lock
//	locked, err := fslock.Locked(repoPath, lockFileName)
//	if err != nil {
//		return nil, xerrors.Errorf("could not check lock status: %w", err)
//	}
//	if locked {
//		return nil, ErrRepoLocked
//	}
//
//	unlockRepo, err := fslock.Lock(repoPath, lockFileName)
//	if err != nil {
//		return nil, xerrors.Errorf("could not lock the repo: %w", err)
//	}
//
//	m.caskMap, err = buildCaskMap(m.cfg)
//	if err != nil {
//		return nil, err
//	}
//	var once sync.Once
//	m.close = func() {
//		once.Do(func() {
//			close(m.closeChan)
//			unlockRepo.Close()
//		})
//	}
//	m.handleCreateCask()
//	return m, nil
//}
//func (m *SliceNode) handleCreateCask() {
//	go func(m *SliceNode) {
//		ids := []uint32{}
//		for {
//			select {
//			case <-m.closeChan:
//				return
//			case req := <-m.createCaskChan:
//				func() {
//					// fmt.Printf("received cask create request, id = %d\n", req.id)
//					if hasId(ids, req.id) {
//						req.done <- ErrNone
//						return
//					}
//					cask := NewCask(req.id)
//					var err error
//					// create vlog file
//					cask.vLog, err = os.OpenFile(filepath.Join(m.cfg.Path, m.vLogName(req.id)), os.O_RDWR|os.O_CREATE, 0644)
//					if err != nil {
//						req.done <- err
//						return
//					}
//					// create hintlog file
//					cask.hintLog, err = os.OpenFile(filepath.Join(m.cfg.Path, m.hintLogName(req.id)), os.O_RDWR|os.O_CREATE, 0644)
//					if err != nil {
//						req.done <- err
//						return
//					}
//					m.caskMap.Add(req.id, cask)
//					ids = append(ids, req.id)
//					req.done <- ErrNone
//				}()
//			}
//		}
//	}(m)
//}
func (m *SliceNode) vLogName(id uint32) string {
	return fmt.Sprintf("%08d%s", id, vLogSuffix)
}

func (m *SliceNode) hintLogName(id uint32) string {
	return fmt.Sprintf("%08d%s", id, hintLogSuffix)
}
func (d DagNode) DeleteBlock(cid cid.Cid) (err error) {
	ctx := context.TODO()
	keyCode := sha256String(cid.String())
	for _, node := range d.nodes {
		_, err := node.Delete(ctx, &proto.DeleteRequest{Key: keyCode})
		if err != nil {
			break
		}
	}
	return err
}

func (d DagNode) Has(cid cid.Cid) (bool, error) {
	_, err := d.GetSize(cid)
	if err != nil {
		if strings.Contains(err.Error(), kv.ErrNotFound.Error()) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (d DagNode) Get(cid cid.Cid) (blocks.Block, error) {
	ctx := context.TODO()
	keyCode := sha256String(cid.String())
	var err error
	var size int
	err = d.db.Get(cid.String(), &size)
	if err != nil {
		return nil, err
	}
	merged := make([][]byte, 0)
	for _, node := range d.nodes {
		res, err := node.Get(ctx, &proto.GetRequest{Key: keyCode})
		if err == nil {
			log.Errorf("mutcask get :%v", err)
		}
		merged = append(merged, res.DataBlock)
	}
	enc, err := NewErasure(d.dataBlocks, d.parityBlocks, int64(size))
	enc.DecodeDataBlocks(merged)
	var data []byte
	data = bytes.Join(merged, []byte(""))
	if err != nil {
		return nil, err
	}
	data = data[:size]
	b, err := blocks.NewBlockWithCid(data, cid)
	if err == blocks.ErrWrongHash {
		return nil, blockstore.ErrHashMismatch
	}
	return b, err
}

func (d DagNode) GetSize(cid cid.Cid) (int, error) {
	ctx := context.TODO()
	keyCode := sha256String(cid.String())
	var err error
	var count int64
	for _, node := range d.nodes {
		size, err := node.Size(ctx, &proto.SizeRequest{
			Key: keyCode,
		})
		if err != nil {
			return 0, err
		}
		count = count + size.Size
	}
	return int(count), err
}

func (d DagNode) Put(block blocks.Block) (err error) {
	ctx := context.TODO()
	err = d.db.Put(block.Cid().String(), len(block.RawData()))
	if err != nil {
		return err
	}
	keyCode := sha256String(block.Cid().String())
	enc, err := NewErasure(d.dataBlocks, d.parityBlocks, int64(len(block.RawData())))
	if err != nil {
		log.Errorf("newErasure fail :%v", err)
		return err
	}
	shards, err := enc.EncodeData(block.RawData())
	if err != nil {
		log.Errorf("encodeData fail :%v", err)
		return err
	}
	ok, err := enc.encoder().Verify(shards)
	if err != nil {
		log.Errorf("encode fail :%v", err)
		return err
	}
	if ok && err == nil {
		log.Infof("encode ok, the data is the same format as Encode. No data is modified")
	}
	for i, node := range d.nodes {
		_, err = node.Put(ctx, &proto.AddRequest{Key: keyCode, DataBlock: shards[i]})
		if err != nil {
			break
		}
	}

	return err
}

func (d DagNode) PutMany(blocks []blocks.Block) error {
	panic("implement me")
}

func (d DagNode) AllKeysChan(ctx context.Context) (<-chan cid.Cid, error) {
	panic("implement me")
}

func (d DagNode) HashOnRead(enabled bool) {
	panic("implement me")
}

func sha256String(s string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
}

type createCaskRequst struct {
	id   uint32
	done chan error
}

func hasId(ids []uint32, id uint32) bool {
	for _, item := range ids {
		if item == id {
			return true
		}
	}
	return false
}
