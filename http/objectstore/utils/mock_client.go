package utils

import (
	"context"
	dagpoolcli "github.com/filedag-project/filedag-storage/dag/pool/client"
	"github.com/golang/mock/gomock"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-merkledag"
)

func NewMockClient(ctrl *gomock.Controller) dagpoolcli.PoolClient {
	m := dagpoolcli.NewMockPoolClient(ctrl)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "user", "pool,pool123")
	var ci cid.Cid
	var node *merkledag.ProtoNode
	m.EXPECT().Add(gomock.AssignableToTypeOf(ctx), gomock.AssignableToTypeOf(node)).Return(nil).AnyTimes()
	m.EXPECT().Get(gomock.AssignableToTypeOf(ctx), gomock.AssignableToTypeOf(ci)).Return(merkledag.NewRawNode([]byte("123456")), nil).AnyTimes()
	m.EXPECT().Close(gomock.AssignableToTypeOf(context.TODO())).AnyTimes()
	m.EXPECT().Remove(gomock.AssignableToTypeOf(ctx), gomock.AssignableToTypeOf(ci)).Return(nil).AnyTimes()
	return m
}
