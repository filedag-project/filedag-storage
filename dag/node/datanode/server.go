package datanode

import (
	"bytes"
	"context"
	"encoding/binary"
	"github.com/filedag-project/filedag-storage/dag/proto"
	"github.com/filedag-project/filedag-storage/kv"
	"github.com/filedag-project/filedag-storage/kv/badger"
	"github.com/filedag-project/filedag-storage/kv/mutcask"
	"github.com/howeyc/crc16"
	logging "github.com/ipfs/go-log/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
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

	// HeaderSize is size of entry header
	HeaderSize = 12
)

// entry item
// | crc (4 bytes) | meta size (4 bytes) | data size (4 bytes) | meta | data |

type server struct {
	proto.UnimplementedDataNodeServer
	kvdb kv.KVDB
}

const healthCheckService = "grpc.health.v1.Health"

type Header struct {
	Checksum uint32
	MetaSize int32
	DataSize int32
}

//Put puts the data by key
func (s *server) Put(ctx context.Context, in *proto.AddRequest) (*emptypb.Empty, error) {
	header := Header{
		MetaSize: int32(len(in.Meta)),
		DataSize: int32(len(in.Data)),
	}
	var buf bytes.Buffer
	buf.Grow(binary.Size(header) + int(header.MetaSize+header.DataSize))
	if err := binary.Write(&buf, binary.LittleEndian, header); err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	buf.Write(in.Meta)
	buf.Write(in.Data)
	data := buf.Bytes()
	header.Checksum = uint32(crc16.Checksum(data[binary.Size(header.Checksum):], crc16.IBMTable))
	var bufCrc bytes.Buffer
	if err := binary.Write(&bufCrc, binary.LittleEndian, header.Checksum); err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	copy(data, bufCrc.Bytes())
	if err := s.kvdb.Put(in.Key, data); err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	return &emptypb.Empty{}, nil
}

//Get gets the data by key
func (s *server) Get(ctx context.Context, in *proto.GetRequest) (*proto.GetResponse, error) {
	data, err := s.kvdb.Get(in.Key)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	buf := bytes.NewBuffer(data)
	header := Header{}
	if err = binary.Read(buf, binary.LittleEndian, &header); err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	// check crc
	sum := crc16.Checksum(data[binary.Size(header.Checksum):], crc16.IBMTable)
	if header.Checksum != uint32(sum) {
		return nil, status.Error(codes.Unknown, "checking crc failed")
	}
	offset := binary.Size(header)
	return &proto.GetResponse{
		Meta: data[offset : offset+int(header.MetaSize)],
		Data: data[offset+int(header.MetaSize) : offset+int(header.MetaSize+header.DataSize)],
	}, nil
}

func (s *server) GetMeta(ctx context.Context, in *proto.GetMetaRequest) (*proto.GetMetaResponse, error) {
	data, err := s.kvdb.Get(in.Key)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	buf := bytes.NewBuffer(data)
	header := Header{}
	if err = binary.Read(buf, binary.LittleEndian, &header); err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	// check crc
	sum := crc16.Checksum(data[binary.Size(header.Checksum):], crc16.IBMTable)
	if header.Checksum != uint32(sum) {
		return nil, status.Error(codes.Unknown, "checking crc failed")
	}
	headerSize := binary.Size(header)
	return &proto.GetMetaResponse{
		Meta: data[headerSize : headerSize+int(header.MetaSize)],
	}, nil
}

//Delete deletes the data by key
func (s *server) Delete(ctx context.Context, in *proto.DeleteRequest) (*emptypb.Empty, error) {
	err := s.kvdb.Delete(in.Key)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	return &emptypb.Empty{}, nil
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

func (s *server) DeleteMany(ctx context.Context, in *proto.DeleteManyRequest) (*emptypb.Empty, error) {
	for _, key := range in.Keys {
		err := s.kvdb.Delete(key)
		if err != nil {
			return nil, status.Error(codes.Unknown, err.Error())
		}
	}
	return &emptypb.Empty{}, nil
}

func (s *server) AllKeysChan(_ *emptypb.Empty, server proto.DataNode_AllKeysChanServer) error {
	ctx, cancel := context.WithCancel(server.Context())
	defer cancel()
	ch, err := s.kvdb.AllKeysChan(ctx)
	if err != nil {
		return status.Error(codes.Unknown, err.Error())
	}
	for key := range ch {
		if err = server.Send(&proto.AllKeysChanResponse{Key: key}); err != nil {
			return status.Error(codes.Unknown, err.Error())
		}
	}
	return nil
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

//StartDataNodeServer is the gRPC server for the MutDataNode
func StartDataNodeServer(listen string, kvType KVType, dataDir string) {
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
