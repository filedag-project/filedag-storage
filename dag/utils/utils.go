package utils

import (
	"github.com/filedag-project/filedag-storage/dag/config"
	"github.com/filedag-project/filedag-storage/dag/proto"
	"github.com/filedag-project/filedag-storage/dag/slotsmgr"
)

func ToDagNodeConfig(node *proto.DagNodeInfo) *config.DagNodeConfig {
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
	return cfg
}

func ToProtoDagNodeInfo(node *config.DagNodeConfig) *proto.DagNodeInfo {
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
	return nodeInfo
}

func ToSlotPairs(pairs []*proto.SlotPair) []slotsmgr.SlotPair {
	newPairs := make([]slotsmgr.SlotPair, 0, len(pairs))
	for _, p := range pairs {
		newPairs = append(newPairs, slotsmgr.SlotPair{Start: uint64(p.Start), End: uint64(p.End)})
	}
	return newPairs
}

func ToProtoSlotPairs(pairs []slotsmgr.SlotPair) []*proto.SlotPair {
	newPairs := make([]*proto.SlotPair, 0, len(pairs))
	for _, p := range pairs {
		newPairs = append(newPairs, &proto.SlotPair{Start: uint32(p.Start), End: uint32(p.End)})
	}
	return newPairs
}
