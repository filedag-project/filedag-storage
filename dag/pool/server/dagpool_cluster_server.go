package server

import (
	"context"
	"github.com/filedag-project/filedag-storage/dag/config"
	"github.com/filedag-project/filedag-storage/dag/pool"
	"github.com/filedag-project/filedag-storage/dag/proto"
	"github.com/filedag-project/filedag-storage/dag/slotsmgr"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// DagPoolClusterServer is used to implement DagPoolClusterServer.
type DagPoolClusterServer struct {
	proto.UnimplementedDagPoolClusterServer
	Cluster pool.Cluster
}

func (s *DagPoolClusterServer) InitSlots(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	if err := s.Cluster.InitSlots(); err != nil {
		return nil, status.Errorf(codes.Unknown, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *DagPoolClusterServer) AddDagNode(ctx context.Context, node *proto.DagNodeInfo) (*emptypb.Empty, error) {
	dataNodes := make([]config.DataNodeConfig, 0, len(node.Nodes))
	for _, nd := range node.Nodes {
		dataNodes = append(dataNodes, config.DataNodeConfig{
			SetIndex:   int(nd.SetIndex),
			RpcAddress: nd.RpcAddress,
		})
	}
	cfg := &config.DagNodeConfig{
		Name:         node.Name,
		Nodes:        dataNodes,
		DataBlocks:   int(node.DataBlocks),
		ParityBlocks: int(node.ParityBlocks),
	}
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
	dataNodes := make([]*proto.DataNodeInfo, 0, len(node.Nodes))
	for _, nd := range node.Nodes {
		dataNodes = append(dataNodes, &proto.DataNodeInfo{
			SetIndex:   int32(nd.SetIndex),
			RpcAddress: nd.RpcAddress,
		})
	}
	nodeInfo := &proto.DagNodeInfo{
		Name:         node.Name,
		Nodes:        dataNodes,
		DataBlocks:   int32(node.DataBlocks),
		ParityBlocks: int32(node.ParityBlocks),
	}

	return nodeInfo, nil
}

func (s *DagPoolClusterServer) RemoveDagNode(ctx context.Context, req *proto.RemoveDagNodeReq) (*proto.DagNodeInfo, error) {
	node, err := s.Cluster.RemoveDagNode(req.Name)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, err.Error())
	}
	dataNodes := make([]*proto.DataNodeInfo, 0, len(node.Nodes))
	for _, nd := range node.Nodes {
		dataNodes = append(dataNodes, &proto.DataNodeInfo{
			SetIndex:   int32(nd.SetIndex),
			RpcAddress: nd.RpcAddress,
		})
	}
	nodeInfo := &proto.DagNodeInfo{
		Name:         node.Name,
		Nodes:        dataNodes,
		DataBlocks:   int32(node.DataBlocks),
		ParityBlocks: int32(node.ParityBlocks),
	}

	return nodeInfo, nil
}

func (s *DagPoolClusterServer) MigrateSlots(ctx context.Context, req *proto.MigrateSlotsReq) (*emptypb.Empty, error) {
	newPairs := make([]slotsmgr.SlotPair, 0, len(req.Pairs))
	for _, p := range req.Pairs {
		newPairs = append(newPairs, slotsmgr.SlotPair{Start: uint64(p.Start), End: uint64(p.End)})
	}
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
	list, err := s.Cluster.Status()
	if err != nil {
		return nil, status.Errorf(codes.Unknown, err.Error())
	}
	return &proto.StatusReply{Statuses: list}, nil
}
