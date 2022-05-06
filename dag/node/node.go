package node

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/filedag-project/filedag-storage/kv"
	"github.com/filedag-project/filedag-storage/kv/mutcask"
	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	"github.com/klauspost/reedsolomon"
	"hash/crc32"
	"sync"
)

type Config struct {
	Batch   int
	Path    string
	CaskNum int
}

var _ blockstore.Blockstore = (*DagNode)(nil)

type DagNode struct {
	sync.Mutex
	cfg            *CaskConfig
	caskMap        *CaskMap
	createCaskChan chan *createCaskRequst
	close          func()
	closeChan      chan struct{}
}

func NewDagNode(cfg *Config) (*blostore, error) {
	if cfg.Batch == 0 {
		cfg.Batch = 4
	}
	if cfg.CaskNum == 0 {
		cfg.CaskNum = 2
	}
	mc, err := mutcask.NewMutcask(mutcask.CaskNumConf(cfg.CaskNum), mutcask.PathConf(cfg.Path))
	if err != nil {
		return nil, err
	}
	return &blostore{
		batch: cfg.Batch,
		kv:    mc,
	}, nil
}

func (d DagNode) DeleteBlock(cid cid.Cid) error {
	keyCode := sha256String(cid.String())
	id := d.fileID(keyCode)
	cask, has := d.caskMap.Get(id)
	if !has {
		return nil
	}
	return cask.Delete(keyCode)
}

func (d DagNode) Has(cid cid.Cid) (bool, error) {
	_, err := d.GetSize(cid)
	if err != nil {
		if err == kv.ErrNotFound {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (d DagNode) Get(cid cid.Cid) (blocks.Block, error) {
	keyCode := sha256String(cid.String())
	id := d.fileID(keyCode)
	cask, has := d.caskMap.Get(id)
	if !has {
		fmt.Println("********")
		return nil, kv.ErrNotFound
	}
	merged := make([][]byte, 8)
	bytes, err := cask.Read(keyCode)
	merged = append(merged, bytes)
	enc, _ := reedsolomon.New(5, 3)
	var data []byte
	err = enc.EncodeIdx(data, 8, merged)
	if err != nil {
		return nil, err
	}
	b, err := blocks.NewBlockWithCid(data, cid)
	if err == blocks.ErrWrongHash {
		return nil, blockstore.ErrHashMismatch
	}
	return b, err
}

func (d DagNode) GetSize(cid cid.Cid) (int, error) {
	keyCode := sha256String(cid.String())
	id := d.fileID(keyCode)
	cask, has := d.caskMap.Get(id)
	if !has {
		return -1, kv.ErrNotFound
	}
	return cask.Size(keyCode)
}

func (d DagNode) Put(block blocks.Block) (err error) {
	keyCode := sha256String(block.Cid().String())
	// Create an encoder with 5 data and 3 parity slices.
	enc, _ := reedsolomon.New(5, 3)
	bytes := block.RawData()
	shards, _ := enc.Split(bytes)
	err = enc.Encode(shards)
	ok, err := enc.Verify(shards)
	if ok && err == nil {
		fmt.Println("encode ok")
	}
	var cask *Cask
	var has bool
	id := d.fileID(keyCode)
	cask, has = d.caskMap.Get(id)
	if !has {
		done := make(chan error)
		d.createCaskChan <- &createCaskRequst{
			id:   id,
			done: done,
		}
		if err := <-done; err != ErrNone {
			return err
		}
		cask, _ = d.caskMap.Get(id)
	}
	for _, shard := range shards {
		err = cask.Put(keyCode, shard)
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

func (d *DagNode) fileID(key string) uint32 {
	crc := crc32.ChecksumIEEE([]byte(key))
	return crc % d.cfg.CaskNum
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

//const lockFileName = "repo.lock"
//
//var _ kv.KVDB = (*dagnode)(nil)
//
//type dagnode struct {
//	sync.Mutex
//	cfg *Config
//	//caskMap        *CaskMap
//	createCaskChan chan *createCaskRequst
//	close          func()
//	closeChan      chan struct{}
//}
//
//func NewDagNode(opts ...Option) (*dagnode, error) {
//	m := &dagnode{
//		cfg:            defaultConfig(),
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
//	//unlockRepo, err := fslock.Lock(repoPath, lockFileName)
//	//if err != nil {
//	//	return nil, xerrors.Errorf("could not lock the repo: %w", err)
//	//}
//	return m, nil
//}
//
//type createCaskRequst struct {
//	id   uint32
//	done chan error
//}
//
//func (d dagnode) Put(s string, bytes []byte) error {
//	panic("implement me")
//}
//
//func (d dagnode) Delete(s string) error {
//	panic("implement me")
//}
//
//func (d dagnode) Get(s string) ([]byte, error) {
//	panic("implement me")
//}
//
//func (d dagnode) Size(s string) (int, error) {
//	panic("implement me")
//}
//
//func (d dagnode) AllKeysChan(ctx context.Context) (chan string, error) {
//	panic("implement me")
//}
//
//func (d dagnode) Close() error {
//	panic("implement me")
//}
