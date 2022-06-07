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
	if !s.DagPool.CheckUserPolicy(in.User.Username, in.User.Pass, userpolicy.OnlyWrite) {
		return &proto.AddReply{Cid: ""}, userpolicy.AccessDenied
	}
	err := s.DagPool.Add(ctx, data)
	if err != nil {
		return &proto.AddReply{Cid: ""}, err
	}
	return &proto.AddReply{Cid: data.Cid().String()}, nil
}
func (s *DagPoolService) Get(ctx context.Context, in *proto.GetReq) (*proto.GetReply, error) {
	if !s.DagPool.CheckUserPolicy(in.User.Username, in.User.Pass, userpolicy.OnlyWrite) {
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
	if !s.DagPool.CheckUserPolicy(in.User.Username, in.User.Pass, userpolicy.OnlyWrite) {
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
	if !dagpooluser.CheckAddUser(in.User.Username, in.User.Pass) {
		return &proto.AddUserReply{Message: "you can not add user"}, xerrors.Errorf("you can not add user")
	}
	if !userpolicy.CheckValid(in.Policy) {
		return &proto.AddUserReply{Message: policyNotRight}, xerrors.Errorf(policyNotRight)
	}
	err := s.DagPool.AddUser(
		dagpooluser.DagPoolUser{
			Username: in.Username,
			Password: in.Password,
			Policy:   userpolicy.DagPoolPolicy(in.Policy),
			Capacity: in.Capacity,
		})
	if err != nil {
		return &proto.AddUserReply{Message: fmt.Sprintf("add user err:%v", err)}, err
	}
	return &proto.AddUserReply{Message: "ok"}, nil
}
func (s *DagPoolService) RemoveUser(ctx context.Context, in *proto.RemoveUserReq) (*proto.RemoveUserReply, error) {
	if !s.DagPool.CheckDeal(in.Username, in.Password) {
		return &proto.RemoveUserReply{Message: "you can not del user"}, xerrors.Errorf("you can not del user")
	}
	err := s.DagPool.RemoveUser(in.Username)
	if err != nil {
		return &proto.RemoveUserReply{Message: fmt.Sprintf("del user err:%v", err)}, err
	}
	return &proto.RemoveUserReply{Message: "ok"}, nil
}
func (s *DagPoolService) QueryUser(ctx context.Context, in *proto.QueryUserReq) (*proto.QueryUserReply, error) {
	if !s.DagPool.CheckDeal(in.Username, in.Password) {
		return &proto.QueryUserReply{}, xerrors.Errorf("you can not get user")
	}
	user, err := s.DagPool.QueryUser(in.Username)
	if err != nil {
		return &proto.QueryUserReply{}, err
	}
	return &proto.QueryUserReply{Username: user.Username, Policy: string(user.Policy), Capacity: user.Capacity}, nil
}
func (s *DagPoolService) UpdateUser(ctx context.Context, in *proto.UpdateUserReq) (*proto.UpdateUserReply, error) {
	if !s.DagPool.CheckDeal(in.User.Username, in.User.Pass) {
		return &proto.UpdateUserReply{Message: "you can not update user"}, xerrors.Errorf("you can not update user")
	}
	if dagpooluser.CheckAddUser(in.User.Username, in.User.Pass) {
		return &proto.UpdateUserReply{Message: "you can not update default user"}, xerrors.Errorf("you can not update default user")
	}
	var user dagpooluser.DagPoolUser
	if in.NewUsername != "" {
		user.Username = in.NewUsername
	}
	if in.NewPassword != "" {
		user.Password = in.NewPassword
	}
	if in.Policy != "" {
		if !userpolicy.CheckValid(in.Policy) {
			return &proto.UpdateUserReply{Message: policyNotRight}, xerrors.Errorf(policyNotRight)
		}
		user.Policy = userpolicy.DagPoolPolicy(in.Policy)
	}
	if in.Capacity != 0 {
		user.Capacity = in.Capacity
	}
	err := s.DagPool.UpdateUser(user)
	if err != nil {
		return &proto.UpdateUserReply{Message: fmt.Sprintf("update user err:%v", err)}, err
	}
	if in.NewUsername != "" {
		err := s.DagPool.RemoveUser(in.User.Username)
		if err != nil {
			return &proto.UpdateUserReply{Message: "update user err"}, err
		}
	}
	return &proto.UpdateUserReply{Message: "ok"}, nil
}
