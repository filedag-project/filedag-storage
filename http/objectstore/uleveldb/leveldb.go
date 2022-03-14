package uleveldb

import (
	"encoding/json"
	logging "github.com/ipfs/go-log/v2"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var log = logging.Logger("leveldb")

//ULeveldb level db store key-struct
type ULeveldb struct {
	DB *leveldb.DB
}

const (
	userDBFILE       = "/tmp/leveldb/user.db"
	policyDBFILE     = "/tmp/leveldb/policy.db"
	userPolicyDBFILE = "/tmp/leveldb/user-policy.db"
	groupDBFILE      = "/tmp/leveldb/group.db"
)

//GlobalUserLevelDB global LevelDB
var (
	GlobalUserLevelDB       = OpenDb(userDBFILE)
	GlobalPolicyLevelDB     = OpenDb(policyDBFILE)
	GlobalUserPolicyLevelDB = OpenDb(userPolicyDBFILE)
	GlobalGroupLevelDB      = OpenDb(groupDBFILE)
)

func OpenDb(path string) *ULeveldb {
	newdb, err := leveldb.OpenFile(path, nil)
	if _, corrupted := err.(*errors.ErrCorrupted); corrupted {
		newdb, err = leveldb.RecoverFile(path, nil)
	}
	if err != nil {
		log.Errorf("Open Db%v", err)
	}
	uleveldb := ULeveldb{}
	uleveldb.DB = newdb
	return &uleveldb
}

func (uleveldb *ULeveldb) Close() {
	uleveldb.DB.Close()
}

// Put
// * @param {interface{}} key
// * @param {interface{}} value
func (uleveldb *ULeveldb) Put(key string, value interface{}) error {

	result, err := json.Marshal(value)
	if err != nil {
		log.Errorf("marshal error%v", err)
		return err
	}
	err = uleveldb.DB.Put([]byte(key), result, nil)
	return err
}

// Get
// * @param {interface{}} key
// * @param {interface{}} value
func (uleveldb *ULeveldb) Get(key, value interface{}) error {
	get, err := uleveldb.DB.Get([]byte(key.(string)), nil)
	if err != nil {
		log.Errorf(" Get error%v", err)
		return err
	}
	err = json.Unmarshal(get, value)
	if err != nil {
		return err
	}
	return err
}

// Delete
// * @param {interface{}} key
// * @param {interface{}} value
func (uleveldb *ULeveldb) Delete(key string) error {

	return uleveldb.DB.Delete([]byte(key), nil)
}

// NewIterator /**
func (uleveldb *ULeveldb) NewIterator(slice *util.Range, ro *opt.ReadOptions) iterator.Iterator {

	return uleveldb.DB.NewIterator(slice, ro)
}

//ReadAll read all key value
func (uleveldb *ULeveldb) ReadAll() (map[string][]byte, error) {
	iter := GlobalUserLevelDB.NewIterator(nil, nil)
	m := make(map[string][]byte)
	for iter.Next() {
		m[string(iter.Key())] = iter.Value()
	}
	iter.Release()
	return m, nil
}
