package diskv

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
	"testing"
)

type kvt struct {
	Key   string
	Value []byte
}

func TestRefdb(t *testing.T) {
	var kvdata = []kvt{
		{"city", []byte("shanghai")},
		{"app", []byte("filedag")},
		{"protocol", []byte("ipfs")},
		{"blockchain", []byte("filecoin")},
	}
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("failed to make temp dir: %s", err)
	}
	rfdb, err := NewRefdb(dir)
	if err != nil {
		t.Fatal(err)
	}
	// put
	for _, item := range kvdata {
		if err := rfdb.Put(item.Key, item.Value); err != nil {
			t.Fatal(err)
		}
	}
	// get
	for _, item := range kvdata {
		bs, err := rfdb.Get(item.Key)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(bs, item.Value) {
			t.Fatalf("mismatched data: %v, expected: %v", bs, item.Value)
		}
	}
	// test all keys chan
	kc, err := rfdb.AllKeysChan(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	keys := make([]string, 0)
	for k := range kc {
		fmt.Println(k)
		keys = append(keys, k)
	}
	expectedKeys := make([]string, 0)
	for _, item := range kvdata {
		expectedKeys = append(expectedKeys, item.Key)
	}
	sort.Slice(expectedKeys, func(i, j int) bool {
		return expectedKeys[i] < expectedKeys[j]
	})
	if strings.Join(keys, "") != strings.Join(expectedKeys, "") {
		t.Fatalf("keys from AllKeysChan not match")
	}
	// delete
	for _, item := range kvdata {
		if err := rfdb.Delete(item.Key); err != nil {
			t.Fatal(err)
		}
	}
	// get
	for _, item := range kvdata {
		_, err := rfdb.Get(item.Key)
		t.Log(err)
		if err == nil {
			t.Fatal(err)
		}
	}
}
