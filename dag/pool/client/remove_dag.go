package client

import (
	"context"
	"github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag"
)

func RemoveDAG(ctx context.Context, dagServ ipld.DAGService, root cid.Cid) error {
	list := make([]cid.Cid, 0, 32)
	visit := func(c cid.Cid) bool {
		list = append(list, c)
		return true
	}
	err := merkledag.Walk(
		ctx, merkledag.GetLinksWithDAG(dagServ), root,
		visit,
		merkledag.Concurrent(),
	)
	if err != nil {
		return err
	}
	for _, c := range list {
		if err := dagServ.Remove(ctx, c); err != nil {
			log.Errorf("remove block failed, error: %v", err)
		}
	}
	return nil
}
