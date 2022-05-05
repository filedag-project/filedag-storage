package store

import (
	"bytes"
	"context"
	"fmt"
	pool "github.com/filedag-project/filedag-storage/dag"
	"github.com/filedag-project/filedag-storage/dag/pool/config"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	"github.com/filedag-project/filedag-storage/http/objectstore/utils"
	"io/ioutil"
	"testing"
)

func TestStorageSys_Object(t *testing.T) {
	var s StorageSys
	var err error
	uleveldb.DBClient, err = uleveldb.OpenDb(utils.TmpDirPath(t))
	if err != nil {
		return
	}
	defer uleveldb.DBClient.Close()
	s.Db = uleveldb.DBClient
	s.DagPool, err = pool.NewSimplePool(&config.SimplePoolConfig{
		StorePath: utils.TmpDirPath(t),
		BatchNum:  4,
		CaskNum:   2,
	})
	if err != nil {
		return
	}
	file, err := ioutil.ReadFile("./store_test.go")
	if err != nil {
		return
	}
	file = bytes.Repeat(file, 10000)
	r := ioutil.NopCloser(bytes.NewReader(file))
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
