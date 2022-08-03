package dnm

import (
	"github.com/filedag-project/filedag-storage/dag/config"
	"github.com/filedag-project/filedag-storage/dag/node"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	logging "github.com/ipfs/go-log/v2"
	"testing"
	"time"
)

func TestHeart_beating(t *testing.T) {
	logging.SetAllLoggers(logging.LevelDebug)
	db, err := uleveldb.OpenDb(t.TempDir())
	if err != nil {
		log.Errorf("err %v", err)
	}
	r := NewRecordSys(db)
	go node.MutDataNodeServer(":9010", node.KVBadge, t.TempDir())
	time.Sleep(time.Second)
	var a []*node.DataNodeClient
	for i := 0; i < 3; i++ {
		datanodeClient, err := node.InitDataNodeClient(config.DataNodeConfig{
			Ip:   "",
			Port: "",
		})
		if err != nil {
			return
		}
		a = append(a, datanodeClient)
	}
	err = r.HandleDagNode(a, "test")
	if err != nil {
		return
	}
	time.Sleep(time.Second * 10)
	log.Debugf("the node : %+v", r.RN)
}
