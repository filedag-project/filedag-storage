package poolservice

import (
	"context"
	"github.com/filedag-project/filedag-storage/dag/node/dagnode"
	"github.com/ipfs/go-cid"
)

// GetNode get the DagNode
func (d *dagPoolService) GetNode(ctx context.Context, c cid.Cid) (*dagnode.DagNode, error) {
	//todo mul node
	get, err := d.nrSys.Get(c.String())
	if err != nil {
		return nil, err
	}
	return d.dagNodes[get], nil
}

// UseNode get the DagNode
func (d *dagPoolService) UseNode(ctx context.Context, c cid.Cid) (*dagnode.DagNode, error) {
	//todo mul node
	dn := d.nrSys.GetCanUseNode()
	err := d.nrSys.Add(c.String(), dn)
	if err != nil {
		return nil, err
	}
	return d.dagNodes[dn], nil
}

// getNodeUseIP get the DagNode
func (d *dagPoolService) getNodeUseIP(ctx context.Context, ip string) (*dagnode.DagNode, error) {
	//todo mul node
	get, err := d.nrSys.GetNameUseIp(ip)
	if err != nil {
		return nil, err
	}
	return d.dagNodes[get], nil
}

//
//// GetNodes get the DagNode
//func (d *Pool) GetNodes(ctx context.Context, cids []cid.Cid) map[*node.DagNode][]cid.Cid {
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
//func (d *Pool) UseNodes(ctx context.Context, c []cid.Cid) *node.DagNode {
//	//todo mul node
//	dn := d.NRSys.GetCanUseNode()
//	err := d.NRSys.Add(c[0].String(), dn)
//	if err != nil {
//		return nil
//	}
//	return d.DagNodes[dn]
//}
