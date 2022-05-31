package client

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/proto"
	logging "github.com/ipfs/go-log/v2"
	"github.com/ipfs/go-merkledag"
	"google.golang.org/grpc"
	"testing"
	"time"
)

func TestPoolClient_Add_Get(t *testing.T) {
	//go server.StartTestDagPoolServer(t)
	time.Sleep(time.Second * 1)
	logging.SetLogLevel("*", "INFO")
	r := bytes.NewReader([]byte("123456"))
	cidBuilder, err := merkledag.PrefixForCidVersion(0)

	addr := flag.String("addr", "localhost:50001", "the address to connect to")
	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Errorf("did not connect: %v", err)
	}
	defer conn.Close()
	c := proto.NewDagPoolClient(conn)
	pc := PoolClient{c, cidBuilder, conn}
	var ctx = context.Background()
	node, err := BalanceNode(ctx, r, pc, cidBuilder)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	fmt.Println("aaaaa", node.Cid().String())
	ctx = context.WithValue(ctx, "user", "pool,pool123")
	get, err := pc.Get(ctx, node.Cid())
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	fmt.Println(get.String())
}
