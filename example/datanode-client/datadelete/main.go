//go:build example
// +build example

package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/proto"
	"github.com/google/martian/log"
	"google.golang.org/grpc"
	"os"
)

//go run main.go dndelete --addr=127.0.0.1:9010 --key=test

func main() {
	var addr, key string
	f := flag.NewFlagSet("dndelete", flag.ExitOnError)
	f.StringVar(&addr, "addr", "", "the addr of data node server eg.127.0.0.1:9010")
	f.StringVar(&key, "key", "", "the data key ")
	switch os.Args[1] {
	case "dndelete":
		f.Parse(os.Args[2:])
		err := delete(addr, key)
		if err != nil {
			fmt.Printf("delete data err %v", err)
			return
		}
	default:
		fmt.Println("expected 'data delete' subcommands")
		os.Exit(1)
	}
}

func delete(addr string, key string) error {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		conn.Close()
		log.Errorf("did not connect: %v", err)
		return err
	}
	defer conn.Close()
	client := proto.NewDataNodeClient(conn)
	_, err = client.Delete(context.TODO(), &proto.DeleteRequest{Key: key})
	if err != nil {
		log.Errorf("%s,keyCode:%s,kvdb delete :%v", addr, key, err)
		return err
	}
	return nil
}
