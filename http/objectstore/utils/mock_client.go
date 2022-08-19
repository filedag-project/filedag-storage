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
	m.EXPECT().Put(gomock.Any(), gomock.AssignableToTypeOf(node)).Return(nil).AnyTimes()
	m.EXPECT().Get(gomock.Any(), gomock.AssignableToTypeOf(cid)).Return(node, nil).AnyTimes()
	m.EXPECT().GetSize(gomock.Any(), gomock.AssignableToTypeOf(cid)).Return(7, nil).AnyTimes()
	m.EXPECT().Has(gomock.Any(), gomock.AssignableToTypeOf(cid)).Return(true, nil).AnyTimes()
	m.EXPECT().DeleteBlock(gomock.Any(), gomock.AssignableToTypeOf(cid)).Return(nil).AnyTimes()
	return m, ctrl.Finish
}
