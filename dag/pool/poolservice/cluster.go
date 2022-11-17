package poolservice

import (
	"errors"
	"github.com/filedag-project/filedag-storage/dag/config"
	"github.com/filedag-project/filedag-storage/dag/node/dagnode"
	"github.com/filedag-project/filedag-storage/dag/proto"
	"github.com/filedag-project/filedag-storage/dag/slotsmgr"
)

var ErrDagNodeNotFound = errors.New("this dag node does not exist")
var ErrDagNodeRemove = errors.New("this dag node still has slots and cannot be removed")

func (d *dagPoolService) InitSlots() error {
	// TODO
	return nil
}

func (d *dagPoolService) AddDagNode(nodeConfig *config.DagNodeConfig) error {
	if _, ok := d.dagNodesMap[nodeConfig.Name]; ok {
		return ErrDagNodeAlreadyExist
	}
	dagNode, err := dagnode.NewDagNode(*nodeConfig)
	if err != nil {
		log.Errorf("new dagnode err:%v", err)
		return err
	}
	go dagNode.RunHeartbeatCheck(d.parentCtx)
	d.dagNodesMap[nodeConfig.Name] = dagNode

	return nil
}

func (d *dagPoolService) GetDagNode(dagNodeName string) (*config.DagNodeConfig, error) {
	if nd, ok := d.dagNodesMap[dagNodeName]; ok {
		return nd.GetConfig(), nil
	}
	return nil, ErrDagNodeNotFound
}

func (d *dagPoolService) RemoveDagNode(dagNodeName string) (*config.DagNodeConfig, error) {
	if nd, ok := d.dagNodesMap[dagNodeName]; ok {
		// make sure the dagnode has no slot
		if nd.GetNumSlots() != 0 {
			return nil, ErrDagNodeRemove
		}
		cfg := nd.GetConfig()
		// update config
		delete(d.dagNodesMap, dagNodeName)
		for i, node := range d.config.Cluster {
			if node.Name == cfg.Name {
				d.config.Cluster = append(d.config.Cluster[:i], d.config.Cluster[(i+1):]...)
				break
			}
		}
		// TODO: save cluster config
		// close dagnode
		nd.Close()
		return cfg, nil
	}
	return nil, ErrDagNodeNotFound
}

func (d *dagPoolService) MigrateSlots(fromDagNodeName, toDagNodeName string, pairs []slotsmgr.SlotPair) error {
	// TODO
	return nil
}

func (d *dagPoolService) BalanceSlots() error {
	// TODO
	return nil
}

func (d *dagPoolService) Status() (*proto.StatusReply, error) {
	list := make([]*proto.DagNodeStatus, 0, len(d.dagNodesMap))
	for _, node := range d.dagNodesMap {
		pairs := node.GetSlotPairs()
		newPairs := make([]*proto.SlotPair, 0, len(pairs))
		for _, p := range pairs {
			newPairs = append(newPairs, &proto.SlotPair{Start: uint32(p.Start), End: uint32(p.End)})
		}
		cfg := node.GetConfig()
		dataNodes := make([]*proto.DataNodeInfo, 0, len(cfg.Nodes))
		for _, nd := range cfg.Nodes {
			state := node.GetDataNodeState(nd.SetIndex)
			dataNodes = append(dataNodes, &proto.DataNodeInfo{
				SetIndex:   int32(nd.SetIndex),
				RpcAddress: nd.RpcAddress,
				State:      &state,
			})
		}
		st := &proto.DagNodeStatus{
			Node: &proto.DagNodeInfo{
				Name:         cfg.Name,
				Nodes:        dataNodes,
				DataBlocks:   int32(cfg.DataBlocks),
				ParityBlocks: int32(cfg.ParityBlocks),
			},
			Pairs: newPairs,
		}
		list = append(list, st)
	}
	return &proto.StatusReply{
		State:    d.state.String(),
		Statuses: list,
	}, nil
}
