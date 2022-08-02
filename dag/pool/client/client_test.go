package client

import (
	"bytes"
	"context"
	"fmt"
	"github.com/ipfs/go-blockservice"
	offline "github.com/ipfs/go-ipfs-exchange-offline"
	logging "github.com/ipfs/go-log/v2"
	"github.com/ipfs/go-merkledag"
	"testing"
	"time"
)

func TestPoolClient_Add_Get(t *testing.T) {
	time.Sleep(time.Second * 1)
	logging.SetAllLoggers(logging.LevelDebug)
	r := bytes.NewReader([]byte("123456"))
	cidBuilder, err := merkledag.PrefixForCidVersion(0)
	poolCli, cancel := NewMockPoolClient(t)
	defer cancel()
	var ctx = context.Background()
	dagServ := merkledag.NewDAGService(blockservice.New(poolCli, offline.Exchange(poolCli)))
	node, err := BalanceNode(r, dagServ, cidBuilder)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	fmt.Println("aaaaa", node.Cid().String())
	get, err := poolCli.Get(ctx, node.Cid())
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	fmt.Println(get.String())
}
