package mutcask

import (
	"context"

	"github.com/filedag-project/filedag-storage/kv"
)

var _ kv.KVDB = (*mutcask)(nil)

type mutcask struct {
	cfg *Config
}

func NewMutcask(opts ...Option) *mutcask {
	m := &mutcask{
		cfg: defaultConfig(),
	}
	for _, opt := range opts {
		opt(m.cfg)
	}
	return m
}

func (m *mutcask) Put(key string, value []byte) error {
	return nil
}

func (m *mutcask) Delete(key string) error {
	return nil
}

func (m *mutcask) Get(key string) ([]byte, error) {
	return nil, nil
}

func (m *mutcask) Size(key string) (int, error) {
	return -1, nil
}

func (m *mutcask) Close() error {
	return nil
}
func (m *mutcask) AllKeysChan(ctx context.Context) (chan string, error) {
	return nil, nil
}
