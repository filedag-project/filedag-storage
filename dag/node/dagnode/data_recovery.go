package dagnode

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"github.com/filedag-project/filedag-storage/dag/proto"
	"github.com/ipfs/go-cid"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
)

// RepairDataNode prepare node repair
func (d *DagNode) RepairDataNode(ctx context.Context, fromNodeIndex int, repairNodeIndex int) error {
	if fromNodeIndex >= len(d.Nodes) {
		return errors.New("index greater than max index of nodes")
	}
	if repairNodeIndex >= len(d.Nodes) {
		return errors.New("repair index greater than max index of nodes")
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	stream, err := d.Nodes[fromNodeIndex].Client.DataClient.AllKeysChan(ctx, &emptypb.Empty{})
	if err != nil {
		return err
	}
	repairNode := d.Nodes[repairNodeIndex]
	for {
		resp, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		key := resp.Key

		if _, err := repairNode.Client.DataClient.GetMeta(ctx, &proto.GetMetaRequest{Key: key}); err == nil {
			continue
		}
		dataCid, err := cid.Decode(key)
		if err != nil {
			log.Errorw("decode cid error", "key", key, "error", err)
			continue
		}
		size, err := d.GetSize(ctx, dataCid)
		if err != nil {
			log.Errorw("get block size error", "key", key, "error", err)
			continue
		}

		merged := make([][]byte, len(d.Nodes))
		for i, node := range d.Nodes {
			if i == repairNodeIndex {
				merged[i] = nil
				continue
			}
			res, err := node.Client.DataClient.Get(ctx, &proto.GetRequest{Key: key})
			if err != nil {
				log.Errorf("this node[%s] err: %v", node.RpcAddress, err)
				merged[i] = nil
				continue
			}
			if len(res.Data) == 0 {
				log.Errorf("There is no data in this node")
				merged[i] = nil
				continue
			}
			merged[i] = res.Data
		}
		enc, err := NewErasure(d.config.DataBlocks, d.config.ParityBlocks, int64(size))
		if err != nil {
			log.Errorf("new erasure fail :%v", err)
			return err
		}
		err = enc.DecodeDataBlocks(merged)
		if err != nil {
			log.Errorf("decode data blocks failed: %v", err)
			return err
		}

		meta := Meta{
			BlockSize: int32(size),
		}
		var metaBuf bytes.Buffer
		if err = binary.Write(&metaBuf, binary.LittleEndian, meta); err != nil {
			log.Errorf("binary.Write failed: %v", err)
			continue
		}
		if _, err = repairNode.Client.DataClient.Put(ctx, &proto.AddRequest{
			Key:  key,
			Meta: metaBuf.Bytes(),
			Data: merged[repairNodeIndex],
		}); err != nil {
			log.Errorf("data node put fail :%v", err)
			return err
		}
	}
}
