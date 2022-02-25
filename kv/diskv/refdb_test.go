package diskv

import (
	"bytes"
	"io/ioutil"
	"testing"
)

type kvt struct {
	Key   []byte
	Value []byte
}

var kvdata = []kvt{
	{[]byte("city"), []byte("shanghai")},
	{[]byte("app"), []byte("filedag")},
	{[]byte("protocol"), []byte("ipfs")},
	{[]byte("blockchain"), []byte("filecoin")},
}

func TestRefdb(t *testing.T) {
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
