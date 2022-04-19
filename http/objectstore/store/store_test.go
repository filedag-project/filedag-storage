package store

import (
	"bytes"
	"context"
	"fmt"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	"github.com/filedag-project/filedag-storage/http/objectstore/utils"
	"io/ioutil"
	"testing"
)

func TestStorageSys_Object(t *testing.T) {
	var s StorageSys
	var err error
	uleveldb.DBClient, err = uleveldb.OpenDb(utils.TmpDirPath(&testing.T{}))
	if err != nil {
		return
	}
	defer uleveldb.DBClient.Close()
	err = s.Init()
	if err != nil {
		return
	}
	r := ioutil.NopCloser(bytes.NewReader([]byte("hello world")))
	object, err := s.StoreObject(context.Background(), "test", "testBucket", "testobject", r)
	if err != nil {
		fmt.Println("StoreObject", err)
		return
	}
	fmt.Println(object)
	res, i, err := s.GetObject(context.Background(), "test", "testBucket", "testobject")
	if err != nil {
		fmt.Println("GetObject", err)
		return
	}
	all, err := ioutil.ReadAll(i)
	if err != nil {
		return
	}
	fmt.Printf("res:%v,\ni:%v", res, string(all))
}
