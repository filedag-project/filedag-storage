package pool

import (
	"bytes"
	"context"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/pool/config"
	"github.com/filedag-project/filedag-storage/http/objectstore/utils"
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
	add, err := dagPool.Add(context.Background(), r)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(add)
}
