package diskv

import (
	"context"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/bluele/gcache"
	"github.com/filedag-project/filedag-storage/kv"
	"golang.org/x/xerrors"
)

const maxLinkDagSize = 8 << 10
const maxWriteTask = 48
const maxReadTask = 64
const maxCacheDags = 4096
const blockpath = "blocks"

var (
	ErrNotFound        = xerrors.New("diskv: not found")
	ErrUnknowOperation = xerrors.New("diskv: unknow operation")
)

var _ kv.KVDB = (*DisKV)(nil)

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
	Key   string
	Value []byte
	Res   chan *opres
}

type DisKV struct {
	Cfg       *Config
	Ref       *Refdb
	closeChan chan struct{}
	close     func()
	opwchan   chan *op
	oprchan   chan *op
	cache     gcache.Cache
}

func NewDisKV(opts ...Option) (*DisKV, error) {
	cfg := &Config{
		MaxLinkDagSize: maxLinkDagSize,
		Shard:          DefaultShardFun,
		MaxRead:        maxReadTask,
		MaxWrite:       maxWriteTask,
		MaxCacheDags:   maxCacheDags,
	}
	for _, opt := range opts {
		opt(cfg)
	}
	ref, err := NewRefdb(filepath.Join(cfg.Dir, refdb_path))
	if err != nil {
		return nil, err
	}
	var o sync.Once
	closeChan := make(chan struct{})
	kv := &DisKV{
		Cfg:       cfg,
		Ref:       ref,
		closeChan: closeChan,
		close: func() {
			o.Do(func() {
				close(closeChan)
			})
		},
		opwchan: make(chan *op),
		oprchan: make(chan *op),
		cache:   gcache.New(cfg.MaxCacheDags).LRU().Expiration(time.Hour * 3).Build(),
	}
	kv.acceptTasks()

	return kv, nil
}

func (di *DisKV) acceptTasks() {
	// set handlers for write operations
	for i := 0; i < di.Cfg.MaxWrite; i++ {
		go func(kv *DisKV) {
			for {
				select {
				case <-kv.closeChan:
					return
				case opt := <-kv.opwchan:
					switch opt.Type {
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
	// set handlers for read operation
	for i := 0; i < di.Cfg.MaxRead; i++ {
		go func(kv *DisKV) {
			for {
				select {
				case <-kv.closeChan:
					return
				case opt := <-kv.oprchan:
					switch opt.Type {
					case opread:
						di.opread(opt)
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
	// find from cache
	if v, err := di.cache.Get(opt.Key); err == nil {
		opt.Res <- &opres{
			Data: v.([]byte),
		}
		return
	}
	// try find reference from refdb
	ref, err := di.getRef(opt.Key)
	if err != nil {
		opt.Res <- &opres{
			Err: err,
		}
		return
	}
	// link-dag keep data in refdb, no need retrive data from disk
	if ref.Type == RefLink {
		opt.Res <- &opres{
			Data: ref.Data,
		}
		// update cache
		di.cache.Set(opt.Key, ref.Data)
		return
	}
	_, p, err := di.pathByKey(opt.Key)
	if err != nil {
		opt.Res <- &opres{
			Err: err,
		}
		return
	}

	f, err := os.OpenFile(p, os.O_RDONLY, 0644)
	if err != nil {
		opt.Res <- &opres{
			Err: err,
		}
		return
	}
	defer f.Close()
	// wait to get read lock
	err = syscall.Flock(int(f.Fd()), syscall.LOCK_SH)
	if err != nil {
		opt.Res <- &opres{
			Err: err,
		}
		return
	}
	defer syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
	d, err := ioutil.ReadAll(f)
	if err != nil {
		opt.Res <- &opres{
			Err: err,
		}
		return
	}
	opt.Res <- &opres{
		Data: d,
	}
	// update cache
	di.cache.Set(opt.Key, ref.Data)
}

func (di *DisKV) opwrite(opt *op) {
	if len(opt.Value) <= di.Cfg.MaxLinkDagSize {
		opt.Res <- &opres{
			Err: di.putRef(opt.Key, opt.Value, true),
		}
		di.cache.Set(opt.Key, opt.Value)
		return
	}

	par, p, err := di.pathByKey(opt.Key)
	if err != nil {
		opt.Res <- &opres{
			Err: err,
		}
		return
	}
	// make sure parent directions has been created
	err = os.MkdirAll(par, 0755)
	if err != nil {
		opt.Res <- &opres{
			Err: err,
		}
		return
	}
	f, err := os.OpenFile(p, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		opt.Res <- &opres{
			Err: err,
		}
		return
	}
	defer f.Close()
	// wait to get read lock
	err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX)
	if err != nil {
		opt.Res <- &opres{
			Err: err,
		}
		return
	}
	defer syscall.Flock(int(f.Fd()), syscall.LOCK_UN)

	n, err := f.Write(opt.Value)
	if err != nil {
		opt.Res <- &opres{
			Err: err,
		}
		return
	}
	err = f.Truncate(int64(n))
	if err != nil {
		opt.Res <- &opres{
			Err: err,
		}
		return
	}
	err = di.putRef(opt.Key, opt.Value, false)
	if err != nil {
		opt.Res <- &opres{
			Err: err,
		}
	} else {
		opt.Res <- &opres{}
		// update cache
		di.cache.Set(opt.Key, opt.Value)
	}
}

func (di *DisKV) opdelete(opt *op) {
	// try find reference from refdb
	ref, err := di.getRef(opt.Key)
	// we does not has the data
	if err != nil {
		opt.Res <- &opres{}
		fmt.Printf("try to delete unkown data: %s\n", opt.Key)
		return
	}
	// link-dag - delete entry in refdb
	if ref.Type == RefLink {
		opt.Res <- &opres{
			Err: di.Ref.Delete(opt.Key),
		}
		return
	}

	_, p, err := di.pathByKey(opt.Key)
	if err != nil {
		opt.Res <- &opres{
			Err: err,
		}
		return
	}

	f, err := os.OpenFile(p, os.O_WRONLY, 0644)
	if err != nil {
		opt.Res <- &opres{
			Err: err,
		}
		return
	}
	defer f.Close()
	// wait to get write lock
	err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX)
	if err != nil {
		opt.Res <- &opres{
			Err: err,
		}

		return
	}
	defer os.Remove(p)
	defer syscall.Flock(int(f.Fd()), syscall.LOCK_UN)

	err = di.Ref.Delete(opt.Key)
	if err != nil {
		opt.Res <- &opres{
			Err: err,
		}
	} else {
		opt.Res <- &opres{}
		// update cache
		di.cache.Remove(opt.Key)
	}
}

func (di *DisKV) pathByKey(key string) (string, string, error) {
	ppath, p, err := di.Cfg.Shard(key)
	if err != nil {
		return "", "", err
	}
	return filepath.Join(di.Cfg.Dir, blockpath, ppath), filepath.Join(di.Cfg.Dir, blockpath, p), nil
}

// Put - write dag node into repo
//
// we use MaxLinkDagSize to divide dag nodes into two categories:
//    - link-dag which size is smaller than or equal to MaxLinkDagSize
//	  - data-dag which size is bigger than MaxLinkDagSize
// we store link-dag into leveldb only, store data-dag into disk and keep an ref whith leveldb
func (di *DisKV) Put(key string, value []byte) error {
	resc := make(chan *opres)
	di.opwchan <- &op{
		Type:  opwrite,
		Key:   key,
		Value: value,
		Res:   resc,
	}
	res := <-resc

	return res.Err
}

func (di *DisKV) Delete(key string) error {
	resc := make(chan *opres)
	di.opwchan <- &op{
		Type: opdelete,
		Key:  key,
		Res:  resc,
	}
	res := <-resc

	return res.Err
}

func (di *DisKV) Get(key string) ([]byte, error) {
	resc := make(chan *opres)
	di.oprchan <- &op{
		Type: opread,
		Key:  key,
		Res:  resc,
	}
	res := <-resc
	return res.Data, res.Err
}

func (di *DisKV) Size(key string) (int, error) {
	ref, err := di.getRef(key)
	if err != nil {
		return -1, err
	}
	return ref.Size, nil
}

func (di *DisKV) AllKeysChan(ctx context.Context) (<-chan string, error) {
	return di.Ref.AllKeysChan(ctx)
}

func (di *DisKV) Close() error {
	di.close()
	return nil
}

func (di *DisKV) getRef(key string) (*DagRef, error) {
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

func (di *DisKV) putRef(key string, value []byte, keepData bool) error {
	ref := &DagRef{
		Code: crc32.ChecksumIEEE(value),
		Size: len(value),
	}
	if keepData {
		ref.Type = RefLink
		ref.Data = value
	} else {
		ref.Type = RefData
	}
	data, err := ref.Bytes()
	if err != nil {
		return err
	}
	return di.Ref.Put(key, data)
}
