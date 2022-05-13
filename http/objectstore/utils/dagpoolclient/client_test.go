package dagpoolclient

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/pool/server"
	"github.com/ipfs/go-merkledag"
	"google.golang.org/grpc"
	"testing"
)

func TestPoolClient_Add(t *testing.T) {
	r := bytes.NewReader([]byte("123456"))
	cidBuilder, err := merkledag.PrefixForCidVersion(0)

	// 建立连接
	addr := flag.String("addr", "localhost:50051", "the address to connect to")
	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	// 实例化client
	c := server.NewDagPoolClient(conn)
	pc := PoolClient{c, cidBuilder}
	var ctx = context.Background()
	ctx = context.WithValue(ctx, "user", "test,test123")
	node, err := BalanceNode(ctx, r, pc, cidBuilder)
	if err != nil {
		return
	}
	fmt.Println("aaaaa", string(node.RawData()))
	get, err := pc.Get(ctx, node.Cid())
	if err != nil {
		return
	}
	fmt.Println(get.String())
}
