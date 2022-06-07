package poolservice

import (
	"context"
	"github.com/ipfs/go-cid"
)

func (d *dagPoolService) Pin(ctx context.Context, c cid.Cid) error {
	links, err := d.GetLinks(ctx, c)
	if err != nil {
		return err
	}
	err = d.refer.AddReference(c.String())
	if err != nil {
		return err
	}
	for _, link := range links {
		err = d.refer.AddReference(link.Cid.String())
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *dagPoolService) UnPin(ctx context.Context, c cid.Cid) error {
	links, err := d.GetLinks(ctx, c)
	if err != nil {
		return err
	}
	err = d.refer.RemoveReference(c.String())
	if err != nil {
		return err
	}
	for _, link := range links {
		err = d.refer.RemoveReference(link.Cid.String())
		if err != nil {
			return err
		}
	}
	return nil
}
