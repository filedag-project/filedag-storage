package node

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/config"
	"github.com/filedag-project/filedag-storage/dag/proto"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	"github.com/filedag-project/filedag-storage/kv"
	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"strings"
	"sync"
)

const lockFileName = "repo.lock"

var _ blockstore.Blockstore = (*DagNode)(nil)

type DagNode struct {
	Nodes                    []DataNode
	db                       *uleveldb.ULevelDB
	dataBlocks, parityBlocks int
}

type DataNode struct {
	sync.Mutex
	Client      proto.DataNodeClient
	HeartClient healthpb.HealthClient
	Ip          string
	Port        string
}

func NewDagNode(cfg config.NodeConfig) (*DagNode, error) {
	var s []DataNode
	for _, c := range cfg.Nodes {
		dateNode := new(DataNode)
		sc, hc, err := InitSliceConn(c.Ip + ":" + c.Port)
		if err != nil {
			return nil, err
		}
		dateNode.Ip = c.Ip
		dateNode.Port = c.Port
		dateNode.Client = sc
		dateNode.HeartClient = hc
		s = append(s, *dateNode)
	}
	db, _ := uleveldb.OpenDb(cfg.LevelDbPath)
	return &DagNode{s, db, cfg.DataBlocks, cfg.ParityBlocks}, nil
}

func InitSliceConn(addr string) (c proto.DataNodeClient, h healthpb.HealthClient, err error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		conn.Close()
		log.Errorf("did not connect: %v", err)
		return c, h, err
	}
	//defer conn.Close()
	// init client
	h = healthpb.NewHealthClient(conn)
	c = proto.NewDataNodeClient(conn)
	return c, h, nil
}

func (d DagNode) GetIP() []string {
	var s []string
	for _, n := range d.Nodes {
		s = append(s, n.Ip)
	}
	return s
}
func (d DagNode) DeleteBlock(cid cid.Cid) (err error) {
	log.Infof("delete block, cid :%v", cid)
	ctx := context.TODO()
	keyCode := sha256String(cid.String())
	wg := sync.WaitGroup{}
	wg.Add(len(d.Nodes))
	for _, node := range d.Nodes {
		go func(node DataNode) {
			defer func() {
				if err := recover(); err != nil {
					log.Errorf("%s:%s, keyCode:%s, delete block err :%v", node.Ip, node.Port, keyCode, err)
				}
				wg.Done()
			}()
			_, err = node.Client.Delete(ctx, &proto.DeleteRequest{Key: keyCode})
			if err != nil {
				log.Errorf("%s:%s, keyCode:%s, delete block err :%v", node.Ip, node.Port, keyCode, err)
			}
		}(node)
	}
	wg.Wait()
	return err
}

func (d DagNode) Has(cid cid.Cid) (bool, error) {
	_, err := d.GetSize(cid)
	if err != nil {
		if strings.Contains(err.Error(), kv.ErrNotFound.Error()) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (d DagNode) Get(cid cid.Cid) (blocks.Block, error) {
	log.Infof("get block, cid :%v", cid)
	ctx := context.TODO()
	keyCode := sha256String(cid.String())
	var err error
	var size int
	err = d.db.Get(cid.String(), &size)
	if err != nil {
		return nil, err
	}
	merged := make([][]byte, len(d.Nodes))
	wg := sync.WaitGroup{}
	wg.Add(len(d.Nodes))
	for i, node := range d.Nodes {
		go func(i int, node DataNode) {
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
				merged[i] = res.DataBlock
			}
		}(i, node)
	}
	wg.Wait()
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
	var data []byte
	data = bytes.Join(merged, []byte(""))
	if err != nil {
		return nil, err
	}
	data = data[:size]
	b, err := blocks.NewBlockWithCid(data, cid)
	if err == blocks.ErrWrongHash {
		return nil, blockstore.ErrHashMismatch
	}
	return b, err
}

func (d DagNode) GetSize(cid cid.Cid) (int, error) {
	ctx := context.TODO()
	keyCode := sha256String(cid.String())
	var err error
	var count int64
	for _, node := range d.Nodes {
		size, err := node.Client.Size(ctx, &proto.SizeRequest{
			Key: keyCode,
		})
		if err != nil {
			return 0, err
		}
		count = count + size.Size
	}
	return int(count), err
}

func (d DagNode) Put(block blocks.Block) (err error) {
	log.Infof("put block, cid :%v", block.Cid())
	ctx := context.TODO()
	//todo store this info in datanode
	err = d.db.Put(block.Cid().String(), len(block.RawData()))
	if err != nil {
		return err
	}
	keyCode := sha256String(block.Cid().String())
	enc, err := NewErasure(d.dataBlocks, d.parityBlocks, int64(len(block.RawData())))
	if err != nil {
		log.Errorf("newErasure fail :%v", err)
		return err
	}
	shards, err := enc.EncodeData(block.RawData())
	if err != nil {
		log.Errorf("encodeData fail :%v", err)
		return err
	}
	ok, err := enc.encoder().Verify(shards)
	if err != nil {
		log.Errorf("encode fail :%v", err)
		return err
	}
	if ok && err == nil {
		log.Debugf("encode ok, the data is the same format as Encode. No data is modified")
	}
	wg := sync.WaitGroup{}
	wg.Add(len(d.Nodes))
	for i, node := range d.Nodes {
		go func(i int, node DataNode) {
			defer func() {
				if err := recover(); err != nil {
					log.Errorf("%s:%s,keyCode:%s,kvdb put :%v", node.Ip, node.Port, keyCode, err)
				}
				wg.Done()
			}()
			_, err = node.Client.Put(ctx, &proto.AddRequest{Key: keyCode, DataBlock: shards[i]})
			if err != nil {
				log.Errorf("%s:%s,keyCode:%s,kvdb put :%v", node.Ip, node.Port, keyCode, err)
			}
		}(i, node)
	}
	wg.Wait()
	return err
}

func (d DagNode) PutMany(blocks []blocks.Block) (err error) {
	for _, block := range blocks {
		err = d.Put(block)
	}
	return err
}

func (d DagNode) AllKeysChan(ctx context.Context) (<-chan cid.Cid, error) {
	panic("implement me")
}

func (d DagNode) HashOnRead(enabled bool) {
	panic("implement me")
}

func sha256String(s string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
}
