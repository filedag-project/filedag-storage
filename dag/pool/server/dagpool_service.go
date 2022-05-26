package server

import (
	"context"
	"github.com/filedag-project/filedag-storage/dag/pool"
	"github.com/filedag-project/filedag-storage/dag/pool/userpolicy"
	"github.com/filedag-project/filedag-storage/dag/proto"
	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("dag-pool-server")

// DagPoolService is used to implement DagPoolServer.
type DagPoolService struct {
	proto.UnimplementedDagPoolServer
	DagPool *pool.DagPool
}

func (s *DagPoolService) Add(ctx context.Context, in *proto.AddReq) (*proto.AddReply, error) {
	data := blocks.NewBlock(in.GetBlock())
	if !s.DagPool.Iam.CheckUserPolicy(in.User.Username, in.User.Pass, userpolicy.OnlyWrite) {
		return &proto.AddReply{Cid: ""}, userpolicy.AccessDenied
	}
	err := s.DagPool.Add(ctx, data)
	if err != nil {
		return &proto.AddReply{Cid: ""}, err
	}
	return &proto.AddReply{Cid: data.Cid().String()}, nil
}
func (s *DagPoolService) Get(ctx context.Context, in *proto.GetReq) (*proto.GetReply, error) {
	if !s.DagPool.Iam.CheckUserPolicy(in.User.Username, in.User.Pass, userpolicy.OnlyWrite) {
		return &proto.GetReply{Block: nil}, userpolicy.AccessDenied
	}
	cid, err := cid.Decode(in.Cid)
	if err != nil {
		return &proto.GetReply{Block: nil}, err
	}
	get, err := s.DagPool.Get(ctx, cid)
	if err != nil {
		return &proto.GetReply{Block: nil}, err
	}
	return &proto.GetReply{Block: get.RawData()}, nil
}
func (s *DagPoolService) Remove(ctx context.Context, in *proto.RemoveReq) (*proto.RemoveReply, error) {
	if !s.DagPool.Iam.CheckUserPolicy(in.User.Username, in.User.Pass, userpolicy.OnlyWrite) {
		return &proto.RemoveReply{Message: ""}, userpolicy.AccessDenied
	}
	c, err := cid.Decode(in.Cid)
	if err != nil {
		return &proto.RemoveReply{Message: ""}, err
	}
	err = s.DagPool.Remove(ctx, c)
	if err != nil {
		return &proto.RemoveReply{Message: ""}, err
	}
	return &proto.RemoveReply{Message: c.String()}, nil
}
