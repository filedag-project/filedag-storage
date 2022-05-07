package store

import (
	"bytes"
	"context"
	"fmt"
	pool "github.com/filedag-project/filedag-storage/dag"
	"github.com/filedag-project/filedag-storage/dag/node"
	"github.com/filedag-project/filedag-storage/dag/pool/config"
	"github.com/filedag-project/filedag-storage/dag/pool/userpolicy"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	"github.com/filedag-project/filedag-storage/http/objectstore/utils"
	"io/ioutil"
	"os"
	"testing"
)

func TestStorageSys_Object(t *testing.T) {
	err := os.Setenv(PoolUser, "test")
	if err != nil {
		return
	}
	err = os.Setenv(PoolPass, "test123")
	if err != nil {
		return
	}
	os.Setenv(PoolDbpath, utils.TmpDirPath(t))
	var s StorageSys
	uleveldb.DBClient, err = uleveldb.OpenDb(utils.TmpDirPath(t))
	if err != nil {
		return
	}
	defer uleveldb.DBClient.Close()
	s.Db = uleveldb.DBClient
	nodec := []node.Config{
		{
			Batch:   4,
			Path:    "testconfig.json",
			CaskNum: 2,
		},
	}
	s.DagPool, err = pool.NewSimplePool(&config.SimplePoolConfig{
		NodesConfig:      nodec,
		LeveldbPath:      os.Getenv(PoolDbpath),
		ImporterBatchNum: defaultPoolBatchNum,
	})
	if err != nil {
		return
	}
	s.DagPool.AddUser(context.TODO(), "test", "test123", userpolicy.ReadWrite, 100000)
	file, err := ioutil.ReadFile("./store_test.go")
	if err != nil {
		return
	}
	file = bytes.Repeat(file, 10000)
	r := ioutil.NopCloser(bytes.NewReader(file))
	ctx := context.Background()
	object, err := s.StoreObject(ctx, "test", "testBucket", "testobject", r)
	if err != nil {
		fmt.Println("StoreObject", err)
		return
	}
	fmt.Println(object)
	res, i, err := s.GetObject(ctx, "test", "testBucket", "testobject")
	if err != nil {
		fmt.Println("GetObject", err)
		return
	}
	all, err := ioutil.ReadAll(i)
	if err != nil {
		return
	}
	fmt.Println(len(file))
	fmt.Printf("res:%v,\ni:%v", res, len(all))
}
