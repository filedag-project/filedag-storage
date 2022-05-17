package pool

import (
	beat "github.com/filedag-project/filedag-storage/dag/node/heart_beat"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
)

type NodeRecordSys struct {
	Db *uleveldb.ULevelDB
	RN []RecordNode
}
type RecordNode struct {
	name   string
	status bool
}

const dagPoolRecord = "dagPoolRecord/"

func NewRecordSys(db *uleveldb.ULevelDB) NodeRecordSys {
	return NodeRecordSys{db, nil}
}
func (r *NodeRecordSys) Add(cid string, theNode int64) error {
	return r.Db.Put(dagPoolRecord+cid, theNode)
}
func (r *NodeRecordSys) AddNode(ips []string, name string) error {
	for _, ip := range ips {
		log.Infof("start listen %v", ip)
		//todo add ip
		go beat.StartListen("7373")
	}
	r.RN = append(r.RN, RecordNode{name: name, status: true})
	return nil
}
func (r *NodeRecordSys) Get(cid string) (int64, error) {
	var theNode int64
	err := r.Db.Get(dagPoolRecord+cid, &theNode)
	if err != nil {
		return -1, err
	}
	return theNode, nil
}
