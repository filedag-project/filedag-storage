package store

import (
	"bytes"
	"context"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/pool/client"
	"github.com/filedag-project/filedag-storage/objectservice/objmetadb"
	"github.com/ipfs/go-blockservice"
	offline "github.com/ipfs/go-ipfs-exchange-offline"
	"github.com/ipfs/go-merkledag"
	"io/ioutil"
	"testing"
)

func TestStorageSys_Object(t *testing.T) {
	poolCli, done := client.NewMockPoolClient(t)
	defer done()
	db, _ := objmetadb.OpenDb(t.TempDir())
	dagServ := merkledag.NewDAGService(blockservice.New(poolCli, offline.Exchange(poolCli)))
	s := NewStorageSys(context.TODO(), dagServ, db)
	mbsys := NewBucketMetadataSys(db)
	mbsys.CreateBucket(context.TODO(), "testbucket", "", "")
	s.SetNewBucketNSLock(mbsys.NewNSLock)
	s.SetHasBucket(mbsys.HasBucket)
	r := ioutil.NopCloser(bytes.NewReader([]byte("123456")))
	ctx := context.TODO()
	object, err := s.StoreObject(ctx, "testbucket", "testobject", r, 6, map[string]string{}, false)
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
func TestGetFolder(t *testing.T) {
	testCases := []struct {
		name   string
		o      []ObjectInfo
		prefix string
		loi    *ListObjectsInfo
		expect []string
	}{
		{
			name: "aaa/",
			o: []ObjectInfo{{
				Name: "aaa/ccc/",
			}},
			prefix: "aaa/",
			loi:    &ListObjectsInfo{},
			expect: []string{"ccc/"},
		},
		{
			name: "nil_frefix",
			o: []ObjectInfo{{
				Name: "aaa/ccc/",
			}},
			prefix: "",
			loi:    &ListObjectsInfo{},
			expect: []string{"aaa/"},
		},
		{
			name: "nil_frefix",
			o: []ObjectInfo{
				{
					Name: "aaa/ccc/",
				},
				{
					Name: "aaa/",
				},
			},
			prefix: "",
			loi:    &ListObjectsInfo{},
			expect: []string{"aaa/"},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			m := make(map[string]struct{})
			for _, aaa := range testCase.o {
				getFolder(aaa, testCase.prefix, testCase.loi, m)
			}
			fmt.Println(m)
		})
	}

}
