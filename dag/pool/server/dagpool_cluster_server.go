package server

import (
	"context"
	"github.com/filedag-project/filedag-storage/dag/pool"
	"github.com/filedag-project/filedag-storage/dag/proto"
	"github.com/filedag-project/filedag-storage/dag/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// DagPoolClusterServer is used to implement DagPoolClusterServer.
type DagPoolClusterServer struct {
	proto.UnimplementedDagPoolClusterServer
	Cluster pool.Cluster
}

func (s *DagPoolClusterServer) AddDagNode(ctx context.Context, node *proto.DagNodeInfo) (*emptypb.Empty, error) {
	cfg := utils.ToDagNodeConfig(node)
	if err := s.Cluster.AddDagNode(cfg); err != nil {
		return nil, status.Errorf(codes.Unknown, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *DagPoolClusterServer) GetDagNode(ctx context.Context, req *proto.GetDagNodeReq) (*proto.DagNodeInfo, error) {
	node, err := s.Cluster.GetDagNode(req.Name)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, err.Error())
	}

	return utils.ToProtoDagNodeInfo(node), nil
}

func (s *DagPoolClusterServer) RemoveDagNode(ctx context.Context, req *proto.RemoveDagNodeReq) (*proto.DagNodeInfo, error) {
	node, err := s.Cluster.RemoveDagNode(req.Name)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, err.Error())
	}

	return utils.ToProtoDagNodeInfo(node), nil
}

func (s *DagPoolClusterServer) MigrateSlots(ctx context.Context, req *proto.MigrateSlotsReq) (*emptypb.Empty, error) {
	newPairs := utils.ToSlotPairs(req.Pairs)
	if err := s.Cluster.MigrateSlots(req.FromDagNodeName, req.ToDagNodeName, newPairs); err != nil {
		return nil, status.Errorf(codes.Unknown, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *DagPoolClusterServer) BalanceSlots(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	if err := s.Cluster.BalanceSlots(); err != nil {
		return nil, status.Errorf(codes.Unknown, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *DagPoolClusterServer) Status(context.Context, *emptypb.Empty) (*proto.StatusReply, error) {
	st, err := s.Cluster.Status()
	if err != nil {
		return nil, status.Errorf(codes.Unknown, err.Error())
	}
	return st, nil
}

func (s *DagPoolClusterServer) RepairDataNode(ctx context.Context, req *proto.RepairDataNodeReq) (*emptypb.Empty, error) {
	if err := s.Cluster.RepairDataNode(ctx, req.DagNodeName, int(req.FromNodeIndex), int(req.RepairNodeIndex)); err != nil {
		return nil, status.Errorf(codes.Unknown, err.Error())
	}
	return &emptypb.Empty{}, nil
}
