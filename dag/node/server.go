package node

import (
	"context"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/proto"
	"github.com/filedag-project/filedag-storage/kv"
	"github.com/filedag-project/filedag-storage/kv/mutcask"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"net"
)

const (
	Host = "HOST"
	Port = "PORT"
	Path = "PATH"
)

type server struct {
	proto.UnimplementedDataNodeServer
	kvdb kv.KVDB
}

const HealthCheckService = "grpc.health.v1.Health"

func (s *server) Put(ctx context.Context, in *proto.AddRequest) (*proto.AddResponse, error) {
	err := s.kvdb.Put(in.Key, in.DataBlock)
	if err != nil {
		return &proto.AddResponse{Message: "failed"}, err
	}
	return &proto.AddResponse{Message: "success"}, nil
}

func (s *server) Get(ctx context.Context, in *proto.GetRequest) (*proto.GetResponse, error) {
	bytes, err := s.kvdb.Get(in.Key)
	if err != nil {
		return nil, err
	}
	return &proto.GetResponse{
		DataBlock: bytes,
	}, nil
}

func (s *server) Delete(ctx context.Context, in *proto.DeleteRequest) (*proto.DeleteResponse, error) {
	err := s.kvdb.Delete(in.Key)
	if err != nil {
		return &proto.DeleteResponse{Message: "failed"}, err
	}
	return &proto.DeleteResponse{Message: "success"}, nil
}

func (s *server) Size(ctx context.Context, in *proto.SizeRequest) (*proto.SizeResponse, error) {
	size, err := s.kvdb.Size(in.Key)
	if err != nil {
		return nil, err
	}
	return &proto.SizeResponse{
		Size: int64(size),
	}, nil
}

func (s *server) Shutdown() error {
	return s.kvdb.Close()
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

func MutDataNodeServer(host, port, path string) {
	log.Infof("datanode start")
	// listen port
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	//HealthCheck
	hs := health.NewServer()
	hs.SetServingStatus(HealthCheckService, healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(s, hs)

	mutc, err := mutcask.NewMutcask(mutcask.PathConf(path), mutcask.CaskNumConf(6))
	proto.RegisterDataNodeServer(s, &server{kvdb: mutc})
	if err != nil {
		return
	}
	log.Infof("listen:%v:%v", host, port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
