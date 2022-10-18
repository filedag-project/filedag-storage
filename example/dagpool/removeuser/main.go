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

//go run -tags example main.go auth remove --address=127.0.0.1:50001 --root-user=dagpool --root-password=dagpool --username=only-read

func main() {
	var addr, rootuser, rootpass, username string
	f := flag.NewFlagSet("auth", flag.ExitOnError)
	f.StringVar(&addr, "address", "127.0.0.1:50001", "the addr of dagpool server eg.127.0.0.1:50001")
	f.StringVar(&rootuser, "root-user", "dagpool", "set root user")
	f.StringVar(&rootpass, "root-password", "dagpool", "set root password")
	f.StringVar(&username, "username", "", "set the username to remove")
	switch os.Args[2] {
	case "remove":
		f.Parse(os.Args[3:])
		err := remove(addr, rootuser, rootpass, username)
		if err != nil {
			fmt.Printf("remove user err %v", err)
			return
		}
	default:
		fmt.Println("expected 'remove' subcommands")
		os.Exit(1)
	}
}

func remove(addr string, rootUser string, rootPassword string, username string) error {
	poolClient, err := client.NewPoolClient(addr, rootUser, rootPassword, false)
	if err != nil {
		fmt.Printf("NewPoolClient err:%v", err)
		return err
	}
	if err = poolClient.RemoveUser(context.TODO(), username); err != nil {
		fmt.Printf("remove user err:%v", err)
		return err
	}
	return nil
}
