package dnm

import (
	"context"
	"github.com/filedag-project/filedag-storage/dag/node"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	logging "github.com/ipfs/go-log/v2"
	"golang.org/x/xerrors"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"strconv"
	"sync"
	"time"
)

//NodeRecordSys is a struct for record the dag node
type NodeRecordSys struct {
	Db       *uleveldb.ULevelDB
	RN       map[string]*dagNodeInfo
	NodeLock sync.Mutex
}
type dagNodeInfo struct {
	status       bool
	dataNodeInfo map[string]*dataNodeInfo
}
type dataNodeInfo struct {
	name   string
	status bool
	ip     string
	port   string
}

const healthCheckService = "grpc.health.v1.Health"
const dagPoolRecord = "dagPoolRecord/"

var log = logging.Logger("data-node-manager")

//NewRecordSys  create a new record system
func NewRecordSys(db *uleveldb.ULevelDB) *NodeRecordSys {
	return &NodeRecordSys{Db: db, RN: make(map[string]*dagNodeInfo)}
}

//Add  a new record
func (r *NodeRecordSys) Add(cid string, name string) error {
	return r.Db.Put(dagPoolRecord+cid, name)
}

//HandleDagNode handle the dag node
func (r *NodeRecordSys) HandleDagNode(cons []node.DataNode, name string) error {
	m := make(map[string]*dataNodeInfo)
	for i, c := range cons {
		var dni = dataNodeInfo{
			name:   strconv.Itoa(i),
			status: true,
			ip:     c.Ip,
			port:   c.Port,
		}
		m[strconv.Itoa(i)] = &dni
	}
	tmp := dagNodeInfo{
		status:       true,
		dataNodeInfo: m,
	}
	r.RN[name] = &tmp
	for i, c := range cons {
		go r.handleConn(&c, name, strconv.Itoa(i))
	}
	return nil
}
func (r *NodeRecordSys) handleConn(c *node.DataNode, name string, dataName string) {

	for {
		r.NodeLock.Lock()
		dni := r.RN[name].dataNodeInfo[dataName]
		log.Debugf("heart")
		check, err := c.HeartClient.Check(context.TODO(), &healthpb.HealthCheckRequest{Service: healthCheckService})
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

//Remove a record
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

//Get a record by cid
func (r *NodeRecordSys) Get(cid string) (string, error) {
	var name string
	err := r.Db.Get(dagPoolRecord+cid, &name)
	if err != nil {
		return "", err
	}
	return name, nil
}

//GetNameUseIp get the name by ip
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

//GetCanUseNode get the can use node
func (r *NodeRecordSys) GetCanUseNode() string {
	for n, st := range r.RN {
		if st.status == true {
			return n
		}
	}
	return ""
}
