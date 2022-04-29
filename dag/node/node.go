package node

import (
	"context"
	storagekv "github.com/filedag-project/filedag-storage/kv"
	"github.com/filedag-project/filedag-storage/kv/mutcask"
	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
)

var _ blockstore.Blockstore = (*DagNode)(nil)

type DagNode struct {
	kv    storagekv.KVDB
	batch int
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
	panic("implement me")
}

func (d DagNode) Has(cid cid.Cid) (bool, error) {
	panic("implement me")
}

func (d DagNode) Get(cid cid.Cid) (blocks.Block, error) {
	panic("implement me")
}

func (d DagNode) GetSize(cid cid.Cid) (int, error) {
	panic("implement me")
}

func (d DagNode) Put(block blocks.Block) error {
	panic("implement me")
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
