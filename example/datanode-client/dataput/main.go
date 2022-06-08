//go:build example
// +build example

package main

import (
	"context"
	"crypto/sha256"
	"flag"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/proto"
	"github.com/google/martian/log"
	blocks "github.com/ipfs/go-block-format"
	"google.golang.org/grpc"
	"io/ioutil"
	"os"
)

//go run main.go dnput --addr=127.0.0.1:9010 --file=./main.go

func main() {
	var addr, file string
	f := flag.NewFlagSet("dnput", flag.ExitOnError)
	f.StringVar(&addr, "addr", "", "the addr of data node server eg.127.0.0.1:9010")
	f.StringVar(&file, "file", "", "the file path that you want add")
	switch os.Args[1] {
	case "dnput":
		f.Parse(os.Args[2:])
		err := put(addr, file)
		if err != nil {
			fmt.Printf("put data err %v", err)
			return
		}
	default:
		fmt.Println("expected 'data put' subcommands")
		os.Exit(1)
	}
}

func put(addr string, file string) error {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		conn.Close()
		log.Errorf("did not connect: %v", err)
		return err
	}
	defer conn.Close()
	client := proto.NewDataNodeClient(conn)
	bytes, err := ioutil.ReadFile(file)
	block := blocks.NewBlock(bytes)
	keyCode := sha256String(block.Cid().String())
	fmt.Println("keyCode:", keyCode)
	_, err = client.Put(context.TODO(), &proto.AddRequest{Key: keyCode, DataBlock: block.RawData()})
	if err != nil {
		log.Errorf("%s,keyCode:%s,kvdb put :%v", addr, keyCode, err)
	}
	return nil
}

func sha256String(s string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
}
