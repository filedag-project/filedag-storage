package dag

import (
	"context"
	"github.com/filedag-project/filedag-storage/dag/pool"
	"github.com/filedag-project/filedag-storage/dag/pool/dagpooluser"
	"github.com/filedag-project/filedag-storage/dag/pool/userpolicy"
	"github.com/filedag-project/filedag-storage/dag/pool/utils"
	logging "github.com/ipfs/go-log/v2"
	"io"

	"github.com/filedrive-team/filehelper/importer"
	"github.com/ipfs/go-cid"
	ufsio "github.com/ipfs/go-unixfs/io"
)

const importerBatchNum = 32

var log = logging.Logger("pool")
var _ DAGPool = (*simplePool)(nil)

type simplePool struct {
	dagserv *pool.DagPool
}

func (p *simplePool) AddUser(ctx context.Context, user, pass string, policy userpolicy.DagPoolPolicy, cap uint64) error {
	return p.dagserv.Iam.AddUser(dagpooluser.DagPoolUser{Username: user, Password: pass, Policy: policy, Capacity: cap})
}

func NewSimplePool() (*simplePool, error) {
	service, err := pool.NewDagPoolService()
	if err != nil {
		return nil, err
	}
	return &simplePool{
		dagserv: service,
	}, nil
}

func (p *simplePool) Add(ctx context.Context, r io.ReadCloser) (cidstr string, err error) {
	nd, err := utils.BalanceNode(ctx, r, p.dagserv, p.dagserv.CidBuilder)
	if err != nil {
		return "", err
	}
	return nd.String(), nil
}

func (p *simplePool) AddWithSize(ctx context.Context, r io.ReadCloser, fsize int64) (cidstr string, err error) {
	ndcid, err := importer.BalanceNode(ctx, r, fsize, p.dagserv, p.dagserv.CidBuilder, p.dagserv.ImporterBatchNum)
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
