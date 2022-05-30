package datapin

import (
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
	err := p.DB.Put(dagPoolPin+pin.Cid.String(), pin)
	if err != nil {
		return err
	}
	return nil
}

// RemovePin remove pin
func (p *BlockPin) RemovePin(cid string) error {
	err := p.DB.Delete(dagPoolPin + cid)
	if err != nil {
		return err
	}
	return nil
}
