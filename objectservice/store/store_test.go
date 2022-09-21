package store

import (
	"bytes"
	"context"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/pool/client"
	"github.com/filedag-project/filedag-storage/objectservice/uleveldb"
	"github.com/ipfs/go-blockservice"
	offline "github.com/ipfs/go-ipfs-exchange-offline"
	"github.com/ipfs/go-merkledag"
	"io/ioutil"
	"testing"
)

func TestStorageSys_Object(t *testing.T) {
	poolCli, done := client.NewMockPoolClient(t)
	defer done()
	db, _ := uleveldb.OpenDb(t.TempDir())
	dagServ := merkledag.NewDAGService(blockservice.New(poolCli, offline.Exchange(poolCli)))
	s := NewStorageSys(dagServ, db)
	r := ioutil.NopCloser(bytes.NewReader([]byte("123456")))
	ctx := context.TODO()
	object, err := s.StoreObject(ctx, "testbucket", "testobject", r, 6, map[string]string{})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("object:%v", object)
	getObject, i, err := s.GetObject(ctx, "testbucket", "testobject")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(getObject)
	all, _ := ioutil.ReadAll(i)
	fmt.Println(string(all))
}
