package mutcask

import (
	"context"
	"flag"
	"fmt"
	"github.com/filedag-project/filedag-storage/proto"
	logging "github.com/ipfs/go-log/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"os"
)

var log = logging.Logger("kv")

const (
	Host = "HOST"
	Port = "PORT"
	Path = "PATH"
)

type server struct {
	proto.UnimplementedMutCaskServer
	mutcask *mutcask
}

const HealthCheckService = "grpc.health.v1.Health"

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

//func (s *server) Check(ctx context.Context, in *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
//	return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}, nil
//}
//func (s *server) Watch(in *healthpb.HealthCheckRequest, w healthpb.Health_WatchServer) error {
//	err := w.Send(&healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING})
//	if err != nil {
//		return err
//	}
//	return nil
//}

func MutServer() {
	fmt.Println("mut cask start")
	flag.Parse()
	// 监听端口
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", os.Getenv(Host), os.Getenv(Port)))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	//HealthCheck
	hs := health.NewServer()
	hs.SetServingStatus(HealthCheckService, healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(s, hs)

	mutc, err := NewMutcask(PathConf(os.Getenv(Path)), CaskNumConf(6))
	proto.RegisterMutCaskServer(s, &server{mutcask: mutc})
	if err != nil {
		return
	}
	log.Infof("listen:%v:%v", os.Getenv(Host), os.Getenv(Port))
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
