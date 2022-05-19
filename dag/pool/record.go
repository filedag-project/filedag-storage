package pool

import (
	"context"
	"github.com/filedag-project/filedag-storage/dag/node"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	"golang.org/x/xerrors"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"sync"
	"time"
)

type NodeRecordSys struct {
	Db       *uleveldb.ULevelDB
	RN       map[string]DagNodeInfo
	NodeLock sync.Mutex
}
type DagNodeInfo struct {
	status bool
	ips    []string
}

const HealthCheckService = "grpc.health.v1.Health"
const dagPoolRecord = "dagPoolRecord/"

func NewRecordSys(db *uleveldb.ULevelDB) NodeRecordSys {
	return NodeRecordSys{Db: db, RN: make(map[string]DagNodeInfo)}
}
func (r *NodeRecordSys) Add(cid string, name string) error {
	return r.Db.Put(dagPoolRecord+cid, name)
}
func (r *NodeRecordSys) HandleDagNode(cons []node.DataNode, name string) error {
	var ips []string
	for _, c := range cons {
		ips = append(ips, c.Ip)
		go r.HandleConn(&c, name)
	}
	tmp := DagNodeInfo{true, ips}
	r.RN[name] = tmp
	return nil
}
func (r *NodeRecordSys) HandleConn(c *node.DataNode, name string) {
	for {
		log.Infof("aaa")
		watch, err := c.HeartClient.Watch(context.TODO(), &healthpb.HealthCheckRequest{Service: HealthCheckService})
		if err != nil {
			log.Errorf("watch err:%v", err)
			r.Remove(name)
			return
		}
		recv, err := watch.Recv()
		if err != nil {
			log.Errorf("Recv err:%v", err)
			r.Remove(name)
			return
		}
		if recv.Status != healthpb.HealthCheckResponse_SERVING {
			log.Errorf("not ser")
			r.Remove(name)
		}
		time.Sleep(time.Second * 2)
	}
}
func (r *NodeRecordSys) Remove(name string) {
	r.NodeLock.Lock()
	tmp := DagNodeInfo{false, r.RN[name].ips}
	r.RN[name] = tmp
	r.NodeLock.Unlock()
}

func (r *NodeRecordSys) Get(cid string) (string, error) {
	var name string
	err := r.Db.Get(dagPoolRecord+cid, &name)
	if err != nil {
		return "", err
	}
	return name, nil
}
func (r *NodeRecordSys) GetNameUseIp(ip string) (string, error) {
	for name, n := range r.RN {
		for _, i := range n.ips {
			if i == ip {
				return name, nil
			}
		}
	}
	return "", xerrors.Errorf("no dagnode")
}
func (r *NodeRecordSys) GetCanUseNode() string {
	for n, st := range r.RN {
		if st.status == true {
			return n
		}
	}
	return ""
}
