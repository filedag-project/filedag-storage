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
	//server.StartTestServer(t)
	r := bytes.NewReader([]byte("123456"))
	cidBuilder, err := merkledag.PrefixForCidVersion(0)

	addr := flag.String("addr", "localhost:9002", "the address to connect to")
	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := server.NewDagPoolClient(conn)
	pc := PoolClient{c, cidBuilder, conn}
	var ctx = context.Background()
	ctx = context.WithValue(ctx, "user", "test,test123")
	node, err := BalanceNode(ctx, r, pc, cidBuilder)
	if err != nil {
		return
	}
	fmt.Println("aaaaa", node.Cid().String())
	get, err := pc.Get(ctx, node.Cid())
	if err != nil {
		return
	}
	fmt.Println(get.String())
}
