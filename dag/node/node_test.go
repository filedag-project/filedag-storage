package node

import (
	"encoding/json"
	"fmt"
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
	dagNode, err := NewDagNode(nc)
	data, err := ioutil.ReadFile("./node.go")
	aa := cid.Cid{}
	b, err := blocks.NewBlockWithCid(data, aa)
	if err == blocks.ErrWrongHash {
		fmt.Println(err)
	}
	err = dagNode.Put(b)
	if err != nil {
		fmt.Println(err)
	}
	dataBlock, err := dagNode.Get(aa)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(dataBlock)
}
