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

type Uleveldb struct {
	DB *leveldb.DB
}

const (
	DBFILE = "/tmp/leveldb2.db"
)

//GlobalLevelDB global LevelDB
var GlobalLevelDB = OpenDb(DBFILE)

func OpenDb(path string) *Uleveldb {
	newdb, err := leveldb.OpenFile(path, nil)
	if _, corrupted := err.(*errors.ErrCorrupted); corrupted {
		newdb, err = leveldb.RecoverFile(path, nil)
	}

	if err != nil {
		panic(err)
	}
	uleveldb := Uleveldb{}
	uleveldb.DB = newdb
	return &uleveldb
}

func (uleveldb *Uleveldb) Close() {
	uleveldb.DB.Close()
}

// Put
// * @param {interface{}} key
// * @param {interface{}} value
func (uleveldb *Uleveldb) Put(key string, value interface{}) error {

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
func (uleveldb *Uleveldb) Get(key, value interface{}) error {
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
func (uleveldb *Uleveldb) Delete(key string) error {

	return uleveldb.DB.Delete([]byte(key), nil)
}

// NewIterator /**
func (uleveldb *Uleveldb) NewIterator(slice *util.Range, ro *opt.ReadOptions) iterator.Iterator {

	return uleveldb.DB.NewIterator(slice, ro)
}

//ReadAll read all key value
func (uleveldb *Uleveldb) ReadAll() (map[string][]byte, error) {
	iter := GlobalLevelDB.NewIterator(nil, nil)
	m := make(map[string][]byte)
	for iter.Next() {
		m[string(iter.Key())] = iter.Value()
	}
	iter.Release()
	return m, nil
}
