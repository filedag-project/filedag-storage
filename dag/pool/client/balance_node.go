package client

import (
	"github.com/ipfs/go-cid"
	chunker "github.com/ipfs/go-ipfs-chunker"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-unixfs/importer/balanced"
	h "github.com/ipfs/go-unixfs/importer/helpers"
	"io"
)

const UnixfsLinksPerLevel = 1 << 10
const UnixfsChunkSize uint64 = 1 << 20

//BalanceNode split the file and store it in DAGService as node
func BalanceNode(f io.Reader, bufDs ipld.DAGService, cidBuilder cid.Builder) (node ipld.Node, err error) {
	params := h.DagBuilderParams{
		Maxlinks:   UnixfsLinksPerLevel,
		RawLeaves:  false,
		CidBuilder: cidBuilder,
		Dagserv:    bufDs,
		NoCopy:     false,
	}
	db, err := params.New(chunker.NewSizeSplitter(f, int64(UnixfsChunkSize)))
	if err != nil {
		return nil, err
	}
	node, err = balanced.Layout(db)
	if err != nil {
		return nil, err
	}
	return
}
