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

var _ format.DAGService = &DagPoolClient{}
var _ PoolClient = &DagPoolClient{}

//go:generate go run github.com/golang/mock/mockgen -destination=mocks/mock_poolclient.go -package=mocks . PoolClient

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
	User     *proto.PoolUser
}

func NewPoolClient(addr, user, password string) (*DagPoolClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Errorf("did not connect: %v", err)
		return nil, err
	}
	c := proto.NewDagPoolClient(conn)
	return &DagPoolClient{
		DPClient: c,
		Conn:     conn,
		User: &proto.PoolUser{
			Username: user,
			Pass:     password,
		},
	}, nil
}

func (p *DagPoolClient) Close(ctx context.Context) {
	p.Conn.Close()
}

func (p *DagPoolClient) Get(ctx context.Context, cid cid.Cid) (format.Node, error) {
	log.Infof(cid.String())
	get, err := p.DPClient.Get(ctx, &proto.GetReq{Cid: cid.String(), User: p.User})
	if err != nil {
		return nil, err
	}
	return legacy.DecodeNode(ctx, blocks.NewBlock(get.Block))
}

func (p *DagPoolClient) Add(ctx context.Context, node format.Node) error {
	_, err := p.DPClient.Add(ctx, &proto.AddReq{Block: node.RawData(), User: p.User})
	if err != nil {
		return err
	}
	return nil
}

func (p *DagPoolClient) Remove(ctx context.Context, cid cid.Cid) error {
	reply, err := p.DPClient.Remove(ctx, &proto.RemoveReq{
		Cid:  cid.String(),
		User: p.User})
	if err != nil {
		return err
	}
	log.Infof("delete sucess %v ", reply.Message)
	return err
}

func (p *DagPoolClient) AddMany(ctx context.Context, nodes []format.Node) error {
	return xerrors.Errorf("implement me")
}

func (p *DagPoolClient) GetMany(ctx context.Context, cids []cid.Cid) <-chan *format.NodeOption {
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

func (p *DagPoolClient) RemoveMany(ctx context.Context, cids []cid.Cid) error {
	return xerrors.Errorf("implement me")
}

var _ format.DAGService = &DagPoolClient{}
var _ PoolClient = &DagPoolClient{}

func (p DagPoolClient) Pin(ctx context.Context, cid cid.Cid) error {
	s := strings.Split((ctx.Value("user")).(string), ",")
	if len(s) != 2 {
		return userpolicy.AccessDenied
	}
	reply, err := p.DPClient.Pin(ctx, &proto.PinReq{
		Cid: cid.String(),
		User: &proto.PoolUser{
			Username: s[0],
			Pass:     s[1],
		}})
	if err != nil {
		return err
	}
	log.Infof("pin sucess %v ", reply.Message)
	return err
}

func (p DagPoolClient) UnPin(ctx context.Context, cid cid.Cid) error {
	s := strings.Split((ctx.Value("user")).(string), ",")
	if len(s) != 2 {
		return userpolicy.AccessDenied
	}
	reply, err := p.DPClient.UnPin(ctx, &proto.UnPinReq{
		Cid: cid.String(),
		User: &proto.PoolUser{
			Username: s[0],
			Pass:     s[1],
		}})
	if err != nil {
		return err
	}
	log.Infof("unpin sucess %v ", reply.Message)
	return err
}
