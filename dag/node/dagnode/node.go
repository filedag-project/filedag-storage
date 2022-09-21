package dagnode

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/config"
	"github.com/filedag-project/filedag-storage/dag/node/datanode"
	"github.com/filedag-project/filedag-storage/dag/proto"
	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
)

var _ blockstore.Blockstore = (*DagNode)(nil)

//DagNode Implemented the Blockstore interface
type DagNode struct {
	Nodes        []*datanode.Client
	dataBlocks   int // Number of data shards
	parityBlocks int // Number of parity shards
}

type Meta struct {
	BlockSize int32
}

//NewDagNode creates a new DagNode
func NewDagNode(cfg config.DagNodeConfig) (*DagNode, error) {
	var clients []*datanode.Client
	for _, c := range cfg.Nodes {
		dateNode, err := datanode.NewClient(c)
		if err != nil {
			return nil, err
		}
		clients = append(clients, dateNode)
	}
	return &DagNode{
		Nodes:        clients,
		dataBlocks:   cfg.DataBlocks,
		parityBlocks: cfg.ParityBlocks,
	}, nil
}

//DeleteBlock deletes a block from the DagNode
func (d *DagNode) DeleteBlock(ctx context.Context, cid cid.Cid) (err error) {
	log.Warnf("delete block, cid :%v", cid)
	keyCode := cid.String()
	wg := sync.WaitGroup{}
	wg.Add(len(d.Nodes))
	for _, node := range d.Nodes {
		go func(node *datanode.Client) {
			defer func() {
				if err := recover(); err != nil {
					log.Errorf("%s:%s, keyCode:%s, delete block err :%v", node.Ip, node.Port, keyCode, err)
				}
				wg.Done()
			}()
			_, err = node.Client.Delete(ctx, &proto.DeleteRequest{Key: keyCode})
			if err != nil {
				log.Debugf("%s:%s, keyCode:%s, delete block err :%v", node.Ip, node.Port, keyCode, err)
			}
		}(node)
	}
	wg.Wait()
	return err
}

//Has returns true if the given cid is in the DagNode
func (d *DagNode) Has(ctx context.Context, cid cid.Cid) (bool, error) {
	if _, err := d.GetSize(ctx, cid); err != nil {
		return false, err
	}

	return true, nil
}

//Get returns the block with the given cid
func (d *DagNode) Get(ctx context.Context, cid cid.Cid) (blocks.Block, error) {
	log.Debugf("get block, cid :%v", cid)
	keyCode := cid.String()
	size, err := d.GetSize(ctx, cid)
	if err != nil {
		return nil, err
	}

	merged := make([][]byte, len(d.Nodes))
	wg := sync.WaitGroup{}
	wg.Add(len(d.Nodes))
	for i, node := range d.Nodes {
		go func(i int, node *datanode.Client) {
			defer func() {
				if err := recover(); err != nil {
					log.Errorf("%s:%s, keyCode:%s, kvdb get err :%v", node.Ip, node.Port, keyCode, err)
				}
				wg.Done()
			}()
			res, err := node.Client.Get(ctx, &proto.GetRequest{Key: keyCode})
			if err != nil {
				log.Errorf("%s:%s, keyCode:%s,kvdb get :%v", node.Ip, node.Port, keyCode, err)
				merged[i] = nil
			} else {
				merged[i] = res.Data
			}
		}(i, node)
	}
	wg.Wait()
	// TODO: After obtaining the shard data that meets the conditions, we can proceed

	enc, err := NewErasure(d.dataBlocks, d.parityBlocks, int64(size))
	if err != nil {
		log.Errorf("new erasure fail :%v", err)
		return nil, err
	}
	err = enc.DecodeDataBlocks(merged)
	if err != nil {
		log.Errorf("decode date blocks fail :%v", err)
		return nil, err
	}
	data := bytes.Join(merged, []byte(""))
	data = data[:size]
	b, err := blocks.NewBlockWithCid(data, cid)
	if err == blocks.ErrWrongHash {
		return nil, blockstore.ErrHashMismatch
	}
	return b, err
}

//GetSize returns the size of the block with the given cid
func (d *DagNode) GetSize(ctx context.Context, cid cid.Cid) (int, error) {
	metas, errs := readAllMeta(ctx, d.Nodes, cid.String())
	entryReadQuorum, _ := d.entryQuorum()
	reducedErr := reduceQuorumErrs(ctx, errs, entryOpIgnoredErrs, len(metas)/2, errErasureReadQuorum)
	if reducedErr != nil {
		return 0, reducedErr
	}
	meta, err := findMetaInQuorum(ctx, metas, entryReadQuorum)
	if err != nil {
		return 0, err
	}
	return int(meta.BlockSize), err
}

//Put adds the given block to the DagNode
func (d *DagNode) Put(ctx context.Context, block blocks.Block) (err error) {
	log.Debugf("put block, cid :%v", block.Cid())
	// copy data from block, because reedsolomon may modify data
	buf := bytes.NewBuffer(nil)
	buf.Write(block.RawData())
	blockData := buf.Bytes()
	blockDataSize := len(blockData)
	keyCode := block.Cid().String()

	meta := Meta{
		BlockSize: int32(blockDataSize),
	}
	var metaBuf bytes.Buffer
	err = binary.Write(&metaBuf, binary.LittleEndian, meta)
	if err != nil {
		return err
	}

	enc, err := NewErasure(d.dataBlocks, d.parityBlocks, int64(blockDataSize))
	if err != nil {
		log.Errorf("newErasure fail :%v", err)
		return err
	}
	shards, err := enc.EncodeData(blockData)
	if err != nil {
		log.Errorf("encodeData fail :%v", err)
		return err
	}
	// TODO: Is verify necessary?
	ok, err := enc.encoder().Verify(shards)
	if err != nil {
		log.Errorf("encode fail :%v", err)
		return err
	}
	if ok {
		//log.Debugf("encode ok, the data is the same format as Encode. No data is modified")
	}
	wg := sync.WaitGroup{}
	wg.Add(len(d.Nodes))
	for i, node := range d.Nodes {
		go func(i int, node *datanode.Client) {
			defer func() {
				if err := recover(); err != nil {
					log.Errorf("%s:%s,keyCode:%s,kvdb put :%v", node.Ip, node.Port, keyCode, err)
				}
				wg.Done()
			}()
			if _, err = node.Client.Put(ctx, &proto.AddRequest{
				Key:  keyCode,
				Meta: metaBuf.Bytes(),
				Data: shards[i],
			}); err != nil {
				log.Errorf("%s:%s,keyCode:%s,kvdb put :%v", node.Ip, node.Port, keyCode, err)
				// TODO: Put failure handling
			}
		}(i, node)
	}
	// TODO: If the specified number of successes is met, the write succeeds,
	// or if the specified number of failures is met, the write fails
	wg.Wait()
	return err
}

//PutMany adds the given blocks to the DagNode
func (d *DagNode) PutMany(ctx context.Context, blocks []blocks.Block) (err error) {
	for _, block := range blocks {
		err = d.Put(ctx, block)
	}
	return err
}

//AllKeysChan returns a channel that will yield every key in the dag
func (d *DagNode) AllKeysChan(ctx context.Context) (<-chan cid.Cid, error) {
	panic("implement me")
}

//HashOnRead tells the dag node to calculate the hash of the block
func (d *DagNode) HashOnRead(enabled bool) {
	panic("implement me")
}

func (d *DagNode) Close() {
	for _, nd := range d.Nodes {
		nd.Conn.Close()
	}
}

// Returns per entry readQuorum and writeQuorum
// readQuorum is the min required nodes to read data.
// writeQuorum is the min required nodes to write data.
func (d *DagNode) entryQuorum() (entryReadQuorum, entryWriteQuorum int) {
	writeQuorum := d.dataBlocks
	if d.dataBlocks == d.parityBlocks {
		writeQuorum++
	}

	return d.dataBlocks, writeQuorum
}

// Reads all metadata as a Meta slice.
// Returns error slice indicating the failed metadata reads.
func readAllMeta(ctx context.Context, nodes []*datanode.Client, key string) ([]Meta, []error) {
	metadataArray := make([]Meta, len(nodes))
	errs := make([]error, len(nodes))
	// Read meta in parallel across nodes.
	var wg sync.WaitGroup
	wg.Add(len(nodes))
	for index := range nodes {
		index := index
		go func(index int) {
			defer wg.Done()

			if nodes[index] == nil {
				errs[index] = errNodeNotFound
				return
			}
			resp, err := nodes[index].Client.GetMeta(ctx, &proto.GetMetaRequest{Key: key})
			if err != nil {
				if st, ok := status.FromError(err); ok && st.Code() == codes.Unknown {
					errs[index] = errors.New(st.Message())
				} else {
					errs[index] = err
				}
				return
			}
			meta := Meta{}
			metaBuf := bytes.NewBuffer(resp.Meta)
			err = binary.Read(metaBuf, binary.LittleEndian, &meta)
			if err != nil {
				errs[index] = err
				return
			}

			metadataArray[index] = meta
		}(index)
	}
	wg.Wait()

	// Return all the metadata.
	return metadataArray, errs
}

func findMetaInQuorum(ctx context.Context, metaArr []Meta, quorum int) (Meta, error) {
	// with less quorum return error.
	if quorum < 2 {
		return Meta{}, errErasureReadQuorum
	}
	metaHashes := make([]string, len(metaArr))
	h := sha256.New()
	for i, meta := range metaArr {
		fmt.Fprint(h, meta.BlockSize)

		metaHashes[i] = hex.EncodeToString(h.Sum(nil))
		h.Reset()
	}

	metaHashCountMap := make(map[string]int)
	for _, hash := range metaHashes {
		if hash == "" {
			continue
		}
		metaHashCountMap[hash]++
	}

	maxHash := ""
	maxCount := 0
	for hash, count := range metaHashCountMap {
		if count > maxCount {
			maxCount = count
			maxHash = hash
		}
	}

	if maxCount < quorum {
		return Meta{}, errErasureReadQuorum
	}

	for i, hash := range metaHashes {
		if hash == maxHash {
			return metaArr[i], nil
		}
	}

	return Meta{}, errErasureReadQuorum
}
