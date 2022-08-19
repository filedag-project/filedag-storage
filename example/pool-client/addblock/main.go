//go:build example
// +build example

package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/pool/client"
	dag "github.com/ipfs/go-merkledag"
	"github.com/ipfs/go-unixfs"
	"os"
)

//go run -tags example main.go addblock --addr=127.0.0.1:50001 --pool-user=dagpool --pool-pass=dagpool --block-data="it's a block data"

func main() {
	var addr, clientuser, clientpass, blockData string
	f := flag.NewFlagSet("addblock", flag.ExitOnError)
	f.StringVar(&addr, "addr", "", "the addr of dagpool server eg.127.0.0.1:50001")
	f.StringVar(&clientuser, "pool-user", "", "the pool user")
	f.StringVar(&clientpass, "pool-pass", "", "the pool user pass")
	f.StringVar(&blockData, "block-data", "", "the block data that you want add,size is usually 1m")
	switch os.Args[1] {
	case "addblock":
		f.Parse(os.Args[2:])
		err := add(addr, clientuser, clientpass, blockData)
		if err != nil {
			fmt.Printf("add block err %v", err)
			return
		}
	default:
		fmt.Println("expected 'addblock' subcommands")
		os.Exit(1)
	}
}

func add(addr string, clientuser string, clientpass string, blockdata string) error {
	poolClient, err := client.NewPoolClient(addr, clientuser, clientpass)
	if err != nil {
		fmt.Printf("NewPoolClient err:%v", err)
		return err
	}

	nd := dag.NodeWithData(unixfs.FilePBData([]byte(blockdata), uint64(len(blockdata))))

	err = poolClient.Put(context.TODO(), nd)
	if err != nil {
		fmt.Printf("add block err:%v", err)
		return err
	}
	fmt.Printf("add block success cid:%v\n", nd.Cid().String())
	return nil
}
