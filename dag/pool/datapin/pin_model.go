package datapin

import (
	"crypto/sha256"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/pool/datapin/types"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	"sync"
)

type BlockPin struct {
	mu sync.RWMutex
	DB *uleveldb.ULevelDB
}

const dagPoolPin = "dagPoolPin/"

// AddPin add pin
func (p *BlockPin) AddPin(pin types.Pin) error {
	defer p.mu.RUnlock()
	cidCode := sha256String(pin.Cid.String())
	p.mu.RLock()
	err := p.DB.Put(dagPoolPin+cidCode, pin)
	if err != nil {
		return err
	}
	return nil
}

// RemovePin remove pin
func (p *BlockPin) RemovePin(cid string) error {
	cidCode := sha256String(cid)
	err := p.DB.Delete(dagPoolPin + cidCode)
	if err != nil {
		return err
	}
	return nil
}

func (p *BlockPin) QueryPin(cid string) (*types.Pin, error) {
	cidCode := sha256String(cid)
	pin := &types.Pin{}
	err := p.DB.Get(dagPoolPin+cidCode, &pin)
	if err != nil {
		return nil, err
	}
	return pin, nil
}

func NewBlockPin(db *uleveldb.ULevelDB) (BlockPin, error) {
	return BlockPin{DB: db}, nil
}

func sha256String(s string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
}
