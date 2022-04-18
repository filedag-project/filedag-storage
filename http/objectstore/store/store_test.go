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

func TestStorageSys_StoreObject(t *testing.T) {
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
		fmt.Println(err)
		return
	}
	fmt.Println(object)
}
