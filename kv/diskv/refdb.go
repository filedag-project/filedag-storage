package diskv

import (
	"context"

	"github.com/fxamacker/cbor/v2"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/iterator"
)

const refdb_path = "refdb"
const (
	RefData int8 = iota
	RefLink
)

type Refdb struct {
	db *leveldb.DB
}

func NewRefdb(dir string) (*Refdb, error) {
	db, err := leveldb.OpenFile(dir, nil)
	if err != nil {
		return nil, err
	}
	return &Refdb{
		db: db,
	}, nil
}

func (ref *Refdb) Put(key string, value []byte) error {
	return ref.db.Put([]byte(key), value, nil)
}

func (ref *Refdb) Get(key string) ([]byte, error) {
	bs, err := ref.db.Get([]byte(key), nil)
	if err == nil {
		return bs, nil
	}
	if err == errors.ErrNotFound {
		return nil, ErrNotFound
	}
	return nil, err
}

func (ref *Refdb) Delete(key string) error {
	return ref.db.Delete([]byte(key), nil)
}

func (ref *Refdb) AllKeysChan(ctx context.Context) (chan string, error) {
	iter := ref.db.NewIterator(nil, nil)
	out := make(chan string, 1)
	go func(iter iterator.Iterator, oc chan string) {
		defer iter.Release()
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			if !iter.Next() {
				return
			}
			out <- string(iter.Key())
		}
		// Todo: log if has iter.Error()
	}(iter, out)
	return out, nil
}

func (ref *Refdb) Close() error {
	return ref.db.Close()
}

type DagRef struct {
	Code uint32 // crc32 checksum
	Size int
	Type int8
	Data []byte
}

func (d *DagRef) Bytes() ([]byte, error) {
	return cbor.Marshal(d)
}

func (d *DagRef) FromBytes(data []byte) error {
	return cbor.Unmarshal(data, d)
}
