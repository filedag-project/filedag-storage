package main

import (
	"encoding/json"
	"github.com/filedag-project/filedag-storage/dag/config"
	"github.com/filedag-project/filedag-storage/dag/pool"
	"github.com/filedag-project/filedag-storage/dag/pool/dagpooluser"
	"github.com/filedag-project/filedag-storage/dag/pool/server"
	"github.com/filedag-project/filedag-storage/dag/pool/userpolicy"
	logging "github.com/ipfs/go-log/v2"
	"google.golang.org/grpc"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	DagPoolLeveldbPath      = "POOL_LEVELDB_PATH"
	DagNodeConfig           = "POOL_IP_OR_PATH"
	DagPoolImporterBatchNum = "POOL_IMPORTER_BATCH_NUM"
	DagPoolAddr             = "POOL_ADDR"
)

var log = logging.Logger("pool-client")

func startDagPoolServer() {
	// listen port
	lis, err := net.Listen("tcp", os.Getenv(DagPoolAddr))
	if err != nil {
		log.Errorf("failed to listen: %v", err)
	}
	// new server
	s := grpc.NewServer()
	con, err := LoadPoolConfig()
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
		Password: "pool",
		Policy:   userpolicy.ReadWrite,
		Capacity: 0,
	})
	server.RegisterDagPoolServer(s, &server.DagPoolService{DagPool: service})
	log.Infof("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Errorf("failed to serve: %v", err)
	}
}
func LoadPoolConfig() (config.PoolConfig, error) {
	p := os.Getenv(DagNodeConfig)
	i := os.Getenv(DagPoolImporterBatchNum)
	importerBatchNum, _ := strconv.Atoi(i)
	var nodeConfigs []config.NodeConfig
	for _, path := range strings.Split(p, ",") {
		var nc config.NodeConfig
		file, err := ioutil.ReadFile(path)
		if err != nil {
			log.Errorf("ReadFile err:%v", err)
			return config.PoolConfig{}, err
		}
		err = json.Unmarshal(file, &nc)
		if err != nil {
			log.Errorf("Unmarshal err:%v", err)
			return config.PoolConfig{}, err
		}
		nodeConfigs = append(nodeConfigs, nc)
	}
	return config.PoolConfig{
		DagNodeConfig:    nodeConfigs,
		LeveldbPath:      os.Getenv(DagPoolLeveldbPath),
		ImporterBatchNum: importerBatchNum,
	}, nil
}
