package pool

import (
	"bytes"
	"context"
	"fmt"
	"github.com/filedag-project/filedag-storage/http/objectstore/utils/hash"
	"testing"
)

func TestSimplePool_Add(t *testing.T) {
	dagPool, err := NewSimplePool(&SimplePoolConfig{
		StorePath: "./",
		BatchNum:  4,
		CaskNum:   2,
	})
	if err != nil {
		fmt.Println(err)
	}
	var a = []byte("12345")
	hashReader, _ := hash.NewReader(bytes.NewReader(a), int64(len(a)), "", "", int64(len(a)))
	add, err := dagPool.Add(context.Background(), hashReader)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(add)
}
