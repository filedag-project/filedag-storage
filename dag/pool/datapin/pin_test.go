package datapin

import (
	"bytes"
	"context"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/pool/client"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	"github.com/filedag-project/filedag-storage/http/objectstore/utils"
	blocks "github.com/ipfs/go-block-format"
	legacy "github.com/ipfs/go-ipld-legacy"
	"github.com/ipfs/go-merkledag"
	"io/ioutil"
	"testing"
)

func TestPin_test(t *testing.T) {
	f, _ := ioutil.ReadFile("pin_service.go")
	dags, _ := client.NewPoolClient(":50001")
	cidBuilder, _ := merkledag.PrefixForCidVersion(0)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "user", "pool,pool123")
	dbh, _ := client.BalanceNode(ctx, bytes.NewReader(bytes.Repeat(f, 2000)), dags, cidBuilder)
	get, _ := dags.Get(ctx, dbh.Cid())
	node, _ := legacy.DecodeNode(ctx, blocks.NewBlock(get.RawData()))
	fmt.Println(node.Cid())

	db, _ := uleveldb.OpenDb(utils.TmpDirPath(&testing.T{}))
	blockPin, err := NewBlockPin(db)
	pinSer := PinService{
		blockPin: blockPin,
	}
	err = pinSer.AddPin(context.TODO(), node.Cid(), blocks.NewBlock(get.RawData()))
	if err != nil {
		fmt.Println(err)
	}
	err = pinSer.RemovePin(context.TODO(), node.Cid(), blocks.NewBlock(get.RawData()))
	if err != nil {
		fmt.Println(err)
	}
}
