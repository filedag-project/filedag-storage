package pool

import (
	"context"
	"io"

	"github.com/filedrive-team/filehelper"
	"github.com/filedrive-team/filehelper/importer"
	"github.com/ipfs/go-cid"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	format "github.com/ipfs/go-ipld-format"
)

const importerBatchNum = 32

var _ DAGPool = (*simplePool)(nil)

type simplePool struct {
	bs         blockstore.Blockstore
	dagserv    format.DAGService
	cidBuilder cid.Builder
}

func (p *simplePool) Add(ctx context.Context, r io.ReadCloser) (cidstr string, err error) {
	nd, err := filehelper.BalanceNode(r, p.dagserv, p.cidBuilder)
	if err != nil {
		return "", err
	}
	return nd.String(), nil
}

func (p *simplePool) AddWithSize(ctx context.Context, r io.ReadCloser, fsize int64) (cidstr string, err error) {
	ndcid, err := importer.BalanceNode(ctx, r, fsize, p.dagserv, p.cidBuilder, importerBatchNum)
	if err != nil {
		return "", err
	}
	return ndcid.String(), nil
}

func (p *simplePool) Get(ctx context.Context, cidstr string) (r io.ReadCloser, err error) {
	return nil, nil
}
