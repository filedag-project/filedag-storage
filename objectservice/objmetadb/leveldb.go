package objmetadb

import (
	"context"
	logging "github.com/ipfs/go-log/v2"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/vmihailenco/msgpack/v5"
	"go.uber.org/zap/buffer"
)

var log = logging.Logger("leveldb")
var _ ObjStoreMetaDBAPI = &objStoreMetaData{}

//objStoreMetaData level db store key-struct
type objStoreMetaData struct {
	DB *leveldb.DB
}

// OpenDb open a db client
func OpenDb(path string) (*objStoreMetaData, error) {
	newDb, err := leveldb.OpenFile(path, nil)
	if _, corrupted := err.(*errors.ErrCorrupted); corrupted {
		newDb, err = leveldb.RecoverFile(path, nil)
	}
	if err != nil {
		log.Errorf("Open Db path: %v,err:%v,", path, err)
		return nil, err
	}
	return &objStoreMetaData{
		DB: newDb,
	}, nil
}

//Close db close
func (l *objStoreMetaData) Close() error {
	return l.DB.Close()
}

// Put
// * @param {string} key
// * @param {interface{}} value
func (l *objStoreMetaData) Put(key string, value interface{}) error {
	result, err := msgpack.Marshal(value)
	if err != nil {
		log.Errorf("marshal error%v", err)
		return err
	}
	return l.DB.Put([]byte(key), result, nil)
}

// Get
// * @param {string} key
// * @param {interface{}} value
func (l *objStoreMetaData) Get(key string, value interface{}) error {
	get, err := l.DB.Get([]byte(key), nil)
	if err != nil {
		return err
	}
	return msgpack.Unmarshal(get, value)
}

// Delete
// * @param {string} key
// * @param {interface{}} value
func (l *objStoreMetaData) Delete(key string) error {
	return l.DB.Delete([]byte(key), nil)
}

// NewIterator /**
func (l *objStoreMetaData) NewIterator(slice *util.Range, ro *opt.ReadOptions) iterator.Iterator {
	return l.DB.NewIterator(slice, ro)
}

type entry struct {
	Key   string
	Value []byte
}

func (e *entry) UnmarshalValue(value interface{}) error {
	return msgpack.Unmarshal(e.Value, value)
}
func (e *entry) GetKey() string {
	return e.Key
}

//ReadAllChan read all key value
func (l *objStoreMetaData) ReadAllChan(ctx context.Context, prefix string, seekKey string) (<-chan *entry, error) {
	ch := make(chan *entry)
	var slice *util.Range
	if prefix != "" {
		slice = util.BytesPrefix([]byte(prefix))
	}
	iter := l.NewIterator(slice, nil)
	if seekKey != "" {
		iter.Seek([]byte(seekKey))
	}
	go func() {
		defer func() {
			iter.Release()
			close(ch)
		}()
		for iter.Next() {
			key := string(iter.Key())
			buf := buffer.Buffer{}
			buf.Write(iter.Value())
			select {
			case <-ctx.Done():
				return
			case ch <- &entry{
				Key:   key,
				Value: buf.Bytes(),
			}:
			}
		}
	}()
	return ch, nil
}
