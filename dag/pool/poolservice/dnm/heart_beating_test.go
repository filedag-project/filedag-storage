package dnm

import (
	"github.com/filedag-project/filedag-storage/dag/node/dagnode"
	"github.com/filedag-project/filedag-storage/dag/node/datanode"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	"github.com/filedag-project/filedag-storage/http/objectstore/utils"
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
	go datanode.MutDataNodeServer(":9010", datanode.KVBadge, utils.TmpDirPath(t))
	time.Sleep(time.Second)
	var a []dagnode.DataNode
	for i := 0; i < 3; i++ {
		conn, h, err := dagnode.InitSliceConn(":9010")
		if err != nil {
			return
		}
		a = append(a, dagnode.DataNode{
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
	time.Sleep(time.Second * 10)
	log.Infof("the node : %+v", r.RN)
}
