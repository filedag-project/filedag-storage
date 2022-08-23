//go:build example
// +build example

package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/pool/client"
	"github.com/ipfs/go-cid"
	"os"
)

//go run -tags example main.go getblock --addr=127.0.0.1:50001 --pool-user=dagpool --pool-pass=dagpool --cid=QmaR7tvZDJgvdXBx59Wf7s1GZRDL1Lqv5ivJDJyUGaHvBY

func main() {
	var addr, clientuser, clientpass, cid string
	f := flag.NewFlagSet("getblock", flag.ExitOnError)
	f.StringVar(&addr, "addr", "", "the addr of dagpool server eg.127.0.0.1:50001")
	f.StringVar(&clientuser, "pool-user", "", "the pool user")
	f.StringVar(&clientpass, "pool-pass", "", "the pool user pass")
	f.StringVar(&cid, "cid", "", "the block cid")
	switch os.Args[1] {
	case "getblock":
		f.Parse(os.Args[2:])
		err := get(addr, clientuser, clientpass, cid)
		if err != nil {
			fmt.Printf("get block err %v", err)
			return
		}
	default:
		fmt.Println("expected 'getblock' subcommands")
		os.Exit(1)
	}
}

func get(addr string, clientuser string, clientpass string, cidStr string) error {
	poolClient, err := client.NewPoolClient(addr, clientuser, clientpass)
	if err != nil {
		fmt.Printf("NewPoolClient err:%v", err)
		return err
	}
	c, err := cid.Decode(cidStr)
	if err != nil {
		return err
	}
	nd, err := poolClient.Get(context.TODO(), c)
	if err != nil {
		fmt.Printf("get block err:%v", err)
		return err
	}
	cidSum, err := c.Prefix().Sum(nd.RawData())
	if err != nil {
		fmt.Printf("check block err:%v", err)
		return err
	}
	if c == cidSum {
		fmt.Println("check block success")
	} else {
		fmt.Println("check block failed")
	}
	fmt.Println("block data: ", string(nd.RawData()))
	return nil
}
