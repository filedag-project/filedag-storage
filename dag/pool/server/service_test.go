package server

import (
	"context"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/proto"
	"github.com/golang/mock/gomock"
	"github.com/ipfs/go-merkledag"
	"testing"
)

//NewMockPoolClient creates a mock of PoolClient
func NewMockDagPoolServer(t *testing.T) (*DagPoolServer, func()) {
	ctrl := gomock.NewController(t)
	m := NewMockDagPool(ctrl)
	node := merkledag.NodeWithData([]byte("\b\u0002\u0012\a1234567\u0018\a"))
	m.EXPECT().Add(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	m.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(node, nil).AnyTimes()
	return &DagPoolServer{DagPool: m}, ctrl.Finish
}
func TestDagPoolServer(t *testing.T) {
	// new server
	ser, _ := NewMockDagPoolServer(t)
	user1 := &proto.PoolUser{User: "user1", Password: "password1"}
	add, err := ser.Add(context.Background(), &proto.AddReq{
		Block: []byte("aaaa"),
		User:  user1,
	})
	if err != nil {
		return
	}
	fmt.Println(add.Cid)
}
