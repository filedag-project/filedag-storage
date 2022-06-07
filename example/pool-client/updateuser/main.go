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

//go run -tags example main.go updateuser --addr=127.0.0.1:50001 --client-user=dagpool --client-pass=dagpool --username=wpg --pass=wpg12345 --capacity=1000 --policy=only-read

func main() {
	var addr, clientuser, clientpass, username, pass, policy string
	var capacity uint64
	f := flag.NewFlagSet("updateuser", flag.ExitOnError)
	f.StringVar(&addr, "addr", "", "the addr of dagpool server eg.127.0.0.1:50001")
	f.StringVar(&clientuser, "client-user", "", "the client user")
	f.StringVar(&clientpass, "client-pass", "", "the client user pass")
	f.StringVar(&username, "username", "", "the username")
	f.StringVar(&pass, "pass", "", "the password")
	f.StringVar(&policy, "policy", "", "the policy")
	f.Uint64Var(&capacity, "capacity", 0, "the capacity")
	switch os.Args[1] {
	case "updateuser":
		f.Parse(os.Args[2:])
		err := updateuser(addr, clientuser, clientpass, username, pass, policy, capacity)
		if err != nil {
			fmt.Printf("update user err %v", err)
			return
		}
	default:
		fmt.Println("expected 'updateuser' subcommands")
		os.Exit(1)
	}
}

func updateuser(addr string, clientuser string, clientpass string, username string, pass string, policy string, capacity uint64) error {
	poolClient, err := client.NewPoolClient(addr, clientuser, clientpass)
	if err != nil {
		fmt.Printf("NewPoolClient err:%v", err)
		return err
	}
	re, err := poolClient.DPClient.UpdateUser(context.TODO(), &proto.UpdateUserReq{
		NewUsername: username,
		NewPassword: pass,
		Policy:      policy,
		Capacity:    capacity,
		User:        poolClient.User,
	})
	if err != nil {
		fmt.Printf("update user err:%v", err)
		return err
	}
	fmt.Printf("update user succes:%v", re.Message)
	return nil
}
