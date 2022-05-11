package dag

import (
	"bytes"
	"context"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/node"
	"github.com/filedag-project/filedag-storage/dag/pool"
	"github.com/filedag-project/filedag-storage/dag/pool/userpolicy"
	"github.com/filedag-project/filedag-storage/http/objectstore/utils"
	chunker "github.com/ipfs/go-ipfs-chunker"
	"io/ioutil"
	"os"
	"testing"
)

func TestSimplePool_Add(t *testing.T) {
	dagPool, ctx := testInit(t)
	r := ioutil.NopCloser(bytes.NewReader([]byte("hello world")))
	add, err := dagPool.Add(ctx, r)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(add)
}
func TestSimplePool_Get(t *testing.T) {
	dagPool, ctx := testInit(t)
	f, _ := ioutil.ReadFile("C:\\Users\\dean\\Downloads\\SunloginClient_12.5.1.45098_x64.exe")
	r := ioutil.NopCloser(bytes.NewReader(f))
	add, err := dagPool.Add(ctx, r)
	if err != nil {
		fmt.Println(err)
	}
	g, err := dagPool.Get(ctx, add)
	if err != nil {
		fmt.Println(err)
	}
	all, err := ioutil.ReadAll(g)
	if err != nil {
		return
	}
	fmt.Println(string(all))
}
func testInit(t *testing.T) (dagPool *simplePool, ctx context.Context) {
	os.Setenv(pool.DagNodeIpOrPath, utils.TmpDirPath(t))
	os.Setenv(pool.DagPoolImporterBatchNum, "4")
	os.Setenv(pool.DagPoolLeveldbPath, utils.TmpDirPath(t))
	os.Setenv(node.NodeConfigPath, "config.json")
	dagPool, err := NewSimplePool()
	if err != nil {
		fmt.Println(err)
	}
	ctx = context.Background()
	ctx = context.WithValue(ctx, "user", "test,test123")
	dagPool.AddUser(context.TODO(), "test", "test123", userpolicy.ReadWrite, 100000)
	return dagPool, ctx
}
func TestSimplePool_a(t *testing.T) {
	f := bytes.NewReader([]byte(""))
	a := chunker.NewSizeSplitter(f, int64(2))
	fmt.Println(a.NextBytes())
}
