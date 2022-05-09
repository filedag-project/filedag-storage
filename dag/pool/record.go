package pool

import "github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"

type RecordSys struct {
	Db *uleveldb.ULevelDB
}

const dagPoolRecord = "dagPoolRecord/"

func NewRecordSys(db *uleveldb.ULevelDB) RecordSys {
	return RecordSys{db}
}
func (r *RecordSys) Add(cid string, theNode int64) error {
	return r.Db.Put(dagPoolRecord+cid, theNode)
}
func (r *RecordSys) Get(cid string) (int64, error) {
	var theNode int64
	err := r.Db.Get(dagPoolRecord+cid, &theNode)
	if err != nil {
		return -1, err
	}
	return theNode, nil
}
