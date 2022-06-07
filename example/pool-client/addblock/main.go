//go:build example
// +build example

package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/pool/client"
	"github.com/filedag-project/filedag-storage/dag/proto"
	"io/ioutil"
	"os"
)

//go run -tags example main.go addblock --addr=127.0.0.1:9985 --client-user=dagpool --client-pass=dagpool --filepath=file.txt

func main() {
	var addr, clientuser, clientpass, filepath string
	f := flag.NewFlagSet("addblock", flag.ExitOnError)
	f.StringVar(&addr, "addr", "", "the addr of dagpool server eg.127.0.0.1:50001")
	f.StringVar(&clientuser, "client-user", "", "the client user")
	f.StringVar(&clientpass, "client-pass", "", "the client user pass")
	f.StringVar(&filepath, "filepath", "", "the block path that you want add,size is usually 1m")
	switch os.Args[1] {
	case "addblock":
		f.Parse(os.Args[2:])
		err := add(addr, clientuser, clientpass, filepath)
		if err != nil {
			fmt.Printf("add user err %v", err)
			return
		}
	default:
		fmt.Println("expected 'addblock' subcommands")
		os.Exit(1)
	}
}

func add(addr string, clientuser string, clientpass string, filepath string) error {
	poolClient, err := client.NewPoolClient(addr, clientuser, clientpass)
	if err != nil {
		fmt.Printf("NewPoolClient err:%v", err)
		return err
	}
	f, err := ioutil.ReadFile(filepath)
	re, err := poolClient.DPClient.Add(context.TODO(), &proto.AddReq{
		Block: f,
		User:  poolClient.User,
	})
	if err != nil {
		fmt.Printf("add block err:%v", err)
		return err
	}
	fmt.Printf("add block succes cid:%v", re.Cid)
	return nil
}
