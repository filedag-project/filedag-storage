package pool

import (
	"context"
	"io"

	blo "github.com/filedag-project/filedag-storage/blockstore"
	"github.com/filedrive-team/filehelper"
	"github.com/filedrive-team/filehelper/importer"
	"github.com/ipfs/go-blockservice"
	"github.com/ipfs/go-cid"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	offline "github.com/ipfs/go-ipfs-exchange-offline"
	format "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag"
	ufsio "github.com/ipfs/go-unixfs/io"
	"golang.org/x/xerrors"
)

const importerBatchNum = 32

var _ DAGPool = (*simplePool)(nil)

type simplePool struct {
	bs               blockstore.Blockstore
	dagserv          format.DAGService
	cidBuilder       cid.Builder
	importerBatchNum int
}

func NewSimplePool(cfg *SimplePoolConfig) (*simplePool, error) {
	if cfg.StorePath == "" {
		return nil, xerrors.New("Need path to set store up for dag pool")
	}
	if cfg.BatchNum == 0 {
		cfg.BatchNum = importerBatchNum
	}

	bs, err := blo.NewMutcaskbs(&blo.Config{
		CaskNum: cfg.CaskNum,
		Batch:   cfg.BatchNum,
		Path:    cfg.StorePath,
	})
	if err != nil {
		return nil, err
	}
	cidBuilder, err := merkledag.PrefixForCidVersion(0)
	if err != nil {
		return nil, err
	}
	p := &simplePool{
		bs:               bs,
		dagserv:          merkledag.NewDAGService(blockservice.New(bs, offline.Exchange(bs))),
		cidBuilder:       cidBuilder,
		importerBatchNum: cfg.BatchNum,
	}
	return p, nil
}

func (p *simplePool) Add(ctx context.Context, r io.ReadCloser) (cidstr string, err error) {
	nd, err := filehelper.BalanceNode(r, p.dagserv, p.cidBuilder)
	if err != nil {
		return "", err
	}
	return nd.String(), nil
}

func (p *simplePool) AddWithSize(ctx context.Context, r io.ReadCloser, fsize int64) (cidstr string, err error) {
	ndcid, err := importer.BalanceNode(ctx, r, fsize, p.dagserv, p.cidBuilder, p.importerBatchNum)
	if err != nil {
		return "", err
	}
	return ndcid.String(), nil
}

func (p *simplePool) Get(ctx context.Context, cidstr string) (r io.ReadSeekCloser, err error) {
	cid, err := cid.Decode(cidstr)
	if err != nil {
		return nil, err
	}

	dagNode, err := p.dagserv.Get(ctx, cid)
	if err != nil {
		return nil, err
	}
	return ufsio.NewDagReader(ctx, dagNode, p.dagserv)
}