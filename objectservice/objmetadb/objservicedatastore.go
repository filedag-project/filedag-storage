package objmetadb

import (
	"context"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// ObjStoreMetaDBAPI object service data store API
type ObjStoreMetaDBAPI interface {
	Close() error
	Put(key string, value interface{}) error
	Get(key string, value interface{}) error
	Delete(key string) error
	NewIterator(slice *util.Range, ro *opt.ReadOptions) iterator.Iterator
	ReadAllChan(ctx context.Context, prefix string, seekKey string) (<-chan *entry, error)
}
