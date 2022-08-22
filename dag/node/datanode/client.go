package datanode

import (
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/config"
	"github.com/filedag-project/filedag-storage/dag/proto"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

//go:generate go run github.com/golang/mock/mockgen -destination=mocks/mock_dnclient.go -package=mocks . DataNodeClient

type DataNodeClient interface {
	proto.DataNodeClient
}

//Client is a node that stores erasure-coded sharded data
type Client struct {
	Client      proto.DataNodeClient
	HeartClient healthpb.HealthClient
	Ip          string
	Port        string
	Conn        *grpc.ClientConn
}

//NewClient creates a grpc connection to a slice
func NewClient(cfg config.DataNodeConfig) (datanode *Client, err error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", cfg.Ip, cfg.Port), grpc.WithInsecure())
	if err != nil {
		log.Errorf("did not connect: %v", err)
		return nil, err
	}
	datanode = &Client{
		Client:      proto.NewDataNodeClient(conn),
		HeartClient: healthpb.NewHealthClient(conn),
		Ip:          cfg.Ip,
		Port:        cfg.Port,
		Conn:        conn,
	}
	return datanode, nil
}
