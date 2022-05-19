package pool

import (
	"github.com/filedag-project/filedag-storage/dag/node"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	"github.com/filedag-project/filedag-storage/http/objectstore/utils"
	"github.com/filedag-project/filedag-storage/kv/mutcask"
	logging "github.com/ipfs/go-log/v2"
	"strconv"
	"testing"
	"time"
)

func TestHeart_beating(t *testing.T) {
	logging.SetLogLevel("*", "DEBUG")
	db, err := uleveldb.OpenDb(utils.TmpDirPath(t))
	if err != nil {
		log.Errorf("err %v", err)
	}
	r := NewRecordSys(db)
	var a []node.DataNode
	for i := 0; i < 3; i++ {
		conn, h, err := node.InitSliceConn(":9010")
		if err != nil {
			return
		}
		a = append(a, node.DataNode{
			Client:      conn,
			HeartClient: h,
			Ip:          "127.0.0.1",
			Port:        strconv.Itoa(9010 + i),
		})
	}
	err = r.HandleDagNode(a, "test")
	if err != nil {
		return
	}
	time.Sleep(time.Second * 30)
	log.Infof("the node : %+v", r.RN)
}
func Test_MutServer(t *testing.T) {
	logging.SetLogLevel("*", "INFO")
	go mutcask.MutServer("127.0.0.1", "9010", utils.TmpDirPath(t))
	go mutcask.MutServer("127.0.0.1", "9011", utils.TmpDirPath(t))
	go mutcask.MutServer("127.0.0.1", "9012", utils.TmpDirPath(t))
	time.Sleep(time.Minute * 10)
}