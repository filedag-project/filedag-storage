package diskv

import (
	"fmt"
	"sync"
	"testing"
)

func TestDisKVPut(t *testing.T) {
	var kvdata = []kvt{
		{[]byte("city"), []byte("shanghai")},
		{[]byte("app"), []byte("filedag")},
		{[]byte("protocol"), []byte("ipfs")},
		{[]byte("blockchain"), []byte("filecoin")},
		{[]byte("fs"), []byte("nfs")},
		{[]byte("kv"), []byte("leveldb")},
		{[]byte("db"), []byte("mongodb")},
		{[]byte("language"), []byte("golang")},
		{[]byte("b"), []byte("1")},
		{[]byte("c"), []byte("2")},
		{[]byte("d"), []byte("3")},
		{[]byte("e"), []byte("4")},
		{[]byte("f"), []byte("5")},
		{[]byte("g"), []byte("6")},
		{[]byte("h"), []byte("7")},
		{[]byte("i"), []byte("8")},
		{[]byte("j"), []byte("9")},
	}
	dkv, err := NewDisKV(func(cfg *Config) {
		cfg.Dir = "/Users/lifeng/testdir/diskv"
		cfg.Parallel = 64
	})
	if err != nil {
		t.Fatal(err)
	}
	defer dkv.Close()
	var wg sync.WaitGroup
	wg.Add(len(kvdata))
	for _, item := range kvdata {
		go func(k, v []byte) {
			defer wg.Done()

			err = dkv.Put(k, v)
			if err != nil {
				fmt.Println(err)
				t.Fail()
			}
		}(item.Key, item.Value)
	}
	wg.Wait()
	//t.Fail()
}
