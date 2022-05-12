package dagpoolclient

import (
	"bytes"
	"context"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/pool/server"
	"github.com/filedag-project/filedag-storage/dag/pool/utils"
	"github.com/ipfs/go-merkledag"
	"google.golang.org/grpc"
	"testing"
)

func TestServer_Add(t *testing.T) {
	cli()
}
func TestPoolClient_Add(t *testing.T) {
	r := bytes.NewReader([]byte("123456"))
	cidBuilder, err := merkledag.PrefixForCidVersion(0)

	// 建立连接
	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	// 实例化client
	c := server.NewDagPoolClient(conn)

	node, err := utils.BalanceNode(context.TODO(), r, PoolClient{c}, cidBuilder)
	if err != nil {
		return
	}
	fmt.Println(string(node.RawData()))
}
