package blockpinner

import (
	"bytes"
	"context"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/pool/client"
	leveldb "github.com/filedag-project/filedag-storage/dag/pool/leveldb_datastore"
	"github.com/filedag-project/filedag-storage/http/objectstore/utils"
	blocks "github.com/ipfs/go-block-format"
	bs "github.com/ipfs/go-blockservice"
	ds "github.com/ipfs/go-datastore"
	dssync "github.com/ipfs/go-datastore/sync"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	offline "github.com/ipfs/go-ipfs-exchange-offline"
	legacy "github.com/ipfs/go-ipld-legacy"
	"github.com/ipfs/go-merkledag"
	"io/ioutil"
	"testing"
)

func TestPin_test(t *testing.T) {
	f, _ := ioutil.ReadFile("pin_service.go")
	dstore := dssync.MutexWrap(ds.NewMapDatastore())
	bstore := blockstore.NewBlockstore(dstore)
	bserv := bs.New(bstore, offline.Exchange(bstore))

	dags := merkledag.NewDAGService(bserv)
	cidBuilder, _ := merkledag.PrefixForCidVersion(0)
	ctx := context.Background()
	dbh, _ := client.BalanceNode(bytes.NewReader(bytes.Repeat(f, 2000)), dags, cidBuilder)
	get, _ := dags.Get(ctx, dbh.Cid())
	node, _ := legacy.DecodeNode(ctx, blocks.NewBlock(get.RawData()))
	fmt.Println(node.Cid())
	datastore, err := leveldb.NewDatastore(utils.TmpDirPath(&testing.T{}), nil)
	if err != nil {
		return
	}
	pinner, err := New(ctx, datastore)
	if err != nil {
		return
	}
	addPin, err := pinner.AddPin(ctx, node.Cid(), Recursive, "")
	if err != nil {
		return
	}
	fmt.Println(addPin)
	ok, err := pinner.RemovePinsForCid(ctx, node.Cid(), Recursive)
	if err != nil {
		return
	}
	fmt.Println(ok)
}
