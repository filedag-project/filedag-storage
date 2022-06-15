package client

import (
	"context"
	"github.com/filedag-project/filedag-storage/dag/proto"
	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	format "github.com/ipfs/go-ipld-format"
	legacy "github.com/ipfs/go-ipld-legacy"
	logging "github.com/ipfs/go-log/v2"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
)

var log = logging.Logger("pool-client")

var _ format.DAGService = &dagPoolClient{}
var _ PoolClient = &dagPoolClient{}
var _ DataPin = &dagPoolClient{}

//go:generate go run github.com/golang/mock/mockgen -destination=mocks/mock_poolclient.go -package=mocks . PoolClient,DataPin

//PoolClient is a DAGService interface
type PoolClient interface {
	Close(ctx context.Context)
	Get(ctx context.Context, cid cid.Cid) (format.Node, error)
	Add(ctx context.Context, node format.Node) error
	Remove(ctx context.Context, cid cid.Cid) error
	AddMany(ctx context.Context, nodes []format.Node) error
	GetMany(ctx context.Context, cids []cid.Cid) <-chan *format.NodeOption
	RemoveMany(ctx context.Context, cids []cid.Cid) error
}

//DataPin is a pin interface
type DataPin interface {
	Pin(ctx context.Context, cid cid.Cid) error
	UnPin(ctx context.Context, cid cid.Cid) error
	IsPin(ctx context.Context, cid cid.Cid) (bool, error)
}

type dagPoolClient struct {
	DPClient proto.DagPoolClient
	Conn     *grpc.ClientConn
	User     *proto.PoolUser
}

//NewPoolClient new a dagPoolClient
func NewPoolClient(addr, user, password string) (*dagPoolClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Errorf("did not connect: %v", err)
		return nil, err
	}
	c := proto.NewDagPoolClient(conn)
	return &dagPoolClient{
		DPClient: c,
		Conn:     conn,
		User: &proto.PoolUser{
			User:     user,
			Password: password,
		},
	}, nil
}

//Close  the client
func (p *dagPoolClient) Close(ctx context.Context) {
	p.Conn.Close()
}

//Get get a node by cid
func (p *dagPoolClient) Get(ctx context.Context, cid cid.Cid) (format.Node, error) {
	log.Infof(cid.String())
	get, err := p.DPClient.Get(ctx, &proto.GetReq{Cid: cid.String(), User: p.User})
	if err != nil {
		return nil, err
	}
	return legacy.DecodeNode(ctx, blocks.NewBlock(get.Block))
}

//Add add a node
func (p *dagPoolClient) Add(ctx context.Context, node format.Node) error {
	_, err := p.DPClient.Add(ctx, &proto.AddReq{Block: node.RawData(), User: p.User})
	if err != nil {
		return err
	}
	return nil
}

//Remove remove a node by cid
func (p *dagPoolClient) Remove(ctx context.Context, cid cid.Cid) error {
	reply, err := p.DPClient.Remove(ctx, &proto.RemoveReq{
		Cid:  cid.String(),
		User: p.User})
	if err != nil {
		return err
	}
	log.Infof("delete sucess %v ", reply.Message)
	return err
}

//AddMany add many nodes
func (p *dagPoolClient) AddMany(ctx context.Context, nodes []format.Node) error {
	return xerrors.Errorf("implement me")
}

//GetMany get many nodes
func (p *dagPoolClient) GetMany(ctx context.Context, cids []cid.Cid) <-chan *format.NodeOption {
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

//RemoveMany remove many nodes
func (p *dagPoolClient) RemoveMany(ctx context.Context, cids []cid.Cid) error {
	return xerrors.Errorf("implement me")
}

//AddUser add a user
func (p *dagPoolClient) AddUser(ctx context.Context, username string, password string, capacity uint64, policy string) error {
	_, err := p.DPClient.AddUser(ctx, &proto.AddUserReq{
		Username: username,
		Password: password,
		Capacity: capacity,
		Policy:   policy,
		User:     p.User,
	})
	return err
}

//QueryUser query user
func (p *dagPoolClient) QueryUser(ctx context.Context, username string) (*proto.QueryUserReply, error) {
	reply, err := p.DPClient.QueryUser(ctx, &proto.QueryUserReq{
		Username: username,
		User:     p.User,
	})
	return reply, err
}

//UpdateUser update user
func (p *dagPoolClient) UpdateUser(ctx context.Context, username string, newPassword string, newCapacity uint64, newPolicy string) error {
	_, err := p.DPClient.UpdateUser(ctx, &proto.UpdateUserReq{
		Username:    username,
		NewPassword: newPassword,
		NewCapacity: newCapacity,
		NewPolicy:   newPolicy,
		User:        p.User,
	})
	return err
}

//RemoveUser remove a user
func (p *dagPoolClient) RemoveUser(ctx context.Context, username string) error {
	_, err := p.DPClient.RemoveUser(ctx, &proto.RemoveUserReq{
		Username: username,
		User:     p.User,
	})
	return err
}

//Pin pin a node
func (p dagPoolClient) Pin(ctx context.Context, cid cid.Cid) error {
	reply, err := p.DPClient.Pin(ctx, &proto.PinReq{
		Cid:  cid.String(),
		User: p.User})
	if err != nil {
		return err
	}
	log.Infof("pin sucess %v ", reply.Message)
	return err
}

//UnPin unpin a node
func (p dagPoolClient) UnPin(ctx context.Context, cid cid.Cid) error {
	reply, err := p.DPClient.UnPin(ctx, &proto.UnPinReq{
		Cid:  cid.String(),
		User: p.User})
	if err != nil {
		return err
	}
	log.Infof("unpin sucess %v ", reply.Message)
	return err
}

//IsPin check if the cid is pinned
func (p dagPoolClient) IsPin(ctx context.Context, cid cid.Cid) (bool, error) {
	reply, err := p.DPClient.IsPin(ctx, &proto.IsPinReq{
		Cid:  cid.String(),
		User: p.User})
	if err != nil {
		return false, err
	}
	log.Infof("ispin sucess %v ", reply.Is)
	return reply.Is, err
}
