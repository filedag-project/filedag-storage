package client

import (
	"context"
	"errors"
	"github.com/filedag-project/filedag-storage/dag/config"
	"github.com/filedag-project/filedag-storage/dag/proto"
	"github.com/filedag-project/filedag-storage/dag/slotsmgr"
	"github.com/filedag-project/filedag-storage/dag/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type dagPoolClusterClient struct {
	DPClusterClient proto.DagPoolClusterClient
	Conn            *grpc.ClientConn
}

//NewPoolClusterClient new a dagPoolClusterClient
func NewPoolClusterClient(addr string) (*dagPoolClusterClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Errorf("did not connect: %v", err)
		return nil, err
	}
	c := proto.NewDagPoolClusterClient(conn)
	return &dagPoolClusterClient{
		DPClusterClient: c,
		Conn:            conn,
	}, nil
}

//Close  the client
func (cli *dagPoolClusterClient) Close(ctx context.Context) {
	cli.Conn.Close()
}

func (cli *dagPoolClusterClient) AddDagNode(ctx context.Context, nodeConfig *config.DagNodeConfig) error {
	nodeInfo := utils.ToProtoDagNodeInfo(nodeConfig)
	_, err := cli.DPClusterClient.AddDagNode(ctx, nodeInfo)
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.Unknown {
			return errors.New(st.Message())
		}
		return err
	}
	return nil
}

func (cli *dagPoolClusterClient) GetDagNode(ctx context.Context, dagNodeName string) (*config.DagNodeConfig, error) {
	nodeInfo, err := cli.DPClusterClient.GetDagNode(ctx, &proto.GetDagNodeReq{Name: dagNodeName})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.Unknown {
			return nil, errors.New(st.Message())
		}
		return nil, err
	}
	return utils.ToDagNodeConfig(nodeInfo), nil
}

func (cli *dagPoolClusterClient) RemoveDagNode(ctx context.Context, dagNodeName string) (*config.DagNodeConfig, error) {
	nodeInfo, err := cli.DPClusterClient.RemoveDagNode(ctx, &proto.RemoveDagNodeReq{Name: dagNodeName})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.Unknown {
			return nil, errors.New(st.Message())
		}
		return nil, err
	}
	return utils.ToDagNodeConfig(nodeInfo), nil
}

func (cli *dagPoolClusterClient) MigrateSlots(ctx context.Context, fromDagNodeName, toDagNodeName string, pairs []slotsmgr.SlotPair) error {
	_, err := cli.DPClusterClient.MigrateSlots(ctx, &proto.MigrateSlotsReq{
		FromDagNodeName: fromDagNodeName,
		ToDagNodeName:   toDagNodeName,
		Pairs:           utils.ToProtoSlotPairs(pairs),
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.Unknown {
			return errors.New(st.Message())
		}
		return err
	}
	return nil
}

func (cli *dagPoolClusterClient) BalanceSlots(ctx context.Context) error {
	_, err := cli.DPClusterClient.BalanceSlots(ctx, &emptypb.Empty{})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.Unknown {
			return errors.New(st.Message())
		}
		return err
	}
	return nil
}

func (cli *dagPoolClusterClient) Status(ctx context.Context) (*proto.StatusReply, error) {
	reply, err := cli.DPClusterClient.Status(ctx, &emptypb.Empty{})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.Unknown {
			return nil, errors.New(st.Message())
		}
		return nil, err
	}
	return reply, nil
}

func (cli *dagPoolClusterClient) RepairDataNode(ctx context.Context, dagNodeName string, fromIndex, repairIndex int) error {
	_, err := cli.DPClusterClient.RepairDataNode(ctx, &proto.RepairDataNodeReq{
		DagNodeName:     dagNodeName,
		FromNodeIndex:   int32(fromIndex),
		RepairNodeIndex: int32(repairIndex),
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.Unknown {
			return errors.New(st.Message())
		}
		return err
	}
	return nil
}
