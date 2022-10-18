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

//go run -tags example main.go auth create --address=127.0.0.1:50001 --root-user=dagpool --root-password=dagpool --username=only-read --password=test123 --capacity=1000 --policy=read-only

func main() {
	var addr, rootuser, rootpass, username, pass, policy string
	var capacity uint64
	f := flag.NewFlagSet("auth", flag.ExitOnError)
	f.StringVar(&addr, "address", "127.0.0.1:50001", "the addr of dagpool server eg.127.0.0.1:50001")
	f.StringVar(&rootuser, "root-user", "dagpool", "set root user")
	f.StringVar(&rootpass, "root-password", "dagpool", "set root password")
	f.StringVar(&username, "username", "", "set the username")
	f.StringVar(&pass, "password", "", "set the password")
	f.StringVar(&policy, "policy", "", "set the policy, enum: read-only, write-only, read-write")
	f.Uint64Var(&capacity, "capacity", 0, "set the capacity")
	switch os.Args[2] {
	case "create":
		f.Parse(os.Args[3:])
		err := create(addr, rootuser, rootpass, username, pass, policy, capacity)
		if err != nil {
			fmt.Printf("add user err %v", err)
			return
		}
	default:
		fmt.Println("expected 'create' subcommands")
		os.Exit(1)
	}
}

func create(addr string, rootUser string, rootPassword string, username string, password string, policy string, capacity uint64) error {
	poolClient, err := client.NewPoolClient(addr, rootUser, rootPassword, false)
	if err != nil {
		fmt.Printf("NewPoolClient err:%v", err)
		return err
	}
	if err = poolClient.AddUser(context.TODO(), username, password, capacity, policy); err != nil {
		fmt.Printf("add user err:%v", err)
		return err
	}
	return nil
}
