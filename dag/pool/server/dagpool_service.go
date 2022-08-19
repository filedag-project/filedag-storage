package server

import (
	"context"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/pool"
	"github.com/filedag-project/filedag-storage/dag/pool/poolservice/dpuser"
	"github.com/filedag-project/filedag-storage/dag/pool/poolservice/dpuser/upolicy"
	"github.com/filedag-project/filedag-storage/dag/proto"
	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	logging "github.com/ipfs/go-log/v2"
	"golang.org/x/xerrors"
)

var log = logging.Logger("dag-pool-server")
var policyNotRight = "policy is illegal, it should be: " +
	fmt.Sprintf("%v,%v,%v", upolicy.ReadOnly, upolicy.WriteOnly, upolicy.ReadWrite)

// DagPoolServer is used to implement DagPoolServer.
type DagPoolServer struct {
	proto.UnimplementedDagPoolServer
	DagPool pool.DagPool
}

//Add is used to add a block to the dag pool server
func (s *DagPoolServer) Add(ctx context.Context, in *proto.AddReq) (*proto.AddReply, error) {
	data := blocks.NewBlock(in.GetBlock())
	err := s.DagPool.Add(ctx, data, in.User.User, in.User.Password)
	if err != nil {
		return &proto.AddReply{Cid: cid.Undef.String()}, err
	}
	return &proto.AddReply{Cid: data.Cid().String()}, nil
}

//Get is used to get a block from the dag pool server
func (s *DagPoolServer) Get(ctx context.Context, in *proto.GetReq) (*proto.GetReply, error) {
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

//GetSize is used to get the size of the block
func (s *DagPoolServer) GetSize(ctx context.Context, in *proto.GetSizeReq) (*proto.GetSizeReply, error) {
	cid, err := cid.Decode(in.Cid)
	if err != nil {
		return &proto.GetSizeReply{Size: 0}, err
	}
	size, err := s.DagPool.GetSize(ctx, cid, in.User.User, in.User.Password)
	if err != nil {
		return &proto.GetSizeReply{Size: 0}, err
	}
	return &proto.GetSizeReply{Size: int32(size)}, nil
}

//Remove is used to remove a block from the dag pool server
func (s *DagPoolServer) Remove(ctx context.Context, in *proto.RemoveReq) (*proto.RemoveReply, error) {
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

//AddUser is used to add a user to the dag pool server
func (s *DagPoolServer) AddUser(ctx context.Context, in *proto.AddUserReq) (*proto.AddUserReply, error) {
	if !upolicy.CheckValid(in.Policy) {
		return &proto.AddUserReply{Message: policyNotRight}, xerrors.Errorf(policyNotRight)
	}
	err := s.DagPool.AddUser(
		dpuser.DagPoolUser{
			Username: in.Username,
			Password: in.Password,
			Policy:   upolicy.DagPoolPolicy(in.Policy),
			Capacity: in.Capacity,
		}, in.User.User, in.User.Password)
	if err != nil {
		return &proto.AddUserReply{Message: fmt.Sprintf("add user err:%v", err)}, err
	}
	return &proto.AddUserReply{Message: "ok"}, nil
}

//RemoveUser is used to remove a user from the dag pool server
func (s *DagPoolServer) RemoveUser(ctx context.Context, in *proto.RemoveUserReq) (*proto.RemoveUserReply, error) {
	err := s.DagPool.RemoveUser(in.Username, in.User.User, in.User.Password)
	if err != nil {
		return &proto.RemoveUserReply{Message: fmt.Sprintf("del user err:%v", err)}, err
	}
	return &proto.RemoveUserReply{Message: "ok"}, nil
}

//QueryUser is used to query a user from the dag pool server
func (s *DagPoolServer) QueryUser(ctx context.Context, in *proto.QueryUserReq) (*proto.QueryUserReply, error) {
	user, err := s.DagPool.QueryUser(in.Username, in.User.User, in.User.Password)
	if err != nil {
		return &proto.QueryUserReply{}, err
	}
	return &proto.QueryUserReply{Username: user.Username, Policy: string(user.Policy), Capacity: user.Capacity}, nil
}

//UpdateUser is used to update a user from the dag pool server
func (s *DagPoolServer) UpdateUser(ctx context.Context, in *proto.UpdateUserReq) (*proto.UpdateUserReply, error) {
	user := dpuser.DagPoolUser{
		Username: in.Username,
		Password: in.NewPassword,
		Capacity: in.NewCapacity,
	}
	if in.NewPolicy != "" {
		if !upolicy.CheckValid(in.NewPolicy) {
			return &proto.UpdateUserReply{Message: policyNotRight}, xerrors.Errorf(policyNotRight)
		}
		user.Policy = upolicy.DagPoolPolicy(in.NewPolicy)
	}
	err := s.DagPool.UpdateUser(user, in.User.User, in.User.Password)
	if err != nil {
		return &proto.UpdateUserReply{Message: fmt.Sprintf("update user err:%v", err)}, err
	}
	return &proto.UpdateUserReply{Message: "ok"}, nil
}
func (s *DagPoolServer) Pin(ctx context.Context, in *proto.PinReq) (*proto.PinReply, error) {
	//if !s.DagPool.CheckUserPolicy(in.User.Username, in.User.Pass, userpolicy.OnlyWrite) {
	//	return &proto.PinReply{Message: ""}, userpolicy.AccessDenied
	//}
	c, err := cid.Decode(in.Cid)
	if err != nil {
		return &proto.PinReply{Message: ""}, err
	}
	err = s.DagPool.Pin(ctx, c, in.User.User, in.User.Password)
	if err != nil {
		return &proto.PinReply{Message: ""}, err
	}
	return &proto.PinReply{Message: c.String()}, nil
}

func (s *DagPoolServer) UnPin(ctx context.Context, in *proto.UnPinReq) (*proto.UnPinReply, error) {
	//if !s.DagPool.CheckUserPolicy(in.User.Username, in.User.Pass, userpolicy.OnlyWrite) {
	//	return &proto.UnPinReply{Message: ""}, userpolicy.AccessDenied
	//}
	c, err := cid.Decode(in.Cid)
	if err != nil {
		return &proto.UnPinReply{Message: ""}, err
	}
	err = s.DagPool.UnPin(ctx, c, in.User.User, in.User.Password)
	if err != nil {
		return &proto.UnPinReply{Message: ""}, err
	}
	return &proto.UnPinReply{Message: c.String()}, nil
}

//func (s *DagPoolService) IsPin(ctx context.Context, in *proto.IsPinReq) (*proto.IsPinReply, error) {
//	if !s.DagPool.CheckUserPolicy(in.User.Username, in.User.Pass, userpolicy.OnlyWrite) {
//		return &proto.IsPinReply{Is: false}, userpolicy.AccessDenied
//	}
//	c, err := cid.Decode(in.Cid)
//	if err != nil {
//		return &proto.IsPinReply{Is: false}, err
//	}
//	ok := s.DagPool.IsPinned(ctx, c)
//	if !ok {
//		return &proto.IsPinReply{Is: false}, err
//	}
//	return &proto.IsPinReply{Is: true}, nil
//}
