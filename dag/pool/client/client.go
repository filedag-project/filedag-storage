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

var _ blockstore.Blockstore = (*DagPoolClient)(nil)
var _ PoolClient = (*DagPoolClient)(nil)

//go:generate go run github.com/golang/mock/mockgen -destination=mocks/mock_poolclient.go -package=mocks . PoolClient

type PoolClient interface {
	blockstore.Blockstore

	Close(ctx context.Context)
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
			User:     user,
			Password: password,
		},
	}, nil
}

func (p *DagPoolClient) Close(ctx context.Context) {
	p.Conn.Close()
}

func (p *DagPoolClient) DeleteBlock(ctx context.Context, cid cid.Cid) error {
	reply, err := p.DPClient.Remove(ctx, &proto.RemoveReq{
		Cid:  cid.String(),
		User: p.User})
	if err != nil {
		return err
	}
	log.Infof("delete sucess %v ", reply.Message)
	return err
}

func (p *DagPoolClient) Has(ctx context.Context, cid cid.Cid) (bool, error) {
	_, err := p.GetSize(ctx, cid)
	if err != nil {
		if xerrors.Is(err, format.ErrNotFound{Cid: cid}) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (p *DagPoolClient) Get(ctx context.Context, cid cid.Cid) (blocks.Block, error) {
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

func (p *DagPoolClient) GetSize(ctx context.Context, cid cid.Cid) (int, error) {
	reply, err := p.DPClient.GetSize(ctx, &proto.GetSizeReq{Cid: cid.String(), User: p.User})
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return 0, format.ErrNotFound{Cid: cid}
		}
		return 0, err
	}
	return int(reply.Size), nil
}

func (p *DagPoolClient) Put(ctx context.Context, blk blocks.Block) error {
	_, err := p.DPClient.Add(ctx, &proto.AddReq{Block: blk.RawData(), User: p.User})
	if err != nil {
		return err
	}
	return nil
}

func (p *DagPoolClient) PutMany(ctx context.Context, blks []blocks.Block) error {
	for _, block := range blks {
		if err := p.Put(ctx, block); err != nil {
			return err
		}
	}
	return nil
}

func (p *DagPoolClient) AllKeysChan(ctx context.Context) (<-chan cid.Cid, error) {
	panic("unimplemented")
}

func (p *DagPoolClient) HashOnRead(enabled bool) {
	panic("unimplemented")
}

//func (p *DagPoolClient) Get(ctx context.Context, cid cid.Cid) (format.Node, error) {
//	log.Infof(cid.String())
//	get, err := p.DPClient.Get(ctx, &proto.GetReq{Cid: cid.String(), User: p.User})
//	if err != nil {
//		return nil, err
//	}
//	return legacy.DecodeNode(ctx, blocks.NewBlock(get.Block))
//}

//func (p *DagPoolClient) Add(ctx context.Context, node format.Node) error {
//	_, err := p.DPClient.Add(ctx, &proto.AddReq{Block: node.RawData(), User: p.User})
//	if err != nil {
//		return err
//	}
//	return nil
//}

//func (p *DagPoolClient) Remove(ctx context.Context, cid cid.Cid) error {
//	reply, err := p.DPClient.Remove(ctx, &proto.RemoveReq{
//		Cid:  cid.String(),
//		User: p.User})
//	if err != nil {
//		return err
//	}
//	log.Infof("delete sucess %v ", reply.Message)
//	return err
//}

//func (p *DagPoolClient) AddMany(ctx context.Context, nodes []format.Node) error {
//	return xerrors.Errorf("implement me")
//}
//
//func (p *DagPoolClient) GetMany(ctx context.Context, cids []cid.Cid) <-chan *format.NodeOption {
//	out := make(chan *format.NodeOption, len(cids))
//	defer close(out)
//	for _, c := range cids {
//		log.Infof(c.String())
//
//		b, err := p.Get(ctx, c)
//		if err != nil {
//			return nil
//		}
//
//		nd, err := legacy.DecodeNode(ctx, b)
//		if err != nil {
//			out <- &format.NodeOption{Err: err}
//			return nil
//		}
//		out <- &format.NodeOption{Node: nd}
//
//	}
//	return out
//}
//
//func (p *DagPoolClient) RemoveMany(ctx context.Context, cids []cid.Cid) error {
//	return xerrors.Errorf("implement me")
//}

func (p *DagPoolClient) AddUser(ctx context.Context, username string, password string, capacity uint64, policy string) error {
	_, err := p.DPClient.AddUser(ctx, &proto.AddUserReq{
		Username: username,
		Password: password,
		Capacity: capacity,
		Policy:   policy,
		User:     p.User,
	})
	return err
}

func (p *DagPoolClient) QueryUser(ctx context.Context, username string) (*proto.QueryUserReply, error) {
	reply, err := p.DPClient.QueryUser(ctx, &proto.QueryUserReq{
		Username: username,
		User:     p.User,
	})
	return reply, err
}

func (p *DagPoolClient) UpdateUser(ctx context.Context, username string, newPassword string, newCapacity uint64, newPolicy string) error {
	_, err := p.DPClient.UpdateUser(ctx, &proto.UpdateUserReq{
		Username:    username,
		NewPassword: newPassword,
		NewCapacity: newCapacity,
		NewPolicy:   newPolicy,
		User:        p.User,
	})
	return err
}

func (p *DagPoolClient) RemoveUser(ctx context.Context, username string) error {
	_, err := p.DPClient.RemoveUser(ctx, &proto.RemoveUserReq{
		Username: username,
		User:     p.User,
	})
	return err
}
