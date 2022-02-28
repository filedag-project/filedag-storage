package diskv

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"sync"
	"testing"
)

func TestConcurrentWriteSameKey(t *testing.T) {
	var kvdata = []kvt{
		{[]byte("QmTwNzgUFg2kCZ47AmsKUDHwnfAhcGj6TB4mNZcott9zWc"), []byte("1234")},
		{[]byte("QmTwNzgUFg2kCZ47AmsKUDHwnfAhcGj6TB4mNZcott9zWc"), []byte("abcd")},
		{[]byte("QmTwNzgUFg2kCZ47AmsKUDHwnfAhcGj6TB4mNZcott9zWc"), []byte("5678")},
		{[]byte("QmTwNzgUFg2kCZ47AmsKUDHwnfAhcGj6TB4mNZcott9zWc"), []byte("efgh")},
		{[]byte("QmTwNzgUFg2kCZ47AmsKUDHwnfAhcGj6TB4mNZcott9zWc"), []byte("9ijk")},
		{[]byte("QmTwNzgUFg2kCZ47AmsKUDHwnfAhcGj6TB4mNZcott9zWc"), []byte("lmno")},
		{[]byte("QmTwNzgUFg2kCZ47AmsKUDHwnfAhcGj6TB4mNZcott9zWc"), []byte("pqrs")},
		{[]byte("QmTwNzgUFg2kCZ47AmsKUDHwnfAhcGj6TB4mNZcott9zWc"), []byte("tuvw")},
		{[]byte("QmTwNzgUFg2kCZ47AmsKUDHwnfAhcGj6TB4mNZcott9zWc"), []byte("xyz0")},
		{[]byte("QmTwNzgUFg2kCZ47AmsKUDHwnfAhcGj6TB4mNZcott9zWc"), []byte("()_+")},
		{[]byte("QmTwNzgUFg2kCZ47AmsKUDHwnfAhcGj6TB4mNZcott9zWc"), []byte("&^%*")},
		{[]byte("QmTwNzgUFg2kCZ47AmsKUDHwnfAhcGj6TB4mNZcott9zWc"), []byte("!@#$")},
		{[]byte("QmTwNzgUFg2kCZ47AmsKUDHwnfAhcGj6TB4mNZcott9zWc"), []byte("%<>/")},
		{[]byte("QmTwNzgUFg2kCZ47AmsKUDHwnfAhcGj6TB4mNZcott9zWc"), []byte(":;',")},
		{[]byte("QmTwNzgUFg2kCZ47AmsKUDHwnfAhcGj6TB4mNZcott9zWc"), []byte("◣◥◢◣")},
		{[]byte("QmTwNzgUFg2kCZ47AmsKUDHwnfAhcGj6TB4mNZcott9zWc"), []byte("➔➘➙➚")},
	}
	dkv, err := NewDisKV(func(cfg *Config) {
		cfg.Dir = tmpdirpath(t)
		cfg.MaxLinkDagSize = 1
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
				t.Fail()
			}
		}(item.Key, item.Value)
	}
	wg.Wait()

	_, err = dkv.Get(kvdata[0].Key)
	if err != nil {
		t.Fail()
	}
	err = dkv.Delete(kvdata[0].Key)
	if err != nil {
		t.Fail()
	}
}

func TestDisKVPutDataDags(t *testing.T) {
	var kvdata = []kvt{
		{[]byte("Qmc35RPEYrW3Mj1mki6thkAjx6a1ZFkU3UYxAyFhMmngr2"), []byte("124567")},
		{[]byte("QmTwNzgUFg2kCZ47AmsKUDHwnfAhcGj6TB4mNZcott9zWc"), []byte("224567")},
		{[]byte("QmYgPV5bT37u56qePZUqLQ15JhnopaSmVx8ao39RUCoZEj"), []byte("324567")},
		{[]byte("QmfVM2KjyzYYRn3geYnqv6EWqSwRZAPpdFcgEhc61ycJRp"), []byte("424567")},
		{[]byte("QmQCTP2mVjwerHuM9CwuHqFEvo9w2BEkmnFNfGvThX5Rai"), []byte("524567")},
		{[]byte("QmeioJd3d9LT2f96VH94WU62AFsB1S1V1qq8sGt7A8L9vN"), []byte("624567")},
		{[]byte("QmPXQHq2un3E4cFsYsGukwYWJs7BrBmm3wNauMuw6EqZMa"), []byte("724567")},
		{[]byte("QmU4tBqMdUe94C3D5wsbe7j6ZP6EboMSRTXdyaxRUb4HQz"), []byte("824567")},
		{[]byte("QmWpN6NyLGpgiUdiy6CZ1AZEhrz9guLDb7iJMupk5LWS9y"), []byte("924567")},
		{[]byte("QmXKztBnVXL6dYzSqDt7pRN67fyK7SiqNLXMvvcK5cjdMc"), []byte("134567")},
		{[]byte("QmZQoGSaHXmJJTchBrqBVQgTJ6nL1mYbR4CDhJBpkeK7Fb"), []byte("278934")},
		{[]byte("QmRoRtbKjZiYqr5yvB6fjTudqrKrwsPkJ9XMfMDzzdGsVK"), []byte("378934")},
		{[]byte("Qmbc3FwKnE36uvL9e44yCwFyKifV5BSZ74t9V2m659Xvg5"), []byte("478934")},
		{[]byte("QmStSiCG7rgDgNU6g1bBdK8jbBBaVtiqRzVgHYQYN2wKWo"), []byte("578934")},
		{[]byte("QmW6EVWYvFEMHFErio7nTU3DhRrjHZn4ednRkHj2fSpTm7"), []byte("678934")},
		{[]byte("QmUPqWa9KJz44skxo8fDD4UFcxTsbTLk2XWQd1HdTdBq1h"), []byte("778934")},
		{[]byte("QmSfVC3EX4Uwa54sJt8F9TFuWDVvRzCbyuxpfDdh6qMgwR"), []byte("878934")},
		{[]byte("QmWBwR7pC2VY9KcFXgLJSYGZrbwuTnpNYHizgfDrtVMPCH"), []byte("978934")},
		{[]byte("QmapgjbPdMSqz6qTWGHesRuzBjQk9btZKSEMzZuEm2BKXt"), []byte("139836")},
		{[]byte("QmYfqhMnqunMjPFYsnUJea8sN65LFmF8ChSZ7kivZiwXi7"), []byte("239836")},
		{[]byte("QmRJvXuzSFRq5Sajd8hesZLsXnaWYe5bScsjWZUj1NLkgz"), []byte("339836")},
		{[]byte("QmQ1xczV6i2GzWv7RnstCs5ThyS9ngTadWiyGLZnBQD4Ry"), []byte("439836")},
		{[]byte("QmYHcpDZAzAW4N8gYDecNDAvk9gpwmPCMJKSCm7U1Eyvna"), []byte("539836")},
		{[]byte("QmQmXdRBn5zRVmq6ZBVS1tFKe3sBf8xuXibzqH7zqi2hp1"), []byte("639836")},
		{[]byte("QmZYCXLAV3wdpiWDfggRnC6ndKboedceDqnJDGqkuDBp3z"), []byte("739836")},
		{[]byte("QmXgEMNz5JbajkQ8tXRJHgbC12aogba9gwTgqTQW2LCK35"), []byte("839836")},
		{[]byte("QmW6esdA2tsRmoiqmAgNx71vdNNtgJEd44CKt4nncUTsur"), []byte("939836")},
	}

	dkv, err := NewDisKV(func(cfg *Config) {
		cfg.Dir = tmpdirpath(t)
		cfg.MaxLinkDagSize = 1
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
				t.Failed()
			}
		}(item.Key, item.Value)
	}
	wg.Wait()

	wg.Add(len(kvdata))
	for _, item := range kvdata {
		go func(k, v []byte) {
			defer wg.Done()

			b, err := dkv.Get(k)
			if err != nil {
				fmt.Println(err)
				t.Fail()
			}
			if !bytes.Equal(b, v) {
				fmt.Printf("%s should equal to %s \n", b, v)
				t.Fail()
			}
		}(item.Key, item.Value)
	}
	wg.Wait()

	wg.Add(len(kvdata))
	for _, item := range kvdata {
		go func(k, v []byte) {
			defer wg.Done()

			err := dkv.Delete(k)
			if err != nil {
				t.Fail()
			}
		}(item.Key, item.Value)
	}
	wg.Wait()

	wg.Add(len(kvdata))
	for _, item := range kvdata {
		go func(k, v []byte) {
			defer wg.Done()

			_, err := dkv.Get(k)
			if err != ErrNotFound {
				t.Fail()
			}

		}(item.Key, item.Value)
	}
	wg.Wait()
}

func TestDisKVPutLinkDags(t *testing.T) {
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
		cfg.Dir = tmpdirpath(t)
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
}

func tmpdirpath(t *testing.T) string {
	tmpdir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("failed to make temp dir: %s", err)
	}
	return tmpdir
}
