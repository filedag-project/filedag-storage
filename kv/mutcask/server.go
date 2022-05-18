package mutcask

import (
	"context"
	"flag"
	"fmt"
	"github.com/filedag-project/filedag-storage/proto"
	logging "github.com/ipfs/go-log/v2"
	"google.golang.org/grpc"
	"net"
)

var log = logging.Logger("kv")

type server struct {
	proto.UnimplementedMutCaskServer
	mutcask *mutcask
}

func (s *server) Put(ctx context.Context, in *proto.AddRequest) (*proto.AddResponse, error) {
	err := s.mutcask.Put(in.Key, in.DataBlock)
	if err != nil {
		return &proto.AddResponse{Message: "failed"}, err
	}
	return &proto.AddResponse{Message: "success"}, nil
}

func (s *server) Get(ctx context.Context, in *proto.GetRequest) (*proto.GetResponse, error) {
	bytes, err := s.mutcask.Get(in.Key)
	if err != nil {
		return nil, err
	}
	return &proto.GetResponse{
		DataBlock: bytes,
	}, nil
}

func (s *server) Delete(ctx context.Context, in *proto.DeleteRequest) (*proto.DeleteResponse, error) {
	err := s.mutcask.Delete(in.Key)
	if err != nil {
		return &proto.DeleteResponse{Message: "failed"}, err
	}
	return &proto.DeleteResponse{Message: "success"}, nil
}

func (s *server) Size(ctx context.Context, in *proto.SizeRequest) (*proto.SizeResponse, error) {
	size, err := s.mutcask.Size(in.Key)
	if err != nil {
		return nil, err
	}
	return &proto.SizeResponse{
		Size: int64(size),
	}, nil
}

func MutServer(ip, port, addr, heartPort string) {
	flag.Parse()
	// 监听端口
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", ip, port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	mutc, err := NewMutcask(PathConf(addr), CaskNumConf(6))
	proto.RegisterMutCaskServer(s, &server{mutcask: mutc})
	if err != nil {
		return
	}
	log.Infof("listen:%v:%v", ip, port)
	//proto.RegisterMutCaskServer(s,mutc)
	go SendHeartBeat(heartPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
