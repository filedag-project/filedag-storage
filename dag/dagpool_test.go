package dag

import (
	"bytes"
	"context"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/node"
	"github.com/filedag-project/filedag-storage/dag/pool/config"
	"github.com/filedag-project/filedag-storage/dag/pool/userpolicy"
	"github.com/filedag-project/filedag-storage/http/objectstore/utils"
	chunker "github.com/ipfs/go-ipfs-chunker"
	"io/ioutil"
	"testing"
)

func TestSimplePool_Add(t *testing.T) {
	nodec := []node.Config{
		{
			Batch:   4,
			Path:    utils.TmpDirPath(t),
			CaskNum: 2,
		},
	}
	dagPool, err := NewSimplePool(&config.SimplePoolConfig{
		NodesConfig: nodec,
		StorePath:   utils.TmpDirPath(t),
		BatchNum:    4,
		CaskNum:     2,
		LeveldbPath: utils.TmpDirPath(t),
	})
	if err != nil {
		fmt.Println(err)
	}
	r := ioutil.NopCloser(bytes.NewReader([]byte("hello world")))
	ctx := context.Background()
	ctx = context.WithValue(ctx, "user", "test,test123")
	dagPool.AddUser(context.TODO(), "test", "test123", userpolicy.ReadWrite, 100000)
	add, err := dagPool.Add(ctx, r)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(add)
}
func TestSimplePool_a(t *testing.T) {
	f := bytes.NewReader([]byte(""))
	a := chunker.NewSizeSplitter(f, int64(2))
	fmt.Println(a.NextBytes())
}
