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
	poolCli, done := utils.NewMockPoolClient(&testing.T{})
	defer done()
	db, _ := uleveldb.OpenDb(utils.TmpDirPath(&testing.T{}))
	s := NewStorageSys(poolCli, db)

	r := ioutil.NopCloser(bytes.NewReader([]byte("123456")))
	ctx := context.WithValue(context.TODO(), "", "")
	object, err := s.StoreObject(ctx, "test", "testbucket", "testobject", r, 6)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("object:%v", object)
	getObject, i, err := s.GetObject(ctx, "test", "testbucket", "testobject")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(getObject)
	all, _ := ioutil.ReadAll(i)
	fmt.Println(string(all))
}
