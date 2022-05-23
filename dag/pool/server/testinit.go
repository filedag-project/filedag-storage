package server

import (
	"github.com/filedag-project/filedag-storage/dag/config"
	"github.com/filedag-project/filedag-storage/dag/pool"
	"github.com/filedag-project/filedag-storage/dag/pool/dagpooluser"
	"github.com/filedag-project/filedag-storage/dag/pool/userpolicy"
	"github.com/filedag-project/filedag-storage/http/objectstore/utils"
	logging "github.com/ipfs/go-log/v2"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"testing"
	"time"
)

// StartTestDagPoolServer only for test
func StartTestDagPoolServer(t *testing.T) {
	logging.SetLogLevel("*", "INFO")
	//go mutcask.MutServer("127.0.0.1", "9010", utils.TmpDirPath(t))
	//go mutcask.MutServer("127.0.0.1", "9011", utils.TmpDirPath(t))
	//go mutcask.MutServer("127.0.0.1", "9012", utils.TmpDirPath(t))
	time.Sleep(time.Millisecond * 500)
	// listen port
	lis, err := net.Listen("tcp", "localhost:50001")
	if err != nil {
		log.Errorf("failed to listen: %v", err)
	}
	// new server
	s := grpc.NewServer()
	con, err := loadTestPoolConfig(t)
	if err != nil {
		return
	}
	service, err := pool.NewDagPoolService(con)
	if err != nil {
		return
	}
	//add default user
	service.Iam.AddUser(dagpooluser.DagPoolUser{
		Username: "pool",
		Password: "pool123",
		Policy:   userpolicy.ReadWrite,
		Capacity: 0,
	})
	RegisterDagPoolServer(s, &DagPoolService{DagPool: service})
	log.Infof("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Errorf("failed to serve: %v", err)
	}
}

func loadTestPoolConfig(t *testing.T) (cfg config.PoolConfig, err error) {
	cfg.LeveldbPath = utils.TmpDirPath(t)
	cfg.ImporterBatchNum = 4
	var caskc []config.CaskConfig
	for i := 0; i < 3; i++ {
		caskc = append(caskc, config.CaskConfig{Ip: "127.0.0.1", Port: strconv.Itoa(9010 + i), HeartAddr: ":" + strconv.Itoa(7373+i)})
	}
	var c = config.NodeConfig{
		Nodes:        caskc,
		DataBlocks:   3,
		ParityBlocks: 2,
		LevelDbPath:  utils.TmpDirPath(t),
	}
	cfg.DagNodeConfig = append(cfg.DagNodeConfig, c)
	return cfg, nil
}
