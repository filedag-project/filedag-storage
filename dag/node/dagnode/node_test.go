package dagnode

import (
	"bytes"
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
	q := &DataNodeClient{
		Client: newDatanode(t),
	}
	var s []*DataNodeClient
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
	content := "123456"
	block := blocks.NewBlock([]byte(content))
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
	if !bytes.Equal(block.RawData(), get.RawData()) {
		t.Fatal("the block from dagnode is not equal the origin block")
	}

	size, err := d.GetSize(ctx, block.Cid())
	if err != nil {
		fmt.Println("size err", err)
		return
	}
	if size != len(content) {
		t.Fatal("the size of block from dagnode is not equal the origin block size")
	}

	err = d.DeleteBlock(ctx, block.Cid())
	if err != nil {
		fmt.Println("del err", err)
		return
	}
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
