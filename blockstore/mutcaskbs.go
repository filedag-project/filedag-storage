package blockstore

import (
	"github.com/filedag-project/filedag-storage/kv/mutcask"
)

type Config struct {
	Batch      int
	HashOnRead bool
	Path       string
	CaskNum    int
}

func NewMutcaskbs(cfg *Config) (*blostore, error) {
	if cfg.Batch == 0 {
		cfg.Batch = default_batch_num
	}
	if cfg.CaskNum == 0 {
		cfg.CaskNum = default_cask_num
	}
	mc, err := mutcask.NewMutcask(mutcask.CaskNumConf(cfg.CaskNum), mutcask.PathConf(cfg.Path))
	if err != nil {
		return nil, err
	}
	return &blostore{
		hasOnRead: cfg.HashOnRead,
		batch:     cfg.Batch,
		kv:        mc,
	}, nil
}
