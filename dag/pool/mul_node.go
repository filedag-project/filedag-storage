package pool

import (
	"context"
	"github.com/filedag-project/filedag-storage/dag/node"
	"github.com/ipfs/go-cid"
)

// GetNode get the DagNode
func (d *DagPool) GetNode(ctx context.Context, c cid.Cid) (*node.DagNode, error) {
	//todo mul node
	get, err := d.NRSys.Get(c.String())
	if err != nil {
		return nil, err
	}
	return d.DagNodes[get], nil
}

// UseNode get the DagNode
func (d *DagPool) UseNode(ctx context.Context, c cid.Cid) (*node.DagNode, error) {
	//todo mul node
	dn := d.NRSys.GetCanUseNode()
	err := d.NRSys.Add(c.String(), dn)
	if err != nil {
		return nil, err
	}
	return d.DagNodes[dn], nil
}

// GetNodeUseIP get the DagNode
func (d *DagPool) GetNodeUseIP(ctx context.Context, ip string) (*node.DagNode, error) {
	//todo mul node
	get, err := d.NRSys.GetNameUseIp(ip)
	if err != nil {
		return nil, err
	}
	return d.DagNodes[get], nil
}

//
//// GetNodes get the DagNode
//func (d *DagPool) GetNodes(ctx context.Context, cids []cid.Cid) map[*node.DagNode][]cid.Cid {
//	//todo mul node
//	//
//	m := make(map[*node.DagNode][]cid.Cid)
//	for _, c := range cids {
//		get, err := d.NRSys.Get(c.String())
//		if err != nil {
//			return nil
//		}
//		m[d.DagNodes[get]] = append(m[d.DagNodes[get]], c)
//	}
//	return m
//}
//
//// UseNodes get the DagNode
//func (d *DagPool) UseNodes(ctx context.Context, c []cid.Cid) *node.DagNode {
//	//todo mul node
//	dn := d.NRSys.GetCanUseNode()
//	err := d.NRSys.Add(c[0].String(), dn)
//	if err != nil {
//		return nil
//	}
//	return d.DagNodes[dn]
//}
