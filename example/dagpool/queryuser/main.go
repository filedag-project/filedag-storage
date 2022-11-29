//go:build example
// +build example

package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/pool/client"
	"os"
)

//go run -tags example main.go auth query --address=127.0.0.1:50001 --root-user=dagpool --root-password=dagpool --username=only-read

func main() {
	var addr, rootuser, rootpass, username string
	f := flag.NewFlagSet("auth", flag.ExitOnError)
	f.StringVar(&addr, "address", "127.0.0.1:50001", "the addr of dagpool server eg.127.0.0.1:50001")
	f.StringVar(&rootuser, "root-user", "dagpool", "set root user")
	f.StringVar(&rootpass, "root-password", "dagpool", "set root password")
	f.StringVar(&username, "username", "", "set the username to query")
	switch os.Args[2] {
	case "query":
		f.Parse(os.Args[3:])
		err := query(addr, rootuser, rootpass, username)
		if err != nil {
			fmt.Printf("get user err %v", err)
			return
		}
	default:
		fmt.Println("expected 'query' subcommands")
		os.Exit(1)
	}
}

func query(addr string, rootUser string, rootPassword, username string) error {
	poolClient, err := client.NewPoolClient(addr, rootUser, rootPassword, false)
	if err != nil {
		fmt.Printf("NewPoolClient err:%v", err)
		return err
	}
	reply, err := poolClient.QueryUser(context.TODO(), username)
	if err != nil {
		fmt.Printf("get user err:%v", err)
		return err
	}
	fmt.Println("username:", reply.Username, "capacity:", reply.Capacity, "policy", reply.Policy)
	return nil
}
