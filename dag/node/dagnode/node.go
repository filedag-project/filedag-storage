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
	"github.com/filedag-project/filedag-storage/dag/slotsmgr"
	"github.com/filedag-project/filedag-storage/dag/utils/paralleltask"
	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
	"sync"
	"time"
)

var _ blockstore.Blockstore = (*DagNode)(nil)

const healthCheckService = "grpc.health.v1.Health"

type StorageNode struct {
	*datanode.Client
	State bool // true: means the data node is health
}

//DagNode Implemented the Blockstore interface
type DagNode struct {
	Nodes    []*StorageNode
	slots    *slotsmgr.SlotsManager
	numSlots int
	config   config.DagNodeConfig
	stopCh   chan struct{}
}

type Meta struct {
	BlockSize int32
}

//NewDagNode creates a new DagNode
func NewDagNode(cfg config.DagNodeConfig) (*DagNode, error) {
	numNodes := len(cfg.Nodes)
	if numNodes != cfg.DataBlocks+cfg.ParityBlocks || numNodes == 0 {
		return nil, errors.New("dag node config is incorrect")
	}
	clients := make([]*StorageNode, 0, cfg.DataBlocks+cfg.ParityBlocks)
	for _, c := range cfg.Nodes {
		dateNode, err := datanode.NewClient(c)
		if err != nil {
			return nil, err
		}
		clients = append(clients, &StorageNode{Client: dateNode})
	}
	return &DagNode{
		Nodes:  clients,
		slots:  slotsmgr.NewSlotsManager(),
		config: cfg,
		stopCh: make(chan struct{}),
	}, nil
}

func (d *DagNode) GetConfig() *config.DagNodeConfig {
	return &d.config
}

func (d *DagNode) GetDataNodeState(setIndex int) bool {
	if setIndex < 0 || setIndex >= len(d.Nodes) {
		log.Fatalf("input setIndex %v is illegal, size of set is %v", setIndex, len(d.Nodes))
	}
	return d.Nodes[setIndex].State
}

// AddSlot Set the slot bit and return the old value
func (d *DagNode) AddSlot(slot uint64) bool {
	old, err := d.slots.Set(slot, true)
	if err != nil {
		log.Fatal(err)
	}
	if !old {
		d.numSlots++
	}
	return old
}

// ClearSlot Clear the slot bit and return the old value
func (d *DagNode) ClearSlot(slot uint64) bool {
	old, err := d.slots.Set(slot, false)
	if err != nil {
		log.Fatal(err)
	}
	if old {
		d.numSlots--
	}
	return old
}

// GetSlot Get the slot bit
func (d *DagNode) GetSlot(slot uint64) bool {
	val, err := d.slots.Get(slot)
	if err != nil {
		log.Fatal(err)
	}
	return val
}

func (d *DagNode) GetSlotPairs() []slotsmgr.SlotPair {
	return d.slots.ToSlotPair()
}

func (d *DagNode) GetNumSlots() int {
	return d.numSlots
}

func (d *DagNode) RunHeartbeatCheck(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	healthCheckAll := func() {
		wg := sync.WaitGroup{}
		for _, node := range d.Nodes {
			wg.Add(1)
			go func(sn *StorageNode) {
				defer wg.Done()
				sn.State = d.healthCheck(ctx, sn.Client)
			}(node)
		}
		wg.Wait()
	}
	healthCheckAll()
	for {
		select {
		case <-ctx.Done():
			return
		case <-d.stopCh:
			cancel()
			return
		case <-ticker.C:
			healthCheckAll()
		}
	}
}

func (d *DagNode) healthCheck(ctx context.Context, cli *datanode.Client) bool {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	check, err := cli.HeartClient.Check(ctx, &healthpb.HealthCheckRequest{Service: healthCheckService})
	if err != nil {
		log.Errorf("Check the rpc address:%v err:%v", cli.RpcAddress, err)
		return false
	}
	if check.Status != healthpb.HealthCheckResponse_SERVING {
		log.Errorf("the rpc server[%v] status: %v", cli.RpcAddress, check.Status)
		return false
	}
	return true
}

//DeleteBlock deletes a block from the DagNode
func (d *DagNode) DeleteBlock(ctx context.Context, cid cid.Cid) (err error) {
	log.Warnf("delete block, cid: %v", cid)
	keyCode := cid.String()
	_, entryWriteQuorum := d.entryQuorum()
	taskCtx := context.Background()
	task := paralleltask.NewParallelTask(taskCtx, entryWriteQuorum, len(d.Nodes)-entryWriteQuorum+1, false)
	for _, snode := range d.Nodes {
		node := snode.Client
		task.Goroutine(func(ctx context.Context) error {
			var err error
			if _, err = node.Client.Delete(ctx, &proto.DeleteRequest{Key: keyCode}); err != nil {
				log.Errorw("delete error", "datanode", node.RpcAddress, "key", keyCode, "error", err)
			}
			return err
		})
	}
	return task.Wait()
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
	meta, _, onlineNodes, err := d.getMetaInfo(ctx, cid)
	if err != nil {
		return nil, err
	}
	size := meta.BlockSize
	entryReadQuorum, _ := d.entryQuorum()

	shards := make([][]byte, len(onlineNodes))
	task := paralleltask.NewParallelTask(ctx, entryReadQuorum, len(onlineNodes)-entryReadQuorum+1, true)
	for i, snode := range onlineNodes {
		index := i
		tnode := snode
		task.Goroutine(func(ctx context.Context) error {
			// is offline node?
			if tnode == nil {
				return errors.New("offline node")
			}
			node := tnode.Client
			var err error
			res, err := node.Client.Get(ctx, &proto.GetRequest{Key: keyCode})
			if err != nil {
				log.Errorw("get error", "datanode", node.RpcAddress, "key", keyCode, "error", err)
			} else {
				shards[index] = res.Data
				return nil
			}
			return err
		})

	}
	if err = task.Wait(); err != nil {
		log.Errorf("task error: %v", err)
		return nil, err
	}

	enc, err := NewErasure(d.config.DataBlocks, d.config.ParityBlocks, int64(size))
	if err != nil {
		log.Errorf("new erasure fail :%v", err)
		return nil, err
	}
	err = enc.DecodeDataBlocks(shards)
	if err != nil {
		log.Errorf("decode date blocks fail :%v", err)
		return nil, err
	}

	// merge to block raw data
	shardSize := int(enc.ShardSize())
	data := make([]byte, d.config.DataBlocks*shardSize)
	for i, shard := range shards {
		if i == d.config.DataBlocks {
			break
		}
		copy(data[i*shardSize:], shard)
	}
	data = data[:size]

	b, err := blocks.NewBlockWithCid(data, cid)
	if err == blocks.ErrWrongHash {
		return nil, blockstore.ErrHashMismatch
	}
	return b, err
}

//GetSize returns the size of the block with the given cid
func (d *DagNode) GetSize(ctx context.Context, cid cid.Cid) (int, error) {
	meta, _, _, err := d.getMetaInfo(ctx, cid)
	return int(meta.BlockSize), err
}

func (d *DagNode) getMetaInfo(ctx context.Context, cid cid.Cid) (meta Meta, metas []Meta, onlineNodes []*StorageNode, err error) {
	var errs []error
	metas, errs = readAllMeta(ctx, d.Nodes, cid.String())
	entryReadQuorum, _ := d.entryQuorum()
	reducedErr := reduceQuorumErrs(ctx, errs, entryOpIgnoredErrs, entryReadQuorum, errErasureReadQuorum)
	if reducedErr != nil {
		return meta, nil, nil, reducedErr
	}
	meta, err = findMetaInQuorum(ctx, metas, entryReadQuorum)
	if err != nil {
		return meta, nil, nil, err
	}
	onlineNodes = make([]*StorageNode, len(metas))
	for i, m := range metas {
		if m.BlockSize == meta.BlockSize {
			onlineNodes[i] = d.Nodes[i]
		} else {
			onlineNodes[i] = nil
		}
	}
	return meta, metas, onlineNodes, err
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

	enc, err := NewErasure(d.config.DataBlocks, d.config.ParityBlocks, int64(blockDataSize))
	if err != nil {
		log.Errorf("newErasure fail :%v", err)
		return err
	}
	shards, err := enc.EncodeData(blockData)
	if err != nil {
		log.Errorf("encodeData fail :%v", err)
		return err
	}

	_, entryWriteQuorum := d.entryQuorum()
	taskCtx := context.Background()
	task := paralleltask.NewParallelTask(taskCtx, entryWriteQuorum, len(d.Nodes)-entryWriteQuorum+1, false)
	for i, snode := range d.Nodes {
		index := i
		node := snode.Client
		task.Goroutine(func(ctx context.Context) error {
			var err error
			if _, err = node.Client.Put(ctx, &proto.AddRequest{
				Key:  keyCode,
				Meta: metaBuf.Bytes(),
				Data: shards[index],
			}); err != nil {
				log.Errorw("put error", "datanode", node.RpcAddress, "key", keyCode, "error", err)
			}
			return err
		})
	}
	// If the specified number of successes is met, the write succeeds,
	// or if the specified number of failures is met, the write fails
	return task.Wait()
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
	close(d.stopCh)
}

// Returns per entry readQuorum and writeQuorum
// readQuorum is the min required nodes to read data.
// writeQuorum is the min required nodes to write data.
func (d *DagNode) entryQuorum() (entryReadQuorum, entryWriteQuorum int) {
	writeQuorum := d.config.DataBlocks
	if d.config.DataBlocks == d.config.ParityBlocks {
		writeQuorum++
	}

	return d.config.DataBlocks, writeQuorum
}

// Reads all metadata as a Meta slice.
// Returns error slice indicating the failed metadata reads.
func readAllMeta(ctx context.Context, nodes []*StorageNode, key string) ([]Meta, []error) {
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
			resp, err := nodes[index].Client.Client.GetMeta(ctx, &proto.GetMetaRequest{Key: key})
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
