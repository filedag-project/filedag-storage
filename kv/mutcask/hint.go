package mutcask

import (
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/xerrors"
)

const vLogSuffix = ".vlog"
const hintLogSuffix = ".hint"

type CaskMap struct {
	sync.Mutex
	m map[uint32]*Cask
}

func (cm *CaskMap) Add(id uint32, cask *Cask) {
	cm.Lock()
	defer cm.Unlock()
	cm.m[id] = cask
}

func (cm *CaskMap) Get(id uint32) (c *Cask, b bool) {
	c, b = cm.m[id]
	return
}

func (cm *CaskMap) CloseAll() {
	for _, cask := range cm.m {
		if cask != nil {
			cask.Close()
		}
	}
}

type KeyMap struct {
	sync.Mutex
	m map[string]*Hint
}

func (km *KeyMap) Add(key string, hint *Hint) {
	km.Lock()
	defer km.Unlock()
	km.m[key] = hint
}

func (km *KeyMap) Get(key string) (h *Hint, b bool) {
	h, b = km.m[key]
	return
}

func buildKeyMap(hint *os.File) (*KeyMap, error) {
	finfo, err := hint.Stat()
	if err != nil {
		return nil, err
	}
	if finfo.Size()%HintEncodeSize != 0 {
		return nil, ErrHintLogBroken
	}
	km := &KeyMap{}
	km.m = make(map[string]*Hint)
	hint.Seek(0, 0)
	offset := uint64(0)
	buf := make([]byte, HintEncodeSize)
	for {
		n, err := hint.Read(buf)
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			// read end file
			break
		}
		// should never happened
		if n != HintEncodeSize {
			return nil, xerrors.Errorf("read hint failed, expected %d bytes, read %d bytes", HintEncodeSize, n)
		}
		h := &Hint{}
		if err = h.From(buf); err != nil {
			return nil, err
		}
		h.KOffset = offset
		offset += HintEncodeSize
		km.Add(h.Key, h)
	}
	return km, nil
}

func buildCaskMap(cfg *Config) (*CaskMap, error) {
	var err error
	dirents, err := os.ReadDir(cfg.Path)
	if err != nil {
		return nil, err
	}
	cm := &CaskMap{}
	cm.m = make(map[uint32]*Cask)
	defer func() {
		if err != nil {
			cm.CloseAll()
		}
	}()

	for _, ent := range dirents {
		if !ent.IsDir() && strings.HasSuffix(ent.Name(), hintLogSuffix) {
			name := strings.TrimSuffix(ent.Name(), hintLogSuffix)
			id, err := strconv.ParseUint(name, 10, 32)
			if err != nil {
				return nil, err
			}
			cask := NewCask()
			cm.Add(uint32(id), cask)
			cask.hintLog, err = os.OpenFile(filepath.Join(cfg.Path, ent.Name()), os.O_RDWR, 0644)
			if err != nil {
				return nil, err
			}
			cask.keyMap, err = buildKeyMap(cask.hintLog)
			if err != nil {
				return nil, err
			}
			cask.vLog, err = os.OpenFile(filepath.Join(cfg.Path, name, vLogSuffix), os.O_RDWR, 0644)
			if err != nil {
				return nil, err
			}
		}
	}

	return cm, nil
}
