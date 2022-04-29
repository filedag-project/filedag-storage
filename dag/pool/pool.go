package pool

import (
	"context"
	"github.com/filedag-project/filedag-storage/dag/node"
	storagekv "github.com/filedag-project/filedag-storage/kv"
	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	"golang.org/x/xerrors"
	"strings"
	"sync"
)

const lockFileName = "repo.lock"

var _ blockstore.Blockstore = (*dagpool)(nil)

type dagpool struct {
	kv    storagekv.KVDB
	batch int
}

func NewDagPool(cfg *Config) (*dagpool, error) {
	//if cfg.Batch == 0 {
	//	cfg.Batch = default_batch_num
	//}
	//if cfg.CaskNum == 0 {
	//	cfg.CaskNum = default_cask_num
	//}
	mc, err := node.NewDagNode(node.CaskNumConf(cfg.CaskNum), node.PathConf(cfg.Path))
	if err != nil {
		return nil, err
	}
	return &dagpool{
		batch: cfg.Batch,
		kv:    mc,
	}, nil
}

func (d *dagpool) DeleteBlock(cid cid.Cid) error {
	return d.kv.Delete(cid.String())
}

func (d *dagpool) Has(cid cid.Cid) (bool, error) {
	_, err := d.kv.Size(cid.String())
	if err != nil {
		if err == storagekv.ErrNotFound {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (d *dagpool) Get(cid cid.Cid) (blocks.Block, error) {
	data, err := d.kv.Get(cid.String())
	if err != nil {
		if err == storagekv.ErrNotFound {
			return nil, blockstore.ErrNotFound
		}
		return nil, err
	}
	b, err := blocks.NewBlockWithCid(data, cid)
	if err == blocks.ErrWrongHash {
		return nil, blockstore.ErrHashMismatch
	}
	return b, err
}

func (d *dagpool) GetSize(cid cid.Cid) (int, error) {
	n, err := d.kv.Size(cid.String())
	if err != nil && err == storagekv.ErrNotFound {
		return -1, blockstore.ErrNotFound
	}
	return n, err
}

func (d *dagpool) Put(block blocks.Block) error {
	return d.kv.Put(block.Cid().String(), block.RawData())
}

func (d *dagpool) PutMany(blos []blocks.Block) error {
	var errlist []string
	var wg sync.WaitGroup
	batchChan := make(chan struct{}, d.batch)
	wg.Add(len(blos))
	for _, blo := range blos {
		go func(d *dagpool, block blocks.Block) {
			defer func() {
				<-batchChan
			}()
			batchChan <- struct{}{}
			err := d.kv.Put(blo.Cid().String(), blo.RawData())
			if err != nil {
				errlist = append(errlist, err.Error())
			}
		}(d, blo)
	}
	wg.Wait()
	if len(errlist) > 0 {
		return xerrors.New(strings.Join(errlist, "\n"))
	}
	return nil
}

func (d dagpool) AllKeysChan(ctx context.Context) (<-chan cid.Cid, error) {
	panic("implement me")
}

func (d dagpool) HashOnRead(enabled bool) {
	panic("implement me")
}

type Config struct {
	Batch   int
	Path    string
	CaskNum int
}
