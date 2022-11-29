package dagnode

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/config"
	"github.com/filedag-project/filedag-storage/dag/node/datanode"
	"github.com/filedag-project/filedag-storage/dag/node/datanode/mocks"
	"github.com/filedag-project/filedag-storage/dag/proto"
	"github.com/golang/mock/gomock"
	blocks "github.com/ipfs/go-block-format"
	"testing"
)

func TestDagNode(t *testing.T) {
	var clients []*StorageNode
	for i := 0; i < 3; i++ {
		cli := &datanode.Client{
			Client: newDatanode(t),
		}
		clients = append(clients, &StorageNode{Client: cli})
	}
	var d = DagNode{
		Nodes: clients,
		config: config.DagNodeConfig{
			DataBlocks:   2,
			ParityBlocks: 1,
		},
	}
	content := "123456"
	block := blocks.NewBlock([]byte(content))
	ctx := context.TODO()
	err := d.Put(ctx, block)
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
	content := "123456"
	block := blocks.NewBlock([]byte(content))
	meta := Meta{
		BlockSize: int32(len(content)),
	}
	var metaBuf bytes.Buffer
	if err := binary.Write(&metaBuf, binary.LittleEndian, meta); err != nil {
		t.Fatalf("binary.Write failed: %v", err)
	}
	ctrl := gomock.NewController(t)
	m := mocks.NewMockDataNodeClient(ctrl)
	m.EXPECT().Put(gomock.AssignableToTypeOf(context.Background()), gomock.AssignableToTypeOf(&proto.AddRequest{})).AnyTimes().Return(nil, nil)
	m.EXPECT().Get(gomock.AssignableToTypeOf(context.Background()), gomock.AssignableToTypeOf(&proto.GetRequest{})).AnyTimes().
		Return(&proto.GetResponse{Data: block.RawData(), Meta: metaBuf.Bytes()}, nil)
	m.EXPECT().GetMeta(gomock.AssignableToTypeOf(context.Background()), gomock.AssignableToTypeOf(&proto.GetMetaRequest{})).AnyTimes().
		Return(&proto.GetMetaResponse{Meta: metaBuf.Bytes()}, nil)
	m.EXPECT().Delete(gomock.AssignableToTypeOf(context.Background()), gomock.AssignableToTypeOf(&proto.DeleteRequest{})).AnyTimes().Return(nil, nil)
	m.EXPECT().Size(gomock.AssignableToTypeOf(context.Background()), gomock.AssignableToTypeOf(&proto.SizeRequest{})).AnyTimes().
		Return(&proto.SizeResponse{Size: int64(datanode.HeaderSize + 4 + len(block.RawData()))}, nil)
	return m
}
