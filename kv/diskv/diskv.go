package diskv

import (
	"hash/crc32"
	"path/filepath"

	"golang.org/x/xerrors"
)

var ErrNotFound = xerrors.New("diskv: not found")

type DisKV struct {
	Cfg *Config
	Ref *Refdb
}

func NewDisKV(opts ...Option) (*DisKV, error) {
	cfg := &Config{}
	for _, opt := range opts {
		opt(cfg)
	}
	ref, err := NewRefdb(filepath.Join(cfg.Dir, refdb_path))
	if err != nil {
		return nil, err
	}
	return &DisKV{
		Cfg: cfg,
		Ref: ref,
	}, nil
}

// Put - write dag node into repo
//
// we use MaxLinkDagSize to divide dag nodes into two categories:
//    - link-dag which size is smaller than or equal to MaxLinkDagSize
//	  - data-dag which size is bigger than MaxLinkDagSize
// we store link-dag into leveldb only, store data-dag into disk and keep an ref whith leveldb
func (di *DisKV) Put(key []byte, value []byte) error {
	vsz := len(value)
	if vsz <= di.Cfg.MaxLinkDagSize {
		return di.putlink(key, value)
	}

	return nil
}

func (di *DisKV) Delete(key []byte) error {
	return nil
}

func (di *DisKV) Get(key []byte) ([]byte, error) {
	ref, err := di.getRef(key)
	if err != nil {
		return nil, err
	}
	if ref.Type == RefLink {
		return ref.Data, nil
	}

	return nil, nil
}

func (di *DisKV) Size(key []byte) (int, error) {
	ref, err := di.getRef(key)
	if err != nil {
		return -1, err
	}
	return ref.Size, nil
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

func (di *DisKV) putlink(key []byte, value []byte) error {
	ref := &DagRef{
		Code: crc32.ChecksumIEEE(value),
		Size: len(value),
		Type: RefLink,
		Data: value,
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
}

type Option func(cfg *Config)
