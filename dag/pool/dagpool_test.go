package pool

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestSimplePool_Add(t *testing.T) {
	dagPool, err := NewSimplePool(&SimplePoolConfig{
		StorePath: "./test",
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
