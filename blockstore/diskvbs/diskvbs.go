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

func (bs *diskvbs) DeleteBlock(ctx context.Context, cid cid.Cid) error {
	return nil
}

func (bs *diskvbs) Has(ctx context.Context, cid cid.Cid) (bool, error) {
	return false, nil
}

func (bs *diskvbs) Get(ctx context.Context, cid cid.Cid) (blocks.Block, error) {
	return nil, nil
}

func (bs *diskvbs) GetSize(ctx context.Context, cid cid.Cid) (int, error) {
	return 0, nil
}

func (bs *diskvbs) Put(ctx context.Context, blo blocks.Block) error {
	return nil
}

func (bs *diskvbs) PutMany(ctx context.Context, blos []blocks.Block) error {
	return nil
}

func (bs *diskvbs) AllKeysChan(ctx context.Context) (<-chan cid.Cid, error) {
	return nil, nil
}

func (bs *diskvbs) HashOnRead(enabled bool) {

}
