package store

import (
	"bytes"
	"context"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/pool/client"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	"github.com/filedag-project/filedag-storage/http/objectstore/utils"
	"github.com/golang/mock/gomock"
	"github.com/ipfs/go-merkledag"
	"google.golang.org/grpc"
	"io/ioutil"
	"testing"
)

func TestStorageSys_Object(t *testing.T) {
	var s StorageSys
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	co, err := grpc.Dial("127.0.0.1:7777", grpc.WithInsecure())
	cidBuilder, err := merkledag.PrefixForCidVersion(0)
	s.CidBuilder = cidBuilder
	s.DagPool = &client.DagPoolClient{DPClient: utils.NewMockDagPoolClient(ctrl), Conn: co}
	s.Db, _ = uleveldb.OpenDb(utils.TmpDirPath(&testing.T{}))
	r := ioutil.NopCloser(bytes.NewReader([]byte("123456")))
	object, err := s.StoreObject(context.TODO(), "test", "testbucket", "testobject", r, 6)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("object:%v", object)
	getObject, i, err := s.GetObject(context.TODO(), "test", "testbucket", "testobject")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(getObject)
	all, _ := ioutil.ReadAll(i)
	fmt.Println(string(all))
	s.DagPool.Close(context.TODO())
}
