package dagnode

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"github.com/filedag-project/filedag-storage/dag/proto"
	"github.com/filedag-project/filedag-storage/dag/utils/paralleltask"
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

		shards := make([][]byte, len(d.Nodes))
		entryReadQuorum, _ := d.entryQuorum()
		task := paralleltask.NewParallelTask(ctx, entryReadQuorum, len(d.Nodes)-entryReadQuorum+1, true)
		for i, snode := range d.Nodes {
			index := i
			tnode := snode
			task.Goroutine(func(ctx context.Context) error {
				if index == repairNodeIndex {
					return errors.New("there is no data in this node")
				}
				res, err := tnode.Client.DataClient.Get(ctx, &proto.GetRequest{Key: key})
				if err != nil {
					log.Errorf("this node[%s] get key err: %v", tnode.RpcAddress, err)
					return err
				}
				if len(res.Data) == 0 {
					err = errors.New("there is no data in this node")
					return err
				}
				shards[index] = res.Data
				return nil
			})
		}
		if err = task.Wait(); err != nil {
			log.Errorw("task error, missing shards", "key", key, "error", err)
			continue
		}

		enc, err := NewErasure(d.config.DataBlocks, d.config.ParityBlocks, int64(size))
		if err != nil {
			log.Errorf("new erasure fail :%v", err)
			return err
		}
		err = enc.DecodeDataAndParityBlocks(shards)
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
			Data: shards[repairNodeIndex],
		}); err != nil {
			log.Errorf("data node put failed: %v", err)
			return err
		}
		log.Infow("repair entry success", "key", key)
	}
}
