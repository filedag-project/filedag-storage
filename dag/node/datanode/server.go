package datanode

import (
	"context"
	"github.com/filedag-project/filedag-storage/dag/proto"
	"github.com/filedag-project/filedag-storage/kv"
	"github.com/filedag-project/filedag-storage/kv/badger"
	"github.com/filedag-project/filedag-storage/kv/mutcask"
	logging "github.com/ipfs/go-log/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var log = logging.Logger("datanode")

//KVType is the type of kv
type KVType string

const (
	//KVBadge is the kv type of badger
	KVBadge KVType = "badger"
	//KVMutcask is the kv type of mutcask
	KVMutcask KVType = "mutcask"
)

type server struct {
	proto.UnimplementedDataNodeServer
	kvdb kv.KVDB
}

const healthCheckService = "grpc.health.v1.Health"

//Put puts the data by key
func (s *server) Put(ctx context.Context, in *proto.AddRequest) (*proto.AddResponse, error) {
	err := s.kvdb.Put(in.Key, in.DataBlock)
	if err != nil {
		return &proto.AddResponse{Message: "failed"}, err
	}
	return &proto.AddResponse{Message: "success"}, nil
}

//Get gets the data by key
func (s *server) Get(ctx context.Context, in *proto.GetRequest) (*proto.GetResponse, error) {
	bytes, err := s.kvdb.Get(in.Key)
	if err != nil {
		return nil, err
	}
	return &proto.GetResponse{
		DataBlock: bytes,
	}, nil
}

//Delete deletes the data by key
func (s *server) Delete(ctx context.Context, in *proto.DeleteRequest) (*proto.DeleteResponse, error) {
	err := s.kvdb.Delete(in.Key)
	if err != nil {
		return &proto.DeleteResponse{Message: "failed"}, err
	}
	return &proto.DeleteResponse{Message: "success"}, nil
}

//Size  returns the size of data by key
func (s *server) Size(ctx context.Context, in *proto.SizeRequest) (*proto.SizeResponse, error) {
	size, err := s.kvdb.Size(in.Key)
	if err != nil {
		return nil, err
	}
	return &proto.SizeResponse{
		Size: int64(size),
	}, nil
}

func (s *server) DeleteMany(ctx context.Context, in *proto.DeleteManyRequest) (*proto.DeleteManyResponse, error) {
	for _, key := range in.Keys {
		err := s.kvdb.Delete(key)
		if err != nil {
			return &proto.DeleteManyResponse{Message: "failed"}, err
		}
	}
	return &proto.DeleteManyResponse{Message: "success"}, nil
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

//MutDataNodeServer is the gRPC server for the MutDataNode
func MutDataNodeServer(listen string, kvType KVType, dataDir string) {
	log.Infof("datanode start...")
	log.Infof("listen %s", listen)
	// listen port
	lis, err := net.Listen("tcp", listen)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	//HealthCheck
	hs := health.NewServer()
	hs.SetServingStatus(healthCheckService, healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(s, hs)

	if err := os.MkdirAll(dataDir, 0777); err != nil {
		log.Fatalf("failed to create directory: %v", err)
	}

	var kvdb kv.KVDB
	switch kvType {
	case KVBadge:
		kvdb, err = badger.NewBadger(dataDir)
	case KVMutcask:
		kvdb, err = mutcask.NewMutcask(mutcask.PathConf(dataDir), mutcask.CaskNumConf(6))
	default:
		log.Fatal("not handle this kv type")
	}
	if err != nil {
		log.Fatalf("failed to load db: %v", err)
	}
	defer kvdb.Close()

	proto.RegisterDataNodeServer(s, &server{kvdb: kvdb})
	if err != nil {
		return
	}
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutdown Server ...")

	s.GracefulStop()

	log.Info("Server exit")
}
