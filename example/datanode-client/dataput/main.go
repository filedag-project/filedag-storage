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

//go run main.go dnput --addr=127.0.0.1:9010 --key=test --data="it's test content"

func main() {
	var addr, key, data string
	f := flag.NewFlagSet("dnput", flag.ExitOnError)
	f.StringVar(&addr, "addr", "", "the addr of data node server eg.127.0.0.1:9010")
	f.StringVar(&key, "key", "", "the data key")
	f.StringVar(&data, "data", "", "the data content")
	switch os.Args[1] {
	case "dnput":
		f.Parse(os.Args[2:])
		err := put(addr, key, data)
		if err != nil {
			fmt.Printf("put data err %v", err)
			return
		}
	default:
		fmt.Println("expected 'data put' subcommands")
		os.Exit(1)
	}
}

func put(addr string, key, data string) error {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		conn.Close()
		log.Errorf("did not connect: %v", err)
		return err
	}
	defer conn.Close()
	client := proto.NewDataNodeClient(conn)
	resp, err := client.Put(context.TODO(), &proto.AddRequest{Key: key, DataBlock: []byte(data)})
	if err != nil {
		log.Errorf("%s,keyCode:%s,kvdb put :%v", addr, key, err)
		return err
	}
	fmt.Println(resp.Message)
	return nil
}
