package store

import (
	"bytes"
	"context"
	"fmt"
	dagpoolcli "github.com/filedag-project/filedag-storage/dag/pool/client"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	"github.com/filedag-project/filedag-storage/http/objectstore/utils"
	"github.com/golang/mock/gomock"
	"github.com/ipfs/go-merkledag"
	"io/ioutil"
	"testing"
)

func TestStorageSys_Object(t *testing.T) {
	var s StorageSys
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s.DagPool = NewMockClient(ctrl)
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
func NewMockClient(ctrl *gomock.Controller) dagpoolcli.PoolClient {
	m := dagpoolcli.NewMockPoolClient(ctrl)
	m.EXPECT().Add(gomock.Any(), gomock.Any()).Return(nil)
	m.EXPECT().Get(gomock.Any(), gomock.Any()).Return(merkledag.NewRawNode([]byte("123456")), nil)
	m.EXPECT().Close(gomock.Any())
	return m
}
