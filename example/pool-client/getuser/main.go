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

//go run -tags example main.go getuser --addr=127.0.0.1:50001 --pool-user=test --pool-pass=test123

func main() {
	var addr, clientuser, clientpass string
	f := flag.NewFlagSet("getuser", flag.ExitOnError)
	f.StringVar(&addr, "addr", "", "the addr of dagpool server eg.127.0.0.1:50001")
	f.StringVar(&clientuser, "pool-user", "", "the pool user")
	f.StringVar(&clientpass, "pool-pass", "", "the pool user pass")
	switch os.Args[1] {
	case "getuser":
		f.Parse(os.Args[2:])
		err := getuser(addr, clientuser, clientpass)
		if err != nil {
			fmt.Printf("get user err %v", err)
			return
		}
	default:
		fmt.Println("expected 'getuser' subcommands")
		os.Exit(1)
	}
}

func getuser(addr string, clientuser string, clientpass string) error {
	poolClient, err := client.NewPoolClient(addr, clientuser, clientpass)
	if err != nil {
		fmt.Printf("NewPoolClient err:%v", err)
		return err
	}
	re, err := poolClient.DPClient.QueryUser(context.TODO(), &proto.QueryUserReq{
		User: &proto.PoolUser{
			User:     clientuser,
			Password: clientpass,
		},
		Username: clientuser,
	})
	if err != nil {
		fmt.Printf("get user err:%v", err)
		return err
	}
	fmt.Printf("get user:%v succes,policy:%v,capacity:%v", re.Username, re.Policy, re.Capacity)
	return nil
}
