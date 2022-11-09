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

//go run -tags example main.go auth update --address=127.0.0.1:50001 --root-user=dagpool --root-password=dagpool --username=only-read --new-password=test123 --new-capacity=2000 --new-policy=read-write

func main() {
	var addr, rootuser, rootpass, username, pass, policy string
	var capacity uint64
	f := flag.NewFlagSet("auth", flag.ExitOnError)
	f.StringVar(&addr, "address", "127.0.0.1:50001", "the addr of dagpool server eg.127.0.0.1:50001")
	f.StringVar(&rootuser, "root-user", "dagpool", "set root user")
	f.StringVar(&rootpass, "root-password", "dagpool", "set root password")
	f.StringVar(&username, "username", "", "set the username to update")
	f.StringVar(&pass, "new-password", "", "set the new password")
	f.StringVar(&policy, "new-policy", "", "set the new capacity")
	f.Uint64Var(&capacity, "new-capacity", 0, "set the new policy, enum: only-read, only-write, read-write")
	switch os.Args[2] {
	case "update":
		f.Parse(os.Args[3:])
		err := update(addr, rootuser, rootpass, username, pass, policy, capacity)
		if err != nil {
			fmt.Printf("update user err %v", err)
			return
		}
	default:
		fmt.Println("expected 'update' subcommands")
		os.Exit(1)
	}
}

func update(addr string, rootUser string, rootPassword string, username string, password string, policy string, capacity uint64) error {
	poolClient, err := client.NewPoolClient(addr, rootUser, rootPassword, false)
	if err != nil {
		fmt.Printf("NewPoolClient err:%v", err)
		return err
	}
	if err = poolClient.UpdateUser(context.TODO(), username, password, capacity, policy); err != nil {
		fmt.Printf("update user err:%v", err)
		return err
	}
	return nil
}
