package client

import (
	"context"
	"github.com/filedag-project/filedag-storage/dag/proto"
	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-blockservice"
	"github.com/ipfs/go-cid"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	offline "github.com/ipfs/go-ipfs-exchange-offline"
	format "github.com/ipfs/go-ipld-format"
	logging "github.com/ipfs/go-log/v2"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
	"strings"
)

var log = logging.Logger("pool-client")

var _ blockstore.Blockstore = (*dagPoolClient)(nil)
var _ PoolClient = (*dagPoolClient)(nil)

//go:generate go run github.com/golang/mock/mockgen -destination=mocks/mock_poolclient.go -package=mocks . PoolClient

//PoolClient is a DAGService interface
type PoolClient interface {
	blockstore.Blockstore

	Close(ctx context.Context)
}

type dagPoolClient struct {
	DPClient  proto.DagPoolClient
	Conn      *grpc.ClientConn
	User      *proto.PoolUser
	enablePin bool
}

func NewBlockService(blkstore blockstore.Blockstore) blockservice.BlockService {
	return blockservice.NewWriteThrough(blkstore, offline.Exchange(blkstore))
}

//NewPoolClient new a dagPoolClient
func NewPoolClient(addr, user, password string, enablePin bool) (*dagPoolClient, error) {
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
		enablePin: enablePin,
	}, nil
}

//Close  the client
func (p *dagPoolClient) Close(ctx context.Context) {
	p.Conn.Close()
}

//DeleteBlock delete a block
func (p *dagPoolClient) DeleteBlock(ctx context.Context, cid cid.Cid) error {
	reply, err := p.DPClient.Remove(ctx, &proto.RemoveReq{
		Cid:   cid.String(),
		User:  p.User,
		Unpin: p.enablePin,
	})
	if err != nil {
		return err
	}
	log.Debugf("delete sucess %v ", reply.Message)
	return err
}

//Has returns if the blockstore has a block with the given key
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

//Get the block with the given key, or nil if not found
func (p *dagPoolClient) Get(ctx context.Context, cid cid.Cid) (blocks.Block, error) {
	log.Debugf(cid.String())
	get, err := p.DPClient.Get(ctx, &proto.GetReq{
		Cid:  cid.String(),
		User: p.User,
	})
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, format.ErrNotFound{Cid: cid}
		}
		return nil, err
	}
	return blocks.NewBlock(get.Block), nil
}

//GetSize get the size of the block with the given key
func (p *dagPoolClient) GetSize(ctx context.Context, cid cid.Cid) (int, error) {
	reply, err := p.DPClient.GetSize(ctx, &proto.GetSizeReq{
		Cid:  cid.String(),
		User: p.User,
	})
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return 0, format.ErrNotFound{Cid: cid}
		}
		return 0, err
	}
	return int(reply.Size), nil
}

//Put  a block
func (p *dagPoolClient) Put(ctx context.Context, blk blocks.Block) error {
	_, err := p.DPClient.Add(ctx, &proto.AddReq{
		Block: blk.RawData(),
		User:  p.User,
		Pin:   p.enablePin,
	})
	if err != nil {
		return err
	}
	return nil
}

//PutMany put many nodes
func (p *dagPoolClient) PutMany(ctx context.Context, blks []blocks.Block) error {
	for _, block := range blks {
		if err := p.Put(ctx, block); err != nil {
			return err
		}
	}
	return nil
}

//AllKeysChan returns a channel from which all keys of the dag can be read.
func (p *dagPoolClient) AllKeysChan(ctx context.Context) (<-chan cid.Cid, error) {
	panic("unimplemented")
}

//HashOnRead hash on read
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
