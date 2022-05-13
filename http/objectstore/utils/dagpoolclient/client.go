package dagpoolclient

import (
	"context"
	"flag"
	"github.com/filedag-project/filedag-storage/dag/pool/server"
	"github.com/filedag-project/filedag-storage/dag/pool/userpolicy"
	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	format "github.com/ipfs/go-ipld-format"
	legacy "github.com/ipfs/go-ipld-legacy"
	logging "github.com/ipfs/go-log/v2"
	"github.com/ipfs/go-merkledag"
	"google.golang.org/grpc"
	"strings"
	"time"
)

var log = logging.Logger("pool-client")

type PoolClient struct {
	pc         server.DagPoolClient
	CidBuilder cid.Builder
}

func NewPoolClient(addr string) (*PoolClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {

		log.Fatalf("did not connect: %v", err)
		return nil, err
	}
	defer conn.Close()
	// 实例化client
	c := server.NewDagPoolClient(conn)
	cidBuilder, err := merkledag.PrefixForCidVersion(0)
	return &PoolClient{c, cidBuilder}, nil
}

func (p PoolClient) Get(ctx context.Context, cid cid.Cid) (format.Node, error) {
	s := strings.Split((ctx.Value("user")).(string), ",")
	if len(s) != 2 {
		return nil, userpolicy.AccessDenied
	}
	get, err := p.pc.Get(ctx, &server.GetRequest{Cid: cid.String(), User: &server.PoolUser{
		Username: s[0],
		Pass:     s[1],
	}})
	if err != nil {
		return nil, err
	}
	return legacy.DecodeNode(ctx, blocks.NewBlock(get.Block))
}

func (p PoolClient) GetMany(ctx context.Context, cids []cid.Cid) <-chan *format.NodeOption {
	panic("implement me")
}

func (p PoolClient) Add(ctx context.Context, node format.Node) error {
	s := strings.Split((ctx.Value("user")).(string), ",")
	if len(s) != 2 {
		return userpolicy.AccessDenied
	}
	_, err := p.pc.Add(ctx, &server.AddRequest{Block: node.RawData(), User: &server.PoolUser{
		Username: s[0],
		Pass:     s[1],
	}})
	if err != nil {
		return err
	}
	return nil
}

func (p PoolClient) AddMany(ctx context.Context, nodes []format.Node) error {
	panic("implement me")
}

func (p PoolClient) Remove(ctx context.Context, cid cid.Cid) error {
	panic("implement me")
}

func (p PoolClient) RemoveMany(ctx context.Context, cids []cid.Cid) error {
	panic("implement me")
}

var _ format.DAGService = &PoolClient{}

func cli() {
	flag.Parse()

	// 建立连接
	addr := flag.String("addr", "localhost:50051", "the address to connect to")
	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	// 实例化client
	c := server.NewDagPoolClient(conn)

	// 调用rpc，等待同步响应
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Add(ctx, &server.AddRequest{Block: []byte("123456"), User: &server.PoolUser{
		Username: "test",
		Pass:     "test",
	}})
	if err != nil {
		log.Fatalf("could not add: %v", err)
	}
	log.Infof("add: %s", r.Cid)
}
