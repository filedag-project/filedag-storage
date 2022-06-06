package node

import (
	"context"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/node/mocks"
	"github.com/filedag-project/filedag-storage/dag/proto"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	"github.com/filedag-project/filedag-storage/http/objectstore/utils"
	"github.com/golang/mock/gomock"
	blocks "github.com/ipfs/go-block-format"
	"testing"
)

func TestDagNode(t *testing.T) {
	q := DataNode{
		Client: newDatanode(t),
	}
	var s []DataNode
	s = append(s, q)
	db, err := uleveldb.OpenDb(utils.TmpDirPath(t))
	if err != nil {
		return
	}
	var d = DagNode{
		Nodes:        s,
		db:           db,
		dataBlocks:   3,
		parityBlocks: 1,
	}
	block := blocks.NewBlock([]byte("123456"))
	ctx := context.TODO()
	err = d.Put(ctx, block)
	if err != nil {
		fmt.Println("put err", err)
		return
	}
	get, err := d.Get(ctx, block.Cid())
	if err != nil {
		fmt.Println("get err", err)
		return
	}
	fmt.Println(get.String())
	err = d.DeleteBlock(ctx, block.Cid())
	if err != nil {
		fmt.Println("del err", err)
		return
	}
	size, err := d.GetSize(ctx, block.Cid())
	if err != nil {
		fmt.Println("size err", err)
		return
	}
	fmt.Println(size)
}
func newDatanode(t *testing.T) *mocks.MockDataNodeClient {
	ctrl := gomock.NewController(t)
	m := mocks.NewMockDataNodeClient(ctrl)
	m.EXPECT().Put(gomock.AssignableToTypeOf(context.Background()), gomock.AssignableToTypeOf(&proto.AddRequest{})).AnyTimes().Return(nil, nil)
	m.EXPECT().Get(gomock.AssignableToTypeOf(context.Background()), gomock.AssignableToTypeOf(&proto.GetRequest{})).AnyTimes().
		Return(&proto.GetResponse{DataBlock: []byte("123456")}, nil)
	m.EXPECT().Delete(gomock.AssignableToTypeOf(context.Background()), gomock.AssignableToTypeOf(&proto.DeleteRequest{})).AnyTimes().Return(nil, nil)
	m.EXPECT().Size(gomock.AssignableToTypeOf(context.Background()), gomock.AssignableToTypeOf(&proto.SizeRequest{})).AnyTimes().
		Return(&proto.SizeResponse{Size: 6}, nil)
	return m
}
