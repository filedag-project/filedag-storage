package server

import (
	"context"
	"fmt"
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
var policyNotRight = "policy not right ,must be:" +
	fmt.Sprintf("%v,%v,%v", userpolicy.OnlyRead, userpolicy.OnlyWrite, userpolicy.ReadWrite)

// DagPoolService is used to implement DagPoolServer.
type DagPoolService struct {
	proto.UnimplementedDagPoolServer
	DagPool pool.DagPool
}

func (s *DagPoolService) Add(ctx context.Context, in *proto.AddReq) (*proto.AddReply, error) {
	data := blocks.NewBlock(in.GetBlock())
	err := s.DagPool.Add(ctx, data, in.User.User, in.User.Password)
	if err != nil {
		return &proto.AddReply{Cid: cid.Undef.String()}, err
	}
	return &proto.AddReply{Cid: data.Cid().String()}, nil
}

func (s *DagPoolService) Get(ctx context.Context, in *proto.GetReq) (*proto.GetReply, error) {
	cid, err := cid.Decode(in.Cid)
	if err != nil {
		return &proto.GetReply{Block: nil}, err
	}
	get, err := s.DagPool.Get(ctx, cid, in.User.User, in.User.Password)
	if err != nil {
		return &proto.GetReply{Block: nil}, err
	}
	return &proto.GetReply{Block: get.RawData()}, nil
}

func (s *DagPoolService) Remove(ctx context.Context, in *proto.RemoveReq) (*proto.RemoveReply, error) {
	c, err := cid.Decode(in.Cid)
	if err != nil {
		return &proto.RemoveReply{Message: ""}, err
	}
	err = s.DagPool.Remove(ctx, c, in.User.User, in.User.Password)
	if err != nil {
		return &proto.RemoveReply{Message: ""}, err
	}
	return &proto.RemoveReply{Message: c.String()}, nil
}

func (s *DagPoolService) AddUser(ctx context.Context, in *proto.AddUserReq) (*proto.AddUserReply, error) {
	if !userpolicy.CheckValid(in.Policy) {
		return &proto.AddUserReply{Message: policyNotRight}, xerrors.Errorf(policyNotRight)
	}
	err := s.DagPool.AddUser(
		dagpooluser.DagPoolUser{
			Username: in.Username,
			Password: in.Password,
			Policy:   userpolicy.DagPoolPolicy(in.Policy),
			Capacity: in.Capacity,
		}, in.User.User, in.User.Password)
	if err != nil {
		return &proto.AddUserReply{Message: fmt.Sprintf("add user err:%v", err)}, err
	}
	return &proto.AddUserReply{Message: "ok"}, nil
}

func (s *DagPoolService) RemoveUser(ctx context.Context, in *proto.RemoveUserReq) (*proto.RemoveUserReply, error) {
	err := s.DagPool.RemoveUser(in.Username, in.User.User, in.User.Password)
	if err != nil {
		return &proto.RemoveUserReply{Message: fmt.Sprintf("del user err:%v", err)}, err
	}
	return &proto.RemoveUserReply{Message: "ok"}, nil
}

func (s *DagPoolService) QueryUser(ctx context.Context, in *proto.QueryUserReq) (*proto.QueryUserReply, error) {
	user, err := s.DagPool.QueryUser(in.Username, in.User.User, in.User.Password)
	if err != nil {
		return &proto.QueryUserReply{}, err
	}
	return &proto.QueryUserReply{Username: user.Username, Policy: string(user.Policy), Capacity: user.Capacity}, nil
}

func (s *DagPoolService) UpdateUser(ctx context.Context, in *proto.UpdateUserReq) (*proto.UpdateUserReply, error) {
	user := dagpooluser.DagPoolUser{
		Username: in.Username,
		Password: in.NewPassword,
		Capacity: in.NewCapacity,
	}
	if in.NewPolicy != "" {
		if !userpolicy.CheckValid(in.NewPolicy) {
			return &proto.UpdateUserReply{Message: policyNotRight}, xerrors.Errorf(policyNotRight)
		}
		user.Policy = userpolicy.DagPoolPolicy(in.NewPolicy)
	}
	err := s.DagPool.UpdateUser(user, in.User.User, in.User.Password)
	if err != nil {
		return &proto.UpdateUserReply{Message: fmt.Sprintf("update user err:%v", err)}, err
	}
	return &proto.UpdateUserReply{Message: "ok"}, nil
}
