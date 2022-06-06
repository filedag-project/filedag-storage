package utils

import (
	"github.com/filedag-project/filedag-storage/dag/pool/client"
	"github.com/filedag-project/filedag-storage/dag/pool/client/mocks"
	"github.com/golang/mock/gomock"
	"github.com/ipfs/go-merkledag"
	"testing"
)

func NewMockPoolClient(t *testing.T) (client.PoolClient, func()) {
	ctrl := gomock.NewController(t)
	m := mocks.NewMockPoolClient(ctrl)
	node := merkledag.NodeWithData([]byte("\b\u0002\u0012\a1234567\u0018\a"))
	cid := node.Cid()
	m.EXPECT().Add(gomock.Any(), gomock.AssignableToTypeOf(node)).Return(nil).AnyTimes()
	m.EXPECT().Get(gomock.Any(), gomock.AssignableToTypeOf(cid)).Return(node, nil).AnyTimes()
	m.EXPECT().Remove(gomock.Any(), gomock.AssignableToTypeOf(cid)).Return(nil).AnyTimes()
	return m, ctrl.Finish
}

func NewMockPinPoolClient(t *testing.T) (client.DataPin, func()) {
	ctrl := gomock.NewController(t)
	m := mocks.NewMockDataPin(ctrl)
	node := merkledag.NodeWithData([]byte("\b\u0002\u0012\a1234567\u0018\a"))
	cid := node.Cid()
	m.EXPECT().Pin(gomock.Any(), gomock.AssignableToTypeOf(node)).Return(nil).AnyTimes()
	m.EXPECT().UnPin(gomock.Any(), gomock.AssignableToTypeOf(cid)).Return(node, nil).AnyTimes()
	m.EXPECT().IsPin(gomock.Any(), gomock.AssignableToTypeOf(cid)).Return(nil).AnyTimes()
	return m, ctrl.Finish
}
