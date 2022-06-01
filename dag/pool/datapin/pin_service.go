package datapin

import (
	"context"
	"github.com/filedag-project/filedag-storage/dag/pool/datapin/types"
	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	legacy "github.com/ipfs/go-ipld-legacy"
	logging "github.com/ipfs/go-log/v2"
)

//type pinner struct {
//	lock        sync.RWMutex
//	recursePin  *cid.Set
//	directPin   *cid.Set
//	internalPin *cid.Set
//	DagPool     *pool.DagPool
//}
var log = logging.Logger("data-pin")

type PinService struct {
	blockPin BlockPin
}

func (s *PinService) AddPin(ctx context.Context, cid cid.Cid, block blocks.Block) error {
	node, err := legacy.DecodeNode(ctx, block)
	if err != nil {
		return err
	}
	var rootPin = Pin{Cid: types.Cid{Cid: cid}}
	poolRootPin, err := svcPinToPoolPin(rootPin)
	if err != nil {
		return err
	}
	err = s.blockPin.AddPin(poolRootPin)
	if err != nil {
		return err
	}
	for _, link := range node.Links() {
		var pin = Pin{Cid: types.Cid{Cid: link.Cid}}
		poolPin, err := svcPinToPoolPin(pin)
		if err != nil {
			return err
		}
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

func svcPinToPoolPin(p Pin) (types.Pin, error) {
	opts := types.PinOptions{
		Name:     string(p.Name),
		Metadata: p.Meta,
		Mode:     types.PinModeDirect,
	}
	return types.PinWithOpts(p.Cid, opts), nil
}
