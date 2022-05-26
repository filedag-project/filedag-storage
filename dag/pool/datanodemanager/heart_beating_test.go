package datanodemanager

import (
	"github.com/filedag-project/filedag-storage/dag/node"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	"github.com/filedag-project/filedag-storage/http/objectstore/utils"
	logging "github.com/ipfs/go-log/v2"
	"os"
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
	os.Setenv(node.Host, "127.0.0.1")
	os.Setenv(node.Port, "9011")
	os.Setenv(node.Path, utils.TmpDirPath(t))
	go node.MutServer()
	os.Setenv(node.Host, "127.0.0.1")
	os.Setenv(node.Port, "9012")
	os.Setenv(node.Path, utils.TmpDirPath(t))
	go node.MutServer()
	os.Setenv(node.Host, "127.0.0.1")
	os.Setenv(node.Port, "9013")
	os.Setenv(node.Path, utils.TmpDirPath(t))
	go node.MutServer()
	time.Sleep(time.Minute * 10)
}
