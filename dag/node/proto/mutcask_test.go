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
	// conn
	addr1 := flag.String("addr1", "localhost:9001", "the address to connect to")
	conn1, err := grpc.Dial(*addr1, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn1.Close()
	// 实例化client
	c1 := proto.NewMutCaskClient(conn1)

	addr2 := flag.String("addr2", "localhost:9002", "the address to connect to")
	conn2, err := grpc.Dial(*addr2, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn1.Close()
	// 实例化client
	c2 := proto.NewMutCaskClient(conn2)

	addr3 := flag.String("addr3", "localhost:9003", "the address to connect to")
	conn3, err := grpc.Dial(*addr3, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn3.Close()
	// 实例化client
	c3 := proto.NewMutCaskClient(conn1)
	clients := make([]proto.MutCaskClient, 0)
	clients = append(clients, c1, c2, c3)
	for _, client := range clients {
		bytes := []byte("123456")
		ctx := context.Background()
		_, err := client.Put(ctx, &proto.AddRequest{Key: "key", DataBlock: bytes})
		if err != nil {
			fmt.Println(err)
		}
	}
	//bytes := []byte("123456")
	//ctx := context.Background()
	//res, err := c.Put(ctx, &proto.AddRequest{Key: "key", DataBlock: bytes})
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(res)
}
