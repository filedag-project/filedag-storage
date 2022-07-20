package client

import (
	"context"
	"github.com/filedag-project/filedag-storage/dag/proto"
	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	format "github.com/ipfs/go-ipld-format"
	logging "github.com/ipfs/go-log/v2"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
	"strings"
)

var log = logging.Logger("pool-client")

var _ blockstore.Blockstore = (*dagPoolClient)(nil)
var _ PoolClient = (*dagPoolClient)(nil)
var _ DataPin = &dagPoolClient{}

//go:generate go run github.com/golang/mock/mockgen -destination=mocks/mock_poolclient.go -package=mocks . PoolClient,DataPin

//PoolClient is a DAGService interface
type PoolClient interface {
	blockstore.Blockstore

	Close(ctx context.Context)
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

//
////Add add a node
//func (p *dagPoolClient) Add(ctx context.Context, node format.Node) error {
//	_, err := p.DPClient.Add(ctx, &proto.AddReq{Block: node.RawData(), User: p.User})
//	if err != nil {
//		return err
//	}
//	return nil
//}

func (p *dagPoolClient) DeleteBlock(ctx context.Context, cid cid.Cid) error {
	reply, err := p.DPClient.Remove(ctx, &proto.RemoveReq{
		Cid:  cid.String(),
		User: p.User})
	if err != nil {
		return err
	}
	log.Infof("delete sucess %v ", reply.Message)
	return err
}

func (p *dagPoolClient) Has(ctx context.Context, cid cid.Cid) (bool, error) {
	_, err := p.GetSize(ctx, cid)
	if err != nil {
		if xerrors.Is(err, format.ErrNotFound{Cid: cid}) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (p *dagPoolClient) Get(ctx context.Context, cid cid.Cid) (blocks.Block, error) {
	log.Infof(cid.String())
	get, err := p.DPClient.Get(ctx, &proto.GetReq{Cid: cid.String(), User: p.User})
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, format.ErrNotFound{Cid: cid}
		}
		return nil, err
	}
	return blocks.NewBlock(get.Block), nil
}

func (p *dagPoolClient) GetSize(ctx context.Context, cid cid.Cid) (int, error) {
	reply, err := p.DPClient.GetSize(ctx, &proto.GetSizeReq{Cid: cid.String(), User: p.User})
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return 0, format.ErrNotFound{Cid: cid}
		}
		return 0, err
	}
	return int(reply.Size), nil
}

func (p *dagPoolClient) Put(ctx context.Context, blk blocks.Block) error {
	_, err := p.DPClient.Add(ctx, &proto.AddReq{Block: blk.RawData(), User: p.User})
	if err != nil {
		return err
	}
	return nil
}

func (p *dagPoolClient) PutMany(ctx context.Context, blks []blocks.Block) error {
	for _, block := range blks {
		if err := p.Put(ctx, block); err != nil {
			return err
		}
	}
	return nil
}

func (p *dagPoolClient) AllKeysChan(ctx context.Context) (<-chan cid.Cid, error) {
	panic("unimplemented")
}

func (p *dagPoolClient) HashOnRead(enabled bool) {
	panic("unimplemented")
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
