package pool

import (
	"github.com/filedag-project/filedag-storage/dag/config"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	"golang.org/x/xerrors"
	"sync"
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

const dagPoolRecord = "dagPoolRecord/"

func NewRecordSys(db *uleveldb.ULevelDB) NodeRecordSys {
	return NodeRecordSys{Db: db, RN: make(map[string]DagNodeInfo)}
}
func (r *NodeRecordSys) Add(cid string, name string) error {
	return r.Db.Put(dagPoolRecord+cid, name)
}
func (r *NodeRecordSys) HandleDagNode(cons []config.CaskConfig, name string) error {
	var ips []string
	for _, c := range cons {
		log.Infof("start listen heartbeat %v", c.HeartAddr)
		go r.StartListen(c.HeartAddr, name)
		ips = append(ips, c.Ip+c.Port)
	}
	tmp := DagNodeInfo{true, ips}
	r.RN[name] = tmp
	return nil
}
func (r *NodeRecordSys) Remove(name string) {
	r.NodeLock.Lock()
	tmp := DagNodeInfo{true, r.RN[name].ips}
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
