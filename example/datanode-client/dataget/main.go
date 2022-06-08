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

//go run main.go dnget --addr=127.0.0.1:9010 --key=5f519eb42bbfac7358812df89186ba8f07aad854383bbb29ef8c48914b62e59e

func main() {
	var addr, key string
	f := flag.NewFlagSet("dnget", flag.ExitOnError)
	f.StringVar(&addr, "addr", "", "the addr of data node server eg.127.0.0.1:9010")
	f.StringVar(&key, "key", "", "the data key ")
	switch os.Args[1] {
	case "dnget":
		f.Parse(os.Args[2:])
		err := get(addr, key)
		if err != nil {
			fmt.Printf("get data err %v", err)
			return
		}
	default:
		fmt.Println("expected 'data get' subcommands")
		os.Exit(1)
	}
}

func get(addr string, key string) error {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		conn.Close()
		log.Errorf("did not connect: %v", err)
		return err
	}
	defer conn.Close()
	client := proto.NewDataNodeClient(conn)
	res, err := client.Get(context.TODO(), &proto.GetRequest{Key: key})
	if err != nil {
		log.Errorf("%s,keyCode:%s,kvdb get :%v", addr, key, err)
		return err
	}
	fmt.Println(string(res.DataBlock))
	return nil
}
