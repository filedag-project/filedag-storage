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
		{"QmTwNzgUFg2kCZ47AmsKUDHwnfAhcGj6TB4mNZcott9zWc", []byte("1234")},
		{"QmTwNzgUFg2kCZ47AmsKUDHwnfAhcGj6TB4mNZcott9zWc", []byte("abcd")},
		{"QmTwNzgUFg2kCZ47AmsKUDHwnfAhcGj6TB4mNZcott9zWc", []byte("5678")},
		{"QmTwNzgUFg2kCZ47AmsKUDHwnfAhcGj6TB4mNZcott9zWc", []byte("efgh")},
		{"QmTwNzgUFg2kCZ47AmsKUDHwnfAhcGj6TB4mNZcott9zWc", []byte("9ijk")},
		{"QmTwNzgUFg2kCZ47AmsKUDHwnfAhcGj6TB4mNZcott9zWc", []byte("lmno")},
		{"QmTwNzgUFg2kCZ47AmsKUDHwnfAhcGj6TB4mNZcott9zWc", []byte("pqrs")},
		{"QmTwNzgUFg2kCZ47AmsKUDHwnfAhcGj6TB4mNZcott9zWc", []byte("tuvw")},
		{"QmTwNzgUFg2kCZ47AmsKUDHwnfAhcGj6TB4mNZcott9zWc", []byte("xyz0")},
		{"QmTwNzgUFg2kCZ47AmsKUDHwnfAhcGj6TB4mNZcott9zWc", []byte("()_+")},
		{"QmTwNzgUFg2kCZ47AmsKUDHwnfAhcGj6TB4mNZcott9zWc", []byte("&^%*")},
		{"QmTwNzgUFg2kCZ47AmsKUDHwnfAhcGj6TB4mNZcott9zWc", []byte("!@#$")},
		{"QmTwNzgUFg2kCZ47AmsKUDHwnfAhcGj6TB4mNZcott9zWc", []byte("%<>/")},
		{"QmTwNzgUFg2kCZ47AmsKUDHwnfAhcGj6TB4mNZcott9zWc", []byte(":;',")},
		{"QmTwNzgUFg2kCZ47AmsKUDHwnfAhcGj6TB4mNZcott9zWc", []byte("◣◥◢◣")},
		{"QmTwNzgUFg2kCZ47AmsKUDHwnfAhcGj6TB4mNZcott9zWc", []byte("➔➘➙➚")},
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
		go func(k string, v []byte) {
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
		{"Qmc35RPEYrW3Mj1mki6thkAjx6a1ZFkU3UYxAyFhMmngr2", []byte("124567")},
		{"QmTwNzgUFg2kCZ47AmsKUDHwnfAhcGj6TB4mNZcott9zWc", []byte("224567")},
		{"QmYgPV5bT37u56qePZUqLQ15JhnopaSmVx8ao39RUCoZEj", []byte("324567")},
		{"QmfVM2KjyzYYRn3geYnqv6EWqSwRZAPpdFcgEhc61ycJRp", []byte("424567")},
		{"QmQCTP2mVjwerHuM9CwuHqFEvo9w2BEkmnFNfGvThX5Rai", []byte("524567")},
		{"QmeioJd3d9LT2f96VH94WU62AFsB1S1V1qq8sGt7A8L9vN", []byte("624567")},
		{"QmPXQHq2un3E4cFsYsGukwYWJs7BrBmm3wNauMuw6EqZMa", []byte("724567")},
		{"QmU4tBqMdUe94C3D5wsbe7j6ZP6EboMSRTXdyaxRUb4HQz", []byte("824567")},
		{"QmWpN6NyLGpgiUdiy6CZ1AZEhrz9guLDb7iJMupk5LWS9y", []byte("924567")},
		{"QmXKztBnVXL6dYzSqDt7pRN67fyK7SiqNLXMvvcK5cjdMc", []byte("134567")},
		{"QmZQoGSaHXmJJTchBrqBVQgTJ6nL1mYbR4CDhJBpkeK7Fb", []byte("278934")},
		{"QmRoRtbKjZiYqr5yvB6fjTudqrKrwsPkJ9XMfMDzzdGsVK", []byte("378934")},
		{"Qmbc3FwKnE36uvL9e44yCwFyKifV5BSZ74t9V2m659Xvg5", []byte("478934")},
		{"QmStSiCG7rgDgNU6g1bBdK8jbBBaVtiqRzVgHYQYN2wKWo", []byte("578934")},
		{"QmW6EVWYvFEMHFErio7nTU3DhRrjHZn4ednRkHj2fSpTm7", []byte("678934")},
		{"QmUPqWa9KJz44skxo8fDD4UFcxTsbTLk2XWQd1HdTdBq1h", []byte("778934")},
		{"QmSfVC3EX4Uwa54sJt8F9TFuWDVvRzCbyuxpfDdh6qMgwR", []byte("878934")},
		{"QmWBwR7pC2VY9KcFXgLJSYGZrbwuTnpNYHizgfDrtVMPCH", []byte("978934")},
		{"QmapgjbPdMSqz6qTWGHesRuzBjQk9btZKSEMzZuEm2BKXt", []byte("139836")},
		{"QmYfqhMnqunMjPFYsnUJea8sN65LFmF8ChSZ7kivZiwXi7", []byte("239836")},
		{"QmRJvXuzSFRq5Sajd8hesZLsXnaWYe5bScsjWZUj1NLkgz", []byte("339836")},
		{"QmQ1xczV6i2GzWv7RnstCs5ThyS9ngTadWiyGLZnBQD4Ry", []byte("439836")},
		{"QmYHcpDZAzAW4N8gYDecNDAvk9gpwmPCMJKSCm7U1Eyvna", []byte("539836")},
		{"QmQmXdRBn5zRVmq6ZBVS1tFKe3sBf8xuXibzqH7zqi2hp1", []byte("639836")},
		{"QmZYCXLAV3wdpiWDfggRnC6ndKboedceDqnJDGqkuDBp3z", []byte("739836")},
		{"QmXgEMNz5JbajkQ8tXRJHgbC12aogba9gwTgqTQW2LCK35", []byte("839836")},
		{"QmW6esdA2tsRmoiqmAgNx71vdNNtgJEd44CKt4nncUTsur", []byte("939836")},
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
		go func(k string, v []byte) {
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
		go func(k string, v []byte) {
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
		go func(k string, v []byte) {
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
		go func(k string, v []byte) {
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
		{"city", []byte("shanghai")},
		{"app", []byte("filedag")},
		{"protocol", []byte("ipfs")},
		{"blockchain", []byte("filecoin")},
		{"fs", []byte("nfs")},
		{"kv", []byte("leveldb")},
		{"db", []byte("mongodb")},
		{"language", []byte("golang")},
		{"b", []byte("1")},
		{"c", []byte("2")},
		{"d", []byte("3")},
		{"e", []byte("4")},
		{"f", []byte("5")},
		{"g", []byte("6")},
		{"h", []byte("7")},
		{"i", []byte("8")},
		{"j", []byte("9")},
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
		go func(k string, v []byte) {
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
