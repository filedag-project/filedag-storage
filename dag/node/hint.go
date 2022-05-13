package node

import (
	"io"
	"os"
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
	// cm.Lock()
	// defer cm.Unlock()
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
	// km.Lock()
	// defer km.Unlock()
	km.m[key] = hint
}

func (km *KeyMap) Get(key string) (h *Hint, b bool) {
	// km.Lock()
	// defer km.Unlock()
	h, b = km.m[key]
	return
}

// func (km *KeyMap) Remove(key string) {
// 	km.Lock()
// 	defer km.Unlock()
// 	delete(km.m, key)
// }

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

func fileSize(f *os.File) (uint64, error) {
	finfo, err := f.Stat()
	if err != nil {
		return 0, err
	}
	return uint64(finfo.Size()), nil
}
