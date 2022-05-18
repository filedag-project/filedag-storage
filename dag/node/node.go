package node

import (
	"bytes"
	"context"
	"crypto/sha256"
	"flag"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/config"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	"github.com/filedag-project/filedag-storage/kv"
	"github.com/filedag-project/filedag-storage/proto"
	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	"google.golang.org/grpc"
	"strings"
	"sync"
)

const lockFileName = "repo.lock"

//var _ blockstore.Blockstore = (*DagNode)(nil)

type DagNode struct {
	nodes                    []DataNode
	db                       *uleveldb.ULevelDB
	dataBlocks, parityBlocks int
}

type DataNode struct {
	sync.Mutex
	Client proto.MutCaskClient
	Ip     string
	Port   string
}

func NewDagNode(cfg config.NodeConfig) (*DagNode, error) {
	var s []DataNode
	for i, c := range cfg.Nodes {
		dateNode := new(DataNode)
		sc, err := InitSliceConn(fmt.Sprint(i), c.Ip, c.Port)
		if err != nil {
			return nil, err
		}
		dateNode.Ip = c.Ip
		dateNode.Port = c.Port
		dateNode.Client = sc
		s = append(s, *dateNode)
	}
	db, _ := uleveldb.OpenDb(cfg.LevelDbPath)
	return &DagNode{s, db, cfg.DataBlocks, cfg.ParityBlocks}, nil
}

func InitSliceConn(index, ip, port string) (c proto.MutCaskClient, err error) {
	addr := flag.String("addr"+fmt.Sprint(index), fmt.Sprintf("%s:%s", ip, port), "the address to connect to")
	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		conn.Close()
		log.Errorf("did not connect: %v", err)
		return c, err
	}
	//defer conn.Close()
	// init client
	c = proto.NewMutCaskClient(conn)
	return c, nil
}

func (d DagNode) GetIP() []string {
	var s []string
	for _, n := range d.nodes {
		s = append(s, n.Ip)
	}
	return s
}
func (d DagNode) DeleteBlock(cid cid.Cid) (err error) {
	ctx := context.TODO()
	keyCode := sha256String(cid.String())
	for _, node := range d.nodes {
		_, err := node.Client.Delete(ctx, &proto.DeleteRequest{Key: keyCode})
		if err != nil {
			break
		}
	}
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
	ctx := context.TODO()
	keyCode := sha256String(cid.String())
	var err error
	var size int
	err = d.db.Get(cid.String(), &size)
	if err != nil {
		return nil, err
	}
	merged := make([][]byte, 0)
	for _, node := range d.nodes {
		res, err := node.Client.Get(ctx, &proto.GetRequest{Key: keyCode})
		if err != nil {
			log.Errorf("mutcask get :%v", err)
		}
		merged = append(merged, res.DataBlock)
	}
	enc, err := NewErasure(d.dataBlocks, d.parityBlocks, int64(size))
	enc.DecodeDataBlocks(merged)
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
	for _, node := range d.nodes {
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
	ctx := context.TODO()
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
		log.Infof("encode ok, the data is the same format as Encode. No data is modified")
	}
	for i, node := range d.nodes {
		_, err = node.Client.Put(ctx, &proto.AddRequest{Key: keyCode, DataBlock: shards[i]})
		if err != nil {
			break
		}
	}

	return err
}

func (d DagNode) PutMany(blocks []blocks.Block) error {
	panic("implement me")
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
