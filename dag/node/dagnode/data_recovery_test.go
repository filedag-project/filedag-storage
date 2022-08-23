package dagnode

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/config"
	"io/ioutil"
	"testing"
)

func TestRecovery_host(t *testing.T) {
	var nc config.DagNodeConfig
	file, err := ioutil.ReadFile("../../../conf/node_config.json")
	if err != nil {

	}
	err = json.Unmarshal(file, &nc)
	dagNode, err := NewDagNode(nc)
	err = dagNode.RepairDataNode(context.TODO(), 1, 2)
	if err != nil {
		fmt.Println(err)
	}
}
