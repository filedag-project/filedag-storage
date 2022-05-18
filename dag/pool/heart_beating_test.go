package pool

import (
	"github.com/filedag-project/filedag-storage/dag/config"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	"github.com/filedag-project/filedag-storage/http/objectstore/utils"
	logging "github.com/ipfs/go-log/v2"
	"testing"
	"time"
)

func TestHeart_beating(t *testing.T) {
	logging.SetLogLevel("*", "INFO")
	db, err := uleveldb.OpenDb(utils.TmpDirPath(t))
	if err != nil {
		log.Errorf("err %v", err)
	}
	r := NewRecordSys(db)
	var a []config.CaskConfig
	a = append(a, config.CaskConfig{HeartAddr: ":7373"})
	err = r.HandleDagNode(a, "test")
	if err != nil {
		return
	}
	//go mutcask.SendHeartBeat("127.0.0.1:7373")
	time.Sleep(time.Minute)
}
