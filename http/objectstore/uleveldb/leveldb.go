package uleveldb

import (
	"encoding/json"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type Uleveldb struct {
	DB *leveldb.DB
}

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
		fmt.Println("error")
		return err
	}
	err = uleveldb.DB.Put([]byte(key), []byte(result), nil)
	return err
}

// Get
// * @param {interface{}} key
// * @param {interface{}} value
func (uleveldb *Uleveldb) Get(key interface{}) ([]byte, error) {

	return uleveldb.DB.Get([]byte(key.(string)), nil)
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
