package proto

import (
	"context"
	"flag"
	"fmt"
	"github.com/filedag-project/filedag-storage/proto"
	"google.golang.org/grpc"
	"log"
	"testing"
)

func TestMutCask_Conn(t *testing.T) {
	// coon
	addr := flag.String("addr", "localhost:9001", "the address to connect to")
	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	// 实例化client
	c := proto.NewMutCaskClient(conn)
	bytes := []byte("123456")
	ctx := context.Background()
	res, err := c.Put(ctx, &proto.AddRequest{Key: "key", DataBlock: bytes})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}
