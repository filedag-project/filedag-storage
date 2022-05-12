package server

import (
	"context"
	"flag"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/node"
	"github.com/filedag-project/filedag-storage/dag/pool"
	"github.com/filedag-project/filedag-storage/dag/pool/dagpooluser"
	"github.com/filedag-project/filedag-storage/dag/pool/userpolicy"
	"github.com/filedrive-team/filehelper/importer"
	pb "github.com/ipfs/go-unixfs/pb"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement DagPoolServer.
type server struct {
	UnimplementedDagPoolServer
	dp *pool.DagPool
}

func (s *server) Add(ctx context.Context, in *AddRequest) (*AddReply, error) {
	data, err := importer.NewDagWithData(in.Block, pb.Data_File, s.dp.CidBuilder)
	if err != nil {
		return &AddReply{Cid: "", Err: err.Error()}, err
	}
	if !s.dp.Iam.CheckUserPolicy(in.User.Username, in.User.Pass, userpolicy.OnlyWrite) {
		return &AddReply{Cid: ""}, err
	}
	err = s.dp.Add(ctx, data)
	if err != nil {
		return &AddReply{Cid: ""}, err
	}
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
	os.Setenv(pool.DagPoolLeveldbPath, "/tmp/leveldb2/pool")

	os.Setenv(pool.DagNodeIpOrPath, "local")

	os.Setenv(pool.DagPoolImporterBatchNum, "4")
	os.Setenv(node.NodeConfigPath, "../config/node_config.json")
	s := grpc.NewServer()
	service, err := pool.NewDagPoolService()
	if err != nil {
		return
	}
	service.Iam.AddUser(dagpooluser.DagPoolUser{
		Username: "test",
		Password: "test",
		Policy:   userpolicy.ReadWrite,
		Capacity: 0,
	})
	RegisterDagPoolServer(s, &server{dp: service})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
