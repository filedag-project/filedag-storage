package poolservice

import (
	"context"
	"github.com/filedag-project/filedag-storage/dag/node/dagnode"
)

func (d *dagPoolService) RepairDataNode(ctx context.Context, dagNodeName string, fromNodeIndex int, repairNodeIndex int) error {
	node, ok := func() (*dagnode.DagNode, bool) {
		d.dagNodesLock.RLock()
		defer d.dagNodesLock.RUnlock()
		nd, ok := d.dagNodesMap[dagNodeName]
		return nd, ok
	}()
	if !ok {
		return ErrDagNodeNotFound
	}

	return node.RepairDataNode(ctx, fromNodeIndex, repairNodeIndex)
}
