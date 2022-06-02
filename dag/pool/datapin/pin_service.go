package datapin

import (
	"context"
	"github.com/filedag-project/filedag-storage/dag/pool"
	"github.com/filedag-project/filedag-storage/dag/pool/datapin/types"
	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	legacy "github.com/ipfs/go-ipld-legacy"
	logging "github.com/ipfs/go-log/v2"
	"sync"
)

type Pinner struct {
	lock        sync.RWMutex
	recursePin  *cid.Set
	directPin   *cid.Set
	internalPin *cid.Set
	DagPool     *pool.DagPool
}

var log = logging.Logger("data-pin")

type PinService struct {
	blockPin BlockPin
	pinner   Pinner
}

func (s *PinService) AddPin(ctx context.Context, cid cid.Cid, block blocks.Block) error {
	node, err := legacy.DecodeNode(ctx, block)
	if err != nil {
		return err
	}
	var rootPin = Pin{Cid: types.Cid{Cid: cid}}
	poolRootPin, err := svcPinToPoolPin(rootPin, types.PinModeDirect)
	s.pinner.directPin.Add(cid)
	if err != nil {
		return err
	}
	err = s.blockPin.AddPin(poolRootPin)
	if err != nil {
		return err
	}
	for _, link := range node.Links() {
		var pin = Pin{Cid: types.Cid{Cid: link.Cid}}
		poolPin, err := svcPinToPoolPin(pin, types.PinModeRecursive)
		if err != nil {
			return err
		}
		s.pinner.recursePin.Add(cid)
		err = s.blockPin.AddPin(poolPin)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *PinService) RemovePin(ctx context.Context, cid cid.Cid, block blocks.Block) error {
	node, err := legacy.DecodeNode(ctx, block)
	if err != nil {
		return err
	}
	err = s.blockPin.RemovePin(cid.String())
	if err != nil {
		return err
	}
	for _, link := range node.Links() {
		err := s.blockPin.RemovePin(link.Cid.String())
		if err != nil {
			return err
		}
	}
	return nil

}

func svcPinToPoolPin(p Pin, mod types.PinMode) (types.Pin, error) {
	opts := types.PinOptions{
		Name:     string(p.Name),
		Metadata: p.Meta,
		Mode:     mod,
	}
	return types.PinWithOpts(p.Cid, opts), nil
}
