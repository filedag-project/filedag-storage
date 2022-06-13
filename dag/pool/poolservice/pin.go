package poolservice

import (
	"context"
	"github.com/filedag-project/filedag-storage/dag/pool/userpolicy"
	"github.com/ipfs/go-cid"
)

func (d *dagPoolService) Pin(ctx context.Context, c cid.Cid) error {
	if !d.iam.CheckUserPolicy("", "", userpolicy.OnlyRead) {
		return userpolicy.AccessDenied
	}
	links, err := d.GetLinks(ctx, c)
	if err != nil {
		return err
	}
	err = d.refer.AddReference(c.String(), true)
	if err != nil {
		return err
	}
	for _, link := range links {
		err = d.refer.AddReference(link.Cid.String(), true)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *dagPoolService) UnPin(ctx context.Context, c cid.Cid) error {
	if !d.iam.CheckUserPolicy("", "", userpolicy.OnlyRead) {
		return userpolicy.AccessDenied
	}
	links, err := d.GetLinks(ctx, c)
	if err != nil {
		return err
	}
	err = d.refer.RemoveReference(c.String(), true)
	if err != nil {
		return err
	}
	for _, link := range links {
		err = d.refer.RemoveReference(link.Cid.String(), true)
		if err != nil {
			return err
		}
	}
	return nil
}
