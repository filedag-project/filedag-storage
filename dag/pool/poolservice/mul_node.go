package poolservice

import (
	"context"
	"github.com/filedag-project/filedag-storage/dag/node"
	"github.com/ipfs/go-cid"
)

// GetNode get the DagNode
func (d *dagPoolService) GetNode(ctx context.Context, c cid.Cid) (*node.DagNode, error) {
	//todo mul node
	get, err := d.nrSys.Get(c.String())
	if err != nil {
		return nil, err
	}
	return d.dagNodes[get], nil
}

// UseNode get the DagNode
func (d *dagPoolService) UseNode(ctx context.Context, c cid.Cid) (*node.DagNode, error) {
	//todo mul node
	dn := d.nrSys.GetCanUseNode()
	err := d.nrSys.Add(c.String(), dn)
	if err != nil {
		return nil, err
	}
	return d.dagNodes[dn], nil
}

// GetNodeUseIP get the DagNode
func (d *dagPoolService) GetNodeUseIP(ctx context.Context, ip string) (*node.DagNode, error) {
	//todo mul node
	get, err := d.nrSys.GetNameUseIp(ip)
	if err != nil {
		return nil, err
	}
	return d.dagNodes[get], nil
}

//
//// GetNodes get the DagNode
//func (d *DagPoolService) GetNodes(ctx context.Context, cids []cid.Cid) map[*node.DagNode][]cid.Cid {
//	//todo mul node
//	//
//	m := make(map[*node.DagNode][]cid.Cid)
//	for _, c := range cids {
//		get, err := d.nrSys.Get(c.String())
//		if err != nil {
//			return nil
//		}
//		m[d.dagNodes[get]] = append(m[d.dagNodes[get]], c)
//	}
//	return m
//}
//
//// UseNodes get the DagNode
//func (d *DagPoolService) UseNodes(ctx context.Context, c []cid.Cid) *node.DagNode {
//	//todo mul node
//	dn := d.nrSys.GetCanUseNode()
//	err := d.nrSys.Add(c[0].String(), dn)
//	if err != nil {
//		return nil
//	}
//	return d.dagNodes[dn]
//}
