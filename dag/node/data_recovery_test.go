package node

import (
	"encoding/json"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/config"
	"io/ioutil"
	"testing"
)

func TestRecovery_host(t *testing.T) {
	var nc config.DagNodeConfig
	file, err := ioutil.ReadFile("./node_config2.json")
	if err != nil {

	}
	err = json.Unmarshal(file, &nc)
	dagNode, err := NewDagNode(nc)
	//err = dagNode.RepairDisk("127.0.0.1", "127.0.0.1", "9010", "9013")
	//if err != nil {
	//	fmt.Println(err)
	//}
	err = dagNode.RepairDisk("127.0.0.1", "9013")
	if err != nil {
		fmt.Println(err)
	}
}
