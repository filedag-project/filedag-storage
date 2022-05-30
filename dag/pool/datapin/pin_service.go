package datapin

import (
	"github.com/filedag-project/filedag-storage/dag/pool"
	"github.com/filedag-project/filedag-storage/dag/pool/datapin/types"
	"github.com/ipfs/go-cid"
	"sync"
)

type pinner struct {
	lock        sync.RWMutex
	recursePin  *cid.Set
	directPin   *cid.Set
	internalPin *cid.Set
	DagPool     *pool.DagPool
}

type PinService struct {
	blockPin BlockPin
}

func (s *PinService) addPin(pin Pin) error {
	poolPin, err := svcPinToPoolPin(pin)
	if err != nil {
		return err
	}
	err = s.blockPin.AddPin(poolPin)
	if err != nil {
		return err
	}
	return nil
}

func (s *PinService) removePin(cid types.Cid) error {
	err := s.blockPin.RemovePin(cid.String())
	if err != nil {
		return err
	}
	return nil

}

func svcPinToPoolPin(p Pin) (types.Pin, error) {
	opts := types.PinOptions{
		Name:     string(p.Name),
		Origins:  p.Origins,
		Metadata: p.Meta,
		Mode:     types.PinModeDirect,
	}
	return types.PinWithOpts(p.Cid, opts), nil
}
