package server

import (
	"context"
	"flag"
	"google.golang.org/grpc"
	"log"
	"time"
)

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func cli() {
	flag.Parse()

	// 建立连接
	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	// 实例化client
	c := NewDagPoolClient(conn)

	// 调用rpc，等待同步响应
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Add(ctx, &AddRequest{Block: []byte("123456")})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Cid)
}
