package mutcask

import (
	"context"
	"fmt"
	"hash/crc32"
	"os"
	"path/filepath"

	"github.com/filedag-project/filedag-storage/kv"
)

var _ kv.KVDB = (*mutcask)(nil)

type mutcask struct {
	cfg     *Config
	caskMap *CaskMap
}

func NewMutcask(opts ...Option) (*mutcask, error) {
	m := &mutcask{
		cfg: defaultConfig(),
	}
	for _, opt := range opts {
		opt(m.cfg)
	}
	repoPath := m.cfg.Path
	if repoPath == "" {
		return nil, ErrPathUndefined
	}
	repo, err := os.Stat(repoPath)
	if err == nil && !repo.IsDir() {
		return nil, ErrPath
	}
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		if err := os.Mkdir(repoPath, 0755); err != nil {
			return nil, err
		}
	}
	m.caskMap, err = buildCaskMap(m.cfg)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (m *mutcask) vLogName(id uint32) string {
	return fmt.Sprintf("%08d%s", id, vLogSuffix)
}

func (m *mutcask) hintLogName(id uint32) string {
	return fmt.Sprintf("%08d%s", id, hintLogSuffix)
}

func (m *mutcask) Put(key string, value []byte) (err error) {
	id := m.fileID(key)
	cask, has := m.caskMap.Get(id)
	if !has {
		cask = NewCask()
		// create vlog file
		cask.vLog, err = os.OpenFile(filepath.Join(m.cfg.Path, m.vLogName(id)), os.O_RDWR, 0644)
		if err != nil {
			return
		}
		cask.hintLog, err = os.OpenFile(filepath.Join(m.cfg.Path, m.hintLogName(id)), os.O_RDWR, 0644)
		if err != nil {
			return
		}
		return
	}
	return cask.Put(key, value)
}

func (m *mutcask) Delete(key string) error {
	id := m.fileID(key)
	cask, has := m.caskMap.Get(id)
	if !has {
		return nil
	}
	return cask.Delete(key)
}

func (m *mutcask) Get(key string) ([]byte, error) {
	id := m.fileID(key)
	cask, has := m.caskMap.Get(id)
	if !has {
		return nil, kv.ErrNotFound
	}

	return cask.Read(key)
}

func (m *mutcask) Size(key string) (int, error) {
	id := m.fileID(key)
	cask, has := m.caskMap.Get(id)
	if !has {
		return -1, kv.ErrNotFound
	}
	return cask.Size(key)
}

func (m *mutcask) Close() error {
	m.caskMap.CloseAll()
	return nil
}
func (m *mutcask) AllKeysChan(ctx context.Context) (chan string, error) {
	kc := make(chan string)
	go func(ctx context.Context, m *mutcask) {
		for _, cask := range m.caskMap.m {
			for key := range cask.keyMap.m {
				select {
				case <-ctx.Done():
					return
				default:
					kc <- key
				}
			}
		}
	}(ctx, m)
	return kc, nil
}

func (m *mutcask) fileID(key string) uint32 {
	crc := crc32.ChecksumIEEE([]byte(key))
	return crc % m.cfg.CaskNum
}
