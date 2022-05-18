package pool

import (
	"context"
	"github.com/filedag-project/filedag-storage/dag/pool/userpolicy"
	bserv "github.com/ipfs/go-blockservice"
	"github.com/ipfs/go-cid"
	"net"
	"os"
	"strings"
	"time"
)

// CheckPolicy check user policy
func (d *DagPool) CheckPolicy(ctx context.Context, policy userpolicy.DagPoolPolicy) bool {
	s := strings.Split((ctx.Value("user")).(string), ",")
	if len(s) != 2 {
		return false
	}
	return d.Iam.CheckUserPolicy(s[0], s[1], policy)
}

// GetNode get the DagNode
func (d *DagPool) GetNode(ctx context.Context, c cid.Cid) bserv.BlockService {
	//todo mul node
	get, err := d.NRSys.Get(c.String())
	if err != nil {
		return nil
	}
	return d.Blocks[get]
}

// UseNode get the DagNode
func (d *DagPool) UseNode(ctx context.Context, c cid.Cid) bserv.BlockService {
	//todo mul node
	dn := d.NRSys.GetCanUseNode()
	err := d.NRSys.Add(c.String(), dn)
	if err != nil {
		return nil
	}
	return d.Blocks[dn]
}

// GetNodes get the DagNode
func (d *DagPool) GetNodes(ctx context.Context, cids []cid.Cid) map[bserv.BlockService][]cid.Cid {
	//todo mul node
	//
	m := make(map[bserv.BlockService][]cid.Cid)
	for _, c := range cids {
		get, err := d.NRSys.Get(c.String())
		if err != nil {
			return nil
		}
		m[d.Blocks[get]] = append(m[d.Blocks[get]], c)
	}
	return m
}

// UseNodes get the DagNode
func (d *DagPool) UseNodes(ctx context.Context, c []cid.Cid) bserv.BlockService {
	//todo mul node
	dn := d.NRSys.GetCanUseNode()
	err := d.NRSys.Add(c[0].String(), dn)
	if err != nil {
		return nil
	}
	return d.Blocks[dn]
}
func (r *NodeRecordSys) StartListen(addr, name string) {
	netListen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Errorf("connect error:%v", err)
		os.Exit(1)
	}
	log.Infof("Waiting for Client ...")
	for {
		conn, err := netListen.Accept()
		if err != nil {
			log.Errorf(conn.RemoteAddr().String(), "Fatal error: ", err)
			continue
		}
		conn.SetReadDeadline(time.Now().Add(time.Duration(10) * time.Second))

		log.Infof("%v,%v", conn.RemoteAddr().String(), "connect success!")
		go r.HandleConnection(conn, name)
	}
}
