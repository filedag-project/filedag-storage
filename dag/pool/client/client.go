package client

import (
	"context"
	"github.com/filedag-project/filedag-storage/dag/pool/userpolicy"
	"github.com/filedag-project/filedag-storage/dag/proto"
	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	format "github.com/ipfs/go-ipld-format"
	legacy "github.com/ipfs/go-ipld-legacy"
	logging "github.com/ipfs/go-log/v2"
	"github.com/ipfs/go-merkledag"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
	"strings"
)

var log = logging.Logger("pool-client")

type PoolClient struct {
	pc         proto.DagPoolClient
	CidBuilder cid.Builder
	conn       *grpc.ClientConn
}

func (p PoolClient) Close(ctx context.Context) {
	p.conn.Close()
}

func NewPoolClient(addr string) (*PoolClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Errorf("did not connect: %v", err)
		return nil, err
	}
	c := proto.NewDagPoolClient(conn)
	cidBuilder, err := merkledag.PrefixForCidVersion(0)
	return &PoolClient{c, cidBuilder, conn}, nil
}

func (p PoolClient) Get(ctx context.Context, cid cid.Cid) (format.Node, error) {
	s := strings.Split((ctx.Value("user")).(string), ",")
	if len(s) != 2 {
		return nil, userpolicy.AccessDenied
	}
	log.Infof(cid.String())
	get, err := p.pc.Get(ctx, &proto.GetReq{Cid: cid.String(), User: &proto.PoolUser{
		Username: s[0],
		Pass:     s[1],
	}})
	if err != nil {
		return nil, err
	}
	return legacy.DecodeNode(ctx, blocks.NewBlock(get.Block))
}

func (p PoolClient) Add(ctx context.Context, node format.Node) error {
	s := strings.Split((ctx.Value("user")).(string), ",")
	if len(s) != 2 {
		return userpolicy.AccessDenied
	}
	_, err := p.pc.Add(ctx, &proto.AddReq{Block: node.RawData(), User: &proto.PoolUser{
		Username: s[0],
		Pass:     s[1],
	}})
	if err != nil {
		return err
	}
	return nil
}

func (p PoolClient) Remove(ctx context.Context, cid cid.Cid) error {
	s := strings.Split((ctx.Value("user")).(string), ",")
	if len(s) != 2 {
		return userpolicy.AccessDenied
	}
	reply, err := p.pc.Remove(ctx, &proto.RemoveReq{
		Cid: cid.String(),
		User: &proto.PoolUser{
			Username: s[0],
			Pass:     s[1],
		}})
	if err != nil {
		return err
	}
	log.Infof("delete sucess %v ", reply.Message)
	return err
}
func (p PoolClient) AddMany(ctx context.Context, nodes []format.Node) error {
	return xerrors.Errorf("implement me")
}
func (p PoolClient) GetMany(ctx context.Context, cids []cid.Cid) <-chan *format.NodeOption {
	out := make(chan *format.NodeOption, len(cids))
	defer close(out)
	for _, c := range cids {
		log.Infof(c.String())

		b, err := p.Get(ctx, c)
		if err != nil {
			return nil
		}

		nd, err := legacy.DecodeNode(ctx, b)
		if err != nil {
			out <- &format.NodeOption{Err: err}
			return nil
		}
		out <- &format.NodeOption{Node: nd}

	}
	return out
}
func (p PoolClient) RemoveMany(ctx context.Context, cids []cid.Cid) error {
	return xerrors.Errorf("implement me")
}

var _ format.DAGService = &PoolClient{}