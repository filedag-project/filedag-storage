package pool

import (
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	"sync"
)

type NodeRecordSys struct {
	Db       *uleveldb.ULevelDB
	RN       map[string]bool
	NodeLock sync.Mutex
}

const dagPoolRecord = "dagPoolRecord/"

func NewRecordSys(db *uleveldb.ULevelDB) NodeRecordSys {
	return NodeRecordSys{Db: db, RN: make(map[string]bool)}
}
func (r *NodeRecordSys) Add(cid string, name string) error {
	return r.Db.Put(dagPoolRecord+cid, name)
}
func (r *NodeRecordSys) HandleDagNode(ips []string, name string) error {
	for _, ip := range ips {
		log.Infof("start listen %v", ip)
		//todo add ip
		go r.StartListen(":7373", name)
	}
	r.RN[name] = true
	return nil
}
func (r *NodeRecordSys) Remove(name string) {
	r.NodeLock.Lock()
	r.RN[name] = false
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
func (r *NodeRecordSys) GetCanUseNode() string {
	for n, st := range r.RN {
		if st == true {
			return n
		}
	}
	return ""
}
