package dagpoolclient

import (
	"context"
	"flag"
	"github.com/filedag-project/filedag-storage/dag/pool/server"
	logging "github.com/ipfs/go-log/v2"
	"google.golang.org/grpc"
	"time"
)

var log = logging.Logger("pool-client")
var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func cli() {
	flag.Parse()

	// 建立连接
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
