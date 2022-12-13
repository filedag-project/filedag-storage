package dnm

import (
	"github.com/filedag-project/filedag-storage/dag/config"
	"github.com/filedag-project/filedag-storage/dag/node/datanode"
	"github.com/filedag-project/filedag-storage/objectservice/objmetadb"
	logging "github.com/ipfs/go-log/v2"
	"testing"
	"time"
)

func TestHeart_beating(t *testing.T) {
	logging.SetLogLevel("*", "DEBUG")
	db, err := objmetadb.OpenDb(t.TempDir())
	if err != nil {
		log.Errorf("err %v", err)
	}
	r := NewRecordSys(db)
	go datanode.StartDataNodeServer(":9010", datanode.KVBadge, t.TempDir())
	time.Sleep(time.Second)
	var a []*datanode.Client

	datanodeClient, err := datanode.NewClient(config.DataNodeConfig{
		Ip:   "",
		Port: ":9010",
	})
	if err != nil {
		return
	}
	a = append(a, datanodeClient)

	err = r.HandleDagNode(a, "test")
	if err != nil {
		return
	}
	time.Sleep(time.Second * 10)
	log.Debugf("the node : %+v", r.RN)
}
