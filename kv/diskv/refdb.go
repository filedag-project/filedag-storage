package diskv

import (
	"encoding/json"

	"github.com/syndtr/goleveldb/leveldb"
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

func (ref *Refdb) Put(key []byte, value []byte) error {
	return ref.db.Put(key, value, nil)
}

func (ref *Refdb) Get(key []byte) ([]byte, error) {
	return ref.db.Get(key, nil)
}

func (ref *Refdb) Delete(key []byte) error {
	return ref.db.Delete(key, nil)
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
	return json.Marshal(d)
}

func (d *DagRef) FromBytes(data []byte) error {
	return json.Unmarshal(data, d)
}
