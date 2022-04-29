package dag

import (
	"bytes"
	"context"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/pool/config"
	"github.com/filedag-project/filedag-storage/http/objectstore/utils"
	chunker "github.com/ipfs/go-ipfs-chunker"
	"io/ioutil"
	"testing"
)

func TestSimplePool_Add(t *testing.T) {
	dagPool, err := NewSimplePool(&config.SimplePoolConfig{
		StorePath: utils.TmpDirPath(t),
		BatchNum:  4,
		CaskNum:   2,
	})
	if err != nil {
		fmt.Println(err)
	}
	r := ioutil.NopCloser(bytes.NewReader([]byte("hello world")))
	add, err := dagPool.Add(context.Background(), r, "", "")
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
