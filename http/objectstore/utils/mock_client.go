package utils

import (
	"context"
	"github.com/filedag-project/filedag-storage/dag/pool/mocks"
	"github.com/filedag-project/filedag-storage/dag/proto"
	"github.com/golang/mock/gomock"
	"github.com/ipfs/go-merkledag"
)

func NewMockDagPoolClient(ctrl *gomock.Controller) proto.DagPoolClient {
	m := mocks.NewMockDagPoolClient(ctrl)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "user", "pool,pool123")
	node := merkledag.NodeWithData([]byte("\b\u0002\u0012\a1234567\u0018\a"))
	var (
		add     proto.AddReq
		addR    proto.AddReply
		get     proto.GetReq
		getR    = proto.GetReply{Block: node.RawData()}
		remove  proto.RemoveReq
		removeR proto.RemoveReply
	)
	m.EXPECT().Add(gomock.AssignableToTypeOf(ctx), gomock.AssignableToTypeOf(&add)).Return(&addR, nil).AnyTimes()
	m.EXPECT().Get(gomock.AssignableToTypeOf(ctx), gomock.AssignableToTypeOf(&get)).Return(&getR, nil).AnyTimes()
	m.EXPECT().Remove(gomock.AssignableToTypeOf(ctx), gomock.AssignableToTypeOf(&remove)).Return(&removeR, nil).AnyTimes()
	return m
}
