package datanode

import (
	"github.com/filedag-project/filedag-storage/dag/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

//go:generate go run github.com/golang/mock/mockgen -destination=mocks/mock_dnclient.go -package=mocks . DataNodeClient

type DataNodeClient interface {
	proto.DataNodeClient
}

// Client is a node that stores erasure-coded sharded data
type Client struct {
	DataClient  proto.DataNodeClient
	HeartClient healthpb.HealthClient
	RpcAddress  string
	Conn        *grpc.ClientConn
}

// NewClient creates a grpc connection to a slice
func NewClient(rpcAddress string) (datanode *Client, err error) {
	conn, err := grpc.Dial(rpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Errorf("did not connect: %v", err)
		return nil, err
	}
	datanode = &Client{
		DataClient:  proto.NewDataNodeClient(conn),
		HeartClient: healthpb.NewHealthClient(conn),
		RpcAddress:  rpcAddress,
		Conn:        conn,
	}
	return datanode, nil
}
