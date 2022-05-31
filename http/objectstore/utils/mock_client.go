package utils

import (
	dagpoolcli "github.com/filedag-project/filedag-storage/dag/pool/client"
	"github.com/golang/mock/gomock"
	"github.com/ipfs/go-merkledag"
)

func NewMockClient(ctrl *gomock.Controller) dagpoolcli.PoolClient {
	m := dagpoolcli.NewMockPoolClient(ctrl)
	m.EXPECT().Add(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	m.EXPECT().Get(gomock.Any(), gomock.Any()).Return(merkledag.NewRawNode([]byte("123456")), nil).AnyTimes()
	m.EXPECT().Close(gomock.Any()).AnyTimes()
	m.EXPECT().Remove(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	return m
}
