package mutcask

import (
	"context"
	"fmt"
	"hash/crc32"
	"os"
	"path/filepath"
	"sync"

	"github.com/filedag-project/filedag-storage/kv"
)

var _ kv.KVDB = (*mutcask)(nil)

type mutcask struct {
	sync.Mutex
	cfg            *Config
	caskMap        *CaskMap
	createCaskChan chan *createCaskRequst
	close          func()
	closeChan      chan struct{}
}

func NewMutcask(opts ...Option) (*mutcask, error) {
	m := &mutcask{
		cfg:            defaultConfig(),
		createCaskChan: make(chan *createCaskRequst),
		closeChan:      make(chan struct{}),
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
	var once sync.Once
	m.close = func() {
		once.Do(func() {
			close(m.closeChan)
		})
	}
	m.handleCreateCask()
	return m, nil
}

func (m *mutcask) handleCreateCask() {
	go func(m *mutcask) {
		ids := []uint32{}
		for {
			select {
			case <-m.closeChan:
				return
			case req := <-m.createCaskChan:
				func() {
					// fmt.Printf("received cask create request, id = %d\n", req.id)
					if hasId(ids, req.id) {
						req.done <- ErrNone
						return
					}
					cask := NewCask(req.id)
					var err error
					// create vlog file
					cask.vLog, err = os.OpenFile(filepath.Join(m.cfg.Path, m.vLogName(req.id)), os.O_RDWR|os.O_CREATE, 0644)
					if err != nil {
						req.done <- err
						return
					}
					// create hintlog file
					cask.hintLog, err = os.OpenFile(filepath.Join(m.cfg.Path, m.hintLogName(req.id)), os.O_RDWR|os.O_CREATE, 0644)
					if err != nil {
						req.done <- err
						return
					}
					m.caskMap.Add(req.id, cask)
					ids = append(ids, req.id)
					req.done <- ErrNone
				}()
			}
		}
	}(m)
}

func (m *mutcask) vLogName(id uint32) string {
	return fmt.Sprintf("%08d%s", id, vLogSuffix)
}

func (m *mutcask) hintLogName(id uint32) string {
	return fmt.Sprintf("%08d%s", id, hintLogSuffix)
}

func (m *mutcask) Put(key string, value []byte) (err error) {
	id := m.fileID(key)
	var cask *Cask
	var has bool
	cask, has = m.caskMap.Get(id)
	if !has {
		done := make(chan error)
		m.createCaskChan <- &createCaskRequst{
			id:   id,
			done: done,
		}
		if err := <-done; err != ErrNone {
			return err
		}
		cask, _ = m.caskMap.Get(id)
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
		fmt.Println("********")
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
	m.close()
	return nil
}
func (m *mutcask) AllKeysChan(ctx context.Context) (chan string, error) {
	kc := make(chan string)
	go func(ctx context.Context, m *mutcask) {
		defer close(kc)
		for _, cask := range m.caskMap.m {
			for key, h := range cask.keyMap.m {
				if h.Deleted {
					continue
				}
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

type createCaskRequst struct {
	id   uint32
	done chan error
}

func hasId(ids []uint32, id uint32) bool {
	for _, item := range ids {
		if item == id {
			return true
		}
	}
	return false
}
