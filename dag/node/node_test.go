package node

import (
	"encoding/json"
	"github.com/filedag-project/filedag-storage/dag/config"
	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	"io/ioutil"
	"testing"
)

func TestDagNode_put(t *testing.T) {
	var nc config.NodeConfig
	file, err := ioutil.ReadFile("./node_config2.json")
	if err != nil {

	}
	err = json.Unmarshal(file, &nc)
	//for i, node := range nc.Nodes {
	//	name := "addr" +fmt.Sprint(i)
	//	addr := flag.String(name, fmt.Sprintf("%s:%s",node.Ip, node.Port), "the address to connect to")
	//	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	//	fmt.Println(conn,err)
	//}

	dagNode, err := NewDagNode(nc)
	data, err := ioutil.ReadFile("./node.go")
	aa := cid.Cid{}
	b, err := blocks.NewBlockWithCid(data, aa)
	//if err == blocks.ErrWrongHash {
	//	return nil, blockstore.ErrHashMismatch
	//}
	dagNode.Put(b)
}
