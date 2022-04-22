package diskvbs

import (
	"context"

	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
)

func NewDiskvBS() (blockstore.Blockstore, error) {
	return nil, nil
}

type diskvbs struct{}

var _ blockstore.Blockstore = (*diskvbs)(nil)

func (bs *diskvbs) DeleteBlock(cid cid.Cid) error {
	return nil
}

func (bs *diskvbs) Has(cid cid.Cid) (bool, error) {
	return false, nil
}

func (bs *diskvbs) Get(cid cid.Cid) (blocks.Block, error) {
	return nil, nil
}

func (bs *diskvbs) GetSize(cid cid.Cid) (int, error) {
	return 0, nil
}

func (bs *diskvbs) Put(blo blocks.Block) error {
	return nil
}

func (bs *diskvbs) PutMany(blos []blocks.Block) error {
	return nil
}

func (bs *diskvbs) AllKeysChan(ctx context.Context) (<-chan cid.Cid, error) {
	return nil, nil
}

func (bs *diskvbs) HashOnRead(enabled bool) {

}
