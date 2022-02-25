package diskv

import (
	"fmt"
	"hash/crc32"
	"path/filepath"

	"golang.org/x/xerrors"
)

const maxLinkDagSize = 8 << 10
const paralletTask = 4

var (
	ErrNotFound        = xerrors.New("diskv: not found")
	ErrUnknowOperation = xerrors.New("diskv: unknow operation")
)

type optype int8

const (
	opread optype = iota
	opwrite
	opdelete
)

func (o optype) String() string {
	switch o {
	case opread:
		return "opread"
	case opwrite:
		return "opwrite"
	case opdelete:
		return "opdelete"
	default:
		return "unknow"
	}
}

type opres struct {
	Err  error
	Data []byte
}

type op struct {
	Type  optype
	Key   []byte
	Value []byte
	Res   chan *opres
}

type DisKV struct {
	Cfg    *Config
	Ref    *Refdb
	close  chan struct{}
	opchan chan *op
}

func NewDisKV(opts ...Option) (*DisKV, error) {
	cfg := &Config{
		MaxLinkDagSize: maxLinkDagSize,
		Shard:          DefaultShardFun,
		Parallel:       paralletTask,
	}
	for _, opt := range opts {
		opt(cfg)
	}
	ref, err := NewRefdb(filepath.Join(cfg.Dir, refdb_path))
	if err != nil {
		return nil, err
	}
	kv := &DisKV{
		Cfg:    cfg,
		Ref:    ref,
		close:  make(chan struct{}),
		opchan: make(chan *op),
	}
	kv.acceptTasks()
	return kv, nil
}

func (di *DisKV) acceptTasks() {
	for i := 0; i < di.Cfg.Parallel; i++ {
		go func(kv *DisKV) {
			for {
				select {
				case <-kv.close:
					return
				case opt := <-kv.opchan:
					fmt.Printf("%s %s %s\n", opt.Type, opt.Key, opt.Value)
					switch opt.Type {
					case opread:
						di.opread(opt)
					case opwrite:
						di.opwrite(opt)
					case opdelete:
						di.opdelete(opt)
					default:
						opt.Res <- &opres{
							Err: ErrUnknowOperation,
						}
					}
				}
			}
		}(di)
	}
}

func (di *DisKV) opread(opt *op) {

}

func (di *DisKV) opwrite(opt *op) {

}

func (di *DisKV) opdelete(opt *op) {

}

// Put - write dag node into repo
//
// we use MaxLinkDagSize to divide dag nodes into two categories:
//    - link-dag which size is smaller than or equal to MaxLinkDagSize
//	  - data-dag which size is bigger than MaxLinkDagSize
// we store link-dag into leveldb only, store data-dag into disk and keep an ref whith leveldb
func (di *DisKV) Put(key []byte, value []byte) error {
	// vsz := len(value)
	// if vsz <= di.Cfg.MaxLinkDagSize {
	// 	return di.putRef(key, value, true)
	// }
	resc := make(chan *opres)
	di.opchan <- &op{
		Type:  opwrite,
		Key:   key,
		Value: value,
		Res:   resc,
	}
	res := <-resc

	return res.Err
}

func (di *DisKV) Delete(key []byte) error {
	resc := make(chan *opres)
	di.opchan <- &op{
		Type: opdelete,
		Key:  key,
		Res:  resc,
	}
	res := <-resc

	return res.Err
}

func (di *DisKV) Get(key []byte) ([]byte, error) {
	// ref, err := di.getRef(key)
	// if err != nil {
	// 	return nil, err
	// }
	// if ref.Type == RefLink {
	// 	return ref.Data, nil
	// }
	resc := make(chan *opres)
	di.opchan <- &op{
		Type: opread,
		Key:  key,
		Res:  resc,
	}
	res := <-resc
	return res.Data, res.Err
}

func (di *DisKV) Size(key []byte) (int, error) {
	ref, err := di.getRef(key)
	if err != nil {
		return -1, err
	}
	return ref.Size, nil
}

func (di *DisKV) Close() error {
	close(di.close)
	return nil
}

func (di *DisKV) getRef(key []byte) (*DagRef, error) {
	data, err := di.Ref.Get(key)
	if err != nil {
		return nil, err
	}
	ref := &DagRef{}
	err = ref.FromBytes(data)
	if err != nil {
		return nil, err
	}
	return ref, nil
}

func (di *DisKV) putRef(key []byte, value []byte, keepData bool) error {
	ref := &DagRef{
		Code: crc32.ChecksumIEEE(value),
		Size: len(value),
		Type: RefLink,
	}
	if keepData {
		ref.Data = value
	}
	data, err := ref.Bytes()
	if err != nil {
		return err
	}
	return di.Ref.Put(key, data)
}

type Config struct {
	Dir            string
	MaxLinkDagSize int
	Shard          ShardFun
	Parallel       int
}

type Option func(cfg *Config)
