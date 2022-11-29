package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/ipfs/go-cid"
	chunker "github.com/ipfs/go-ipfs-chunker"
	ipld "github.com/ipfs/go-ipld-format"
	dag "github.com/ipfs/go-merkledag"
	ft "github.com/ipfs/go-unixfs"
	"github.com/ipfs/go-unixfs/importer/balanced"
	h "github.com/ipfs/go-unixfs/importer/helpers"
	"io"
)

const unixfsLinksPerLevel = 1 << 10
const unixfsChunkSize uint64 = 1 << 20

//BalanceNode split the file and store it in DAGService as node
func BalanceNode(f io.Reader, bufDs ipld.DAGService, cidBuilder cid.Builder) (node ipld.Node, err error) {
	params := h.DagBuilderParams{
		Maxlinks:   unixfsLinksPerLevel,
		RawLeaves:  false,
		CidBuilder: cidBuilder,
		Dagserv:    bufDs,
		NoCopy:     false,
	}
	db, err := params.New(chunker.NewSizeSplitter(f, int64(unixfsChunkSize)))
	if err != nil {
		return nil, err
	}
	node, err = balanced.Layout(db)
	if err != nil {
		return nil, err
	}
	return
}

type LinkInfo struct {
	Link     *ipld.Link
	FileSize uint64
}

type unixfsNode struct {
	dag  *dag.ProtoNode
	file *ft.FSNode
}

func NewUnixfsNodeFromDag(nd *dag.ProtoNode) (*unixfsNode, error) {
	mb, err := ft.FSNodeFromBytes(nd.Data())
	if err != nil {
		return nil, err
	}

	return &unixfsNode{
		dag:  nd,
		file: mb,
	}, nil
}

func (n *unixfsNode) AddChild(child *ipld.Link, fileSize uint64) error {
	err := n.dag.AddRawLink("", child)
	if err != nil {
		return err
	}

	n.file.AddBlockSize(fileSize)

	return nil
}

func (n *unixfsNode) Commit() (ipld.Node, error) {
	fileData, err := n.file.GetBytes()
	if err != nil {
		return nil, err
	}
	n.dag.SetData(fileData)

	return n.dag, nil
}

func BuildDataCidByLinks(ctx context.Context, dagServ ipld.DAGService, cidBuilder cid.Builder, links []LinkInfo) (cid.Cid, error) {
	var linkList = make([]LinkInfo, 0)
	var needAdd = make([]ipld.Node, 0)

	for len(links) > 1 {
		nd := ft.EmptyFileNode()
		nd.SetCidBuilder(cidBuilder)
		od, err := NewUnixfsNodeFromDag(nd)
		if err != nil {
			return cid.Undef, err
		}
		count := 0
		for _, link := range links {
			if count >= unixfsLinksPerLevel {
				nnd, err := od.Commit()
				if err != nil {
					return cid.Undef, err
				}
				needAdd = append(needAdd, nnd)
				lk, err := ipld.MakeLink(nnd)
				if err != nil {
					return cid.Undef, err
				}
				linkList = append(linkList, LinkInfo{
					Link:     lk,
					FileSize: od.file.FileSize(),
				})

				nd = ft.EmptyFileNode()
				nd.SetCidBuilder(cidBuilder)
				od, err = NewUnixfsNodeFromDag(nd)
				if err != nil {
					return cid.Undef, err
				}
				count = 0
			}
			if err := od.AddChild(link.Link, link.FileSize); err != nil {
				return cid.Undef, err
			}
			count++
		}
		if len(nd.Links()) > 0 {
			nnd, err := od.Commit()
			if err != nil {
				return cid.Undef, err
			}
			needAdd = append(needAdd, nnd)
			lk, err := ipld.MakeLink(nnd)
			if err != nil {
				return cid.Undef, err
			}
			linkList = append(linkList, LinkInfo{
				Link:     lk,
				FileSize: od.file.FileSize(),
			})
		}
		links = linkList
		linkList = linkList[:0]
	}

	if len(needAdd) > 0 {
		if err := dagServ.AddMany(ctx, needAdd); err != nil {
			return cid.Undef, err
		}
	}
	return links[0].Link.Cid, nil
}

func CreateLinkInfo(ctx context.Context, dagServ ipld.DAGService, c cid.Cid) (LinkInfo, error) {
	nd, err := dagServ.Get(ctx, c)
	if err != nil {
		return LinkInfo{}, err
	}
	pn, ok := nd.(*dag.ProtoNode)
	if !ok {
		return LinkInfo{}, errors.New(fmt.Sprintf("node %s is not ProtoNode", c.String()))
	}
	ftn, err := ft.FSNodeFromBytes(pn.Data())
	if err != nil {
		return LinkInfo{}, err
	}
	lk, err := ipld.MakeLink(nd)
	if err != nil {
		return LinkInfo{}, err
	}
	return LinkInfo{
		Link:     lk,
		FileSize: ftn.FileSize(),
	}, nil
}
