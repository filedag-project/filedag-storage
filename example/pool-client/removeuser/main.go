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

//go run -tags example main.go removeuser --addr=127.0.0.1:50001 --pool-user=dagpool --pool-pass=dagpool --username=test

func main() {
	var addr, clientuser, clientpass, username string
	f := flag.NewFlagSet("removeuser", flag.ExitOnError)
	f.StringVar(&addr, "addr", "", "the addr of dagpool server eg.127.0.0.1:50001")
	f.StringVar(&clientuser, "pool-user", "", "the pool user")
	f.StringVar(&clientpass, "pool-pass", "", "the pool user pass")
	f.StringVar(&username, "username", "", "the username")
	switch os.Args[1] {
	case "removeuser":
		f.Parse(os.Args[2:])
		err := removeuser(addr, clientuser, clientpass, username)
		if err != nil {
			fmt.Printf("remove user err %v", err)
			return
		}
	default:
		fmt.Println("expected 'removeuser' subcommands")
		os.Exit(1)
	}
}

func removeuser(addr string, clientuser string, clientpass string, username string) error {
	poolClient, err := client.NewPoolClient(addr, clientuser, clientpass)
	if err != nil {
		fmt.Printf("NewPoolClient err:%v", err)
		return err
	}
	re, err := poolClient.DPClient.RemoveUser(context.TODO(), &proto.RemoveUserReq{
		User:     poolClient.User,
		Username: username,
	})
	if err != nil {
		fmt.Printf("remove user err:%v", err)
		return err
	}
	fmt.Printf("remove user:%v succes %v", username, re.Message)
	return nil
}
