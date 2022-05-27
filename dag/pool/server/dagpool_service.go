package server

import (
	"context"
	"github.com/filedag-project/filedag-storage/dag/pool"
	"github.com/filedag-project/filedag-storage/dag/pool/dagpooluser"
	"github.com/filedag-project/filedag-storage/dag/pool/userpolicy"
	"github.com/filedag-project/filedag-storage/dag/proto"
	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	logging "github.com/ipfs/go-log/v2"
	"golang.org/x/xerrors"
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
func (s *DagPoolService) AddUser(ctx context.Context, in *proto.AddUserReq) (*proto.AddUserReply, error) {
	if !dagpooluser.CheckAddUser(in.User, in.Pass) {
		return &proto.AddUserReply{Message: ""}, xerrors.Errorf("you can not add user")
	}
	err := s.DagPool.Iam.AddUser(
		dagpooluser.DagPoolUser{
			Username: in.Username,
			Password: in.Password,
			Policy:   userpolicy.DagPoolPolicy(in.Policy),
			Capacity: in.Capacity,
		})
	if err != nil {
		return &proto.AddUserReply{Message: ""}, err
	}
	return &proto.AddUserReply{Message: "ok"}, nil
}
func (s *DagPoolService) RemoveUser(ctx context.Context, in *proto.RemoveUserReq) (*proto.RemoveUserReply, error) {
	if !s.DagPool.Iam.CheckDeal(in.Username, in.Password) {
		return &proto.RemoveUserReply{Message: ""}, xerrors.Errorf("you can not del user")
	}
	err := s.DagPool.Iam.RemoveUser(in.Username)
	if err != nil {
		return &proto.RemoveUserReply{Message: ""}, err
	}
	return &proto.RemoveUserReply{Message: "ok"}, nil
}
func (s *DagPoolService) QueryUser(ctx context.Context, in *proto.QueryUserReq) (*proto.QueryUserReply, error) {
	if !s.DagPool.Iam.CheckDeal(in.Username, in.Password) {
		return &proto.QueryUserReply{}, xerrors.Errorf("you can not get user")
	}
	user, err := s.DagPool.Iam.QueryUser(in.Username)
	if err != nil {
		return &proto.QueryUserReply{}, err
	}
	return &proto.QueryUserReply{Username: user.Username, Policy: string(user.Policy), Capacity: user.Capacity}, nil
}
func (s *DagPoolService) UpdateUser(ctx context.Context, in *proto.UpdateUserReq) (*proto.UpdateUserReply, error) {
	if !s.DagPool.Iam.CheckDeal(in.Username, in.Password) {
		return &proto.UpdateUserReply{Message: ""}, xerrors.Errorf("you can not update user")
	}
	var user dagpooluser.DagPoolUser
	if in.Password != "" {
		user.Password = in.Password
	}
	if in.Policy != "" {
		user.Policy = userpolicy.DagPoolPolicy(in.Policy)
	}
	if in.Capacity != 0 {
		user.Capacity = in.Capacity
	}
	err := s.DagPool.Iam.UpdateUser(user)
	if err != nil {
		return &proto.UpdateUserReply{Message: ""}, err
	}
	return &proto.UpdateUserReply{Message: "ok"}, nil
}
