package datapin

import (
	"crypto/sha256"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/pool/datapin/types"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	"golang.org/x/xerrors"
	"sync"
)

type BlockPin struct {
	mu sync.Mutex
	DB *uleveldb.ULevelDB
}

const dagPoolPin = "dagPoolPin/"

// AddPin add pin
func (p *BlockPin) AddPin(pin types.Pin) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.Has(pin.Cid.String()) {
		return nil
	}
	cidCode := sha256String(pin.Cid.String())
	err := p.DB.Put(dagPoolPin+cidCode, pin)
	if err != nil {
		return err
	}
	return nil
}

// RemovePin remove pin
func (p *BlockPin) RemovePin(cid string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.Has(cid) {
		return xerrors.Errorf("not found : %v", cid)
	}
	cidCode := sha256String(cid)
	err := p.DB.Delete(dagPoolPin + cidCode)
	if err != nil {
		return err
	}
	return nil
}

func (p *BlockPin) QueryPin(cid string) (*types.Pin, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	cidCode := sha256String(cid)
	pin := &types.Pin{}
	err := p.DB.Get(dagPoolPin+cidCode, &pin)
	if err != nil {
		return nil, err
	}
	return pin, nil
}
func (p *BlockPin) Has(cid string) bool {
	cidCode := sha256String(cid)
	pin := &types.Pin{}
	err := p.DB.Get(dagPoolPin+cidCode, &pin)
	if err != nil {
		return false
	}
	return true
}

func NewBlockPin(db *uleveldb.ULevelDB) (BlockPin, error) {
	return BlockPin{DB: db}, nil
}

func sha256String(s string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
}
