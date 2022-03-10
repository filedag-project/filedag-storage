package uleveldb

import (
	"encoding/json"
	logging "github.com/ipfs/go-log/v2"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
	"strings"
)

var log = logging.Logger("leveldb")

// DBClient start a levelDb
var DBClient *ULevelDB

//ULevelDB level db store key-struct
type ULevelDB struct {
	DB *leveldb.DB
}

// OpenDb open a db client
func OpenDb(path string) *ULevelDB {
	newDb, err := leveldb.OpenFile(path, nil)
	if _, corrupted := err.(*errors.ErrCorrupted); corrupted {
		newDb, err = leveldb.RecoverFile(path, nil)
	}
	if err != nil {
		log.Errorf("Open Db path: %v,err:%v,", path, err)
	}
	uLevelDB := ULevelDB{}
	uLevelDB.DB = newDb
	return &uLevelDB
}

//Close db close
func (l *ULevelDB) Close() {
	l.DB.Close()
}

// Put
// * @param {interface{}} key
// * @param {interface{}} value
func (l *ULevelDB) Put(key string, value interface{}) error {

	result, err := json.Marshal(value)
	if err != nil {
		log.Errorf("marshal error%v", err)
		return err
	}
	err = l.DB.Put([]byte(key), result, nil)
	return err
}

// Get
// * @param {interface{}} key
// * @param {interface{}} value
func (l *ULevelDB) Get(key, value interface{}) error {
	get, err := l.DB.Get([]byte(key.(string)), nil)
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
func (l *ULevelDB) Delete(key string) error {

	return l.DB.Delete([]byte(key), nil)
}

// NewIterator /**
func (l *ULevelDB) NewIterator(slice *util.Range, ro *opt.ReadOptions) iterator.Iterator {

	return l.DB.NewIterator(slice, ro)
}

//ReadAll read all key value
func (l *ULevelDB) ReadAll(prefix string) (map[string][]byte, error) {
	iter := l.NewIterator(nil, nil)
	m := make(map[string][]byte)
	for iter.Next() {
		if strings.Contains(string(iter.Key()), prefix) {
			m[string(iter.Key())] = iter.Value()
		}
	}
	iter.Release()
	return m, nil
}
