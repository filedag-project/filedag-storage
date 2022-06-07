//go:build example
// +build example

package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/pool/client"
	"github.com/filedag-project/filedag-storage/dag/proto"
	"os"
)

//go run -tags example main.go getblock --addr=127.0.0.1:50001 --client-user=dagpool --client-pass=dagpool --cid=QmZikYuqANVBRWcbb1zHAHEXzX6CsWbPz2mqRCoy92Jcge

func main() {
	var addr, clientuser, clientpass, cid string
	f := flag.NewFlagSet("getblock", flag.ExitOnError)
	f.StringVar(&addr, "addr", "", "the addr of dagpool server eg.127.0.0.1:50001")
	f.StringVar(&clientuser, "client-user", "", "the client user")
	f.StringVar(&clientpass, "client-pass", "", "the client user pass")
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

func get(addr string, clientuser string, clientpass string, cid string) error {
	poolClient, err := client.NewPoolClient(addr, clientuser, clientpass)
	if err != nil {
		fmt.Printf("NewPoolClient err:%v", err)
		return err
	}
	re, err := poolClient.DPClient.Get(context.TODO(), &proto.GetReq{
		Cid:  cid,
		User: poolClient.User,
	})
	if err != nil {
		fmt.Printf("get block err:%v", err)
		return err
	}
	fmt.Printf("get block succes cid:%v", re.Block)
	return nil
}
