package store

import (
	"bytes"
	"context"
	"fmt"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	"github.com/filedag-project/filedag-storage/http/objectstore/utils"
	"github.com/filedag-project/filedag-storage/http/objectstore/utils/dagpoolclient"
	"io/ioutil"
	"os"
	"testing"
)

func TestStorageSys_Object(t *testing.T) {
	//server.StartTestServer(t)
	var s StorageSys
	s.DagPool, _ = dagpoolclient.NewPoolClient("localhost:9002")
	s.Db, _ = uleveldb.OpenDb(utils.TmpDirPath(&testing.T{}))
	os.Setenv(PoolUser, "pool")
	os.Setenv(PoolPass, "pool")
	r := ioutil.NopCloser(bytes.NewReader([]byte("")))
	object, err := s.StoreObject(context.TODO(), "test", "testbucket", "testobject", r)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("object:%v", object)
	//getObject, i, err := s.GetObject(context.TODO(), "test", "testbucket", "testobject")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(getObject)
	//all,_:=ioutil.ReadAll(i)
	//fmt.Println(string(all))
	s.DagPool.Close(context.TODO())
}
