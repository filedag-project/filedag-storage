package pool

import (
	"context"
	"github.com/filedag-project/filedag-storage/dag/node"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	"golang.org/x/xerrors"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"strconv"
	"sync"
	"time"
)

type NodeRecordSys struct {
	Db       *uleveldb.ULevelDB
	RN       map[string]*DagNodeInfo
	NodeLock sync.Mutex
}
type DagNodeInfo struct {
	status       bool
	dataNodeInfo map[string]*DataNodeInfo
}
type DataNodeInfo struct {
	name   string
	status bool
	ip     string
	port   string
}

const HealthCheckService = "grpc.health.v1.Health"
const dagPoolRecord = "dagPoolRecord/"

func NewRecordSys(db *uleveldb.ULevelDB) NodeRecordSys {
	return NodeRecordSys{Db: db, RN: make(map[string]*DagNodeInfo)}
}
func (r *NodeRecordSys) Add(cid string, name string) error {
	return r.Db.Put(dagPoolRecord+cid, name)
}
func (r *NodeRecordSys) HandleDagNode(cons []node.DataNode, name string) error {
	m := make(map[string]*DataNodeInfo)
	for i, c := range cons {
		var dni = DataNodeInfo{
			name:   strconv.Itoa(i),
			status: true,
			ip:     c.Ip,
			port:   c.Port,
		}
		m[strconv.Itoa(i)] = &dni
		go r.HandleConn(&c, name, dni.name)
	}
	tmp := DagNodeInfo{
		status:       true,
		dataNodeInfo: m,
	}
	r.RN[name] = &tmp
	return nil
}
func (r *NodeRecordSys) HandleConn(c *node.DataNode, name string, dataName string) {

	for {
		r.NodeLock.Lock()
		dni := r.RN[name].dataNodeInfo[dataName]
		log.Infof("heart")
		check, err := c.HeartClient.Check(context.TODO(), &healthpb.HealthCheckRequest{Service: HealthCheckService})
		if err != nil {
			log.Errorf("Check the %v ip:%v,port:%v err:%v", dni.name, dni.ip, dni.port, err)
			r.Remove(name, dataName)
			return
		}
		if check.Status != healthpb.HealthCheckResponse_SERVING {
			log.Errorf("the %v ip:%v,port:%v not ser", dni.name, dni.ip, dni.port)
			r.Remove(name, dataName)
		}
		r.NodeLock.Unlock()
		time.Sleep(time.Second * 10)
	}
}
func (r *NodeRecordSys) Remove(name, dataName string) {
	r.RN[name].dataNodeInfo[dataName].status = false
	count := 0
	for _, info := range r.RN[name].dataNodeInfo {
		if !info.status {
			count++
		}
		if count >= 2 {
			r.RN[name].status = false
			break
		}
	}
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
		for _, i := range n.dataNodeInfo {
			if i.ip == ip {
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
