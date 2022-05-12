package server

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	UnimplementedDagPoolServer
}

func (s *server) Add(ctx context.Context, in *AddRequest) (*AddReply, error) {
	return &AddReply{Cid: "Hello" + string(in.GetBlock())}, nil
}

func ser() {
	flag.Parse()
	// 监听端口
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// 实例化server
	s := grpc.NewServer()
	RegisterDagPoolServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
