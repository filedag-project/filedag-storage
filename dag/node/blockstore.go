package node

import (
	"context"
	storagekv "github.com/filedag-project/filedag-storage/kv"
	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
)

type blostore struct {
	kv    storagekv.KVDB
	batch int
}

func (b blostore) DeleteBlock(cid cid.Cid) error {
	panic("implement me")
}

func (b blostore) Has(cid cid.Cid) (bool, error) {
	panic("implement me")
}

func (b blostore) Get(cid cid.Cid) (blocks.Block, error) {
	panic("implement me")
}

func (b blostore) GetSize(cid cid.Cid) (int, error) {
	panic("implement me")
}

func (b blostore) Put(block blocks.Block) error {
	panic("implement me")
}

func (b blostore) PutMany(blocks []blocks.Block) error {
	panic("implement me")
}

func (b blostore) AllKeysChan(ctx context.Context) (<-chan cid.Cid, error) {
	panic("implement me")
}

func (b blostore) HashOnRead(enabled bool) {
	panic("implement me")
}
