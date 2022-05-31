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
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
	"strings"
)

var log = logging.Logger("pool-client")

type PoolClient interface {
	Close(ctx context.Context)
	Get(ctx context.Context, cid cid.Cid) (format.Node, error)
	Add(ctx context.Context, node format.Node) error
	Remove(ctx context.Context, cid cid.Cid) error
	AddMany(ctx context.Context, nodes []format.Node) error
	GetMany(ctx context.Context, cids []cid.Cid) <-chan *format.NodeOption
	RemoveMany(ctx context.Context, cids []cid.Cid) error
}
type DagPoolClient struct {
	DPClient proto.DagPoolClient
	Conn     *grpc.ClientConn
}

func (p DagPoolClient) Close(ctx context.Context) {
	p.Conn.Close()
}

func NewPoolClient(addr string) (*DagPoolClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Errorf("did not connect: %v", err)
		return nil, err
	}
	c := proto.NewDagPoolClient(conn)
	return &DagPoolClient{c, conn}, nil
}

func (p DagPoolClient) Get(ctx context.Context, cid cid.Cid) (format.Node, error) {
	s := strings.Split((ctx.Value("user")).(string), ",")
	if len(s) != 2 {
		return nil, userpolicy.AccessDenied
	}
	log.Infof(cid.String())
	get, err := p.DPClient.Get(ctx, &proto.GetReq{Cid: cid.String(), User: &proto.PoolUser{
		Username: s[0],
		Pass:     s[1],
	}})
	if err != nil {
		return nil, err
	}
	return legacy.DecodeNode(ctx, blocks.NewBlock(get.Block))
}

func (p DagPoolClient) Add(ctx context.Context, node format.Node) error {
	s := strings.Split((ctx.Value("user")).(string), ",")
	if len(s) != 2 {
		return userpolicy.AccessDenied
	}
	_, err := p.DPClient.Add(ctx, &proto.AddReq{Block: node.RawData(), User: &proto.PoolUser{
		Username: s[0],
		Pass:     s[1],
	}})
	if err != nil {
		return err
	}
	return nil
}

func (p DagPoolClient) Remove(ctx context.Context, cid cid.Cid) error {
	s := strings.Split((ctx.Value("user")).(string), ",")
	if len(s) != 2 {
		return userpolicy.AccessDenied
	}
	reply, err := p.DPClient.Remove(ctx, &proto.RemoveReq{
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
func (p DagPoolClient) AddMany(ctx context.Context, nodes []format.Node) error {
	return xerrors.Errorf("implement me")
}
func (p DagPoolClient) GetMany(ctx context.Context, cids []cid.Cid) <-chan *format.NodeOption {
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
func (p DagPoolClient) RemoveMany(ctx context.Context, cids []cid.Cid) error {
	return xerrors.Errorf("implement me")
}

var _ format.DAGService = &DagPoolClient{}
var _ PoolClient = &DagPoolClient{}
