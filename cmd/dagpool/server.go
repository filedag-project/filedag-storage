package main

import (
	"encoding/json"
	"github.com/filedag-project/filedag-storage/dag/config"
	"github.com/filedag-project/filedag-storage/dag/pool"
	"github.com/filedag-project/filedag-storage/dag/pool/dagpooluser"
	"github.com/filedag-project/filedag-storage/dag/pool/server"
	"github.com/filedag-project/filedag-storage/dag/pool/userpolicy"
	"github.com/filedag-project/filedag-storage/dag/proto"
	logging "github.com/ipfs/go-log/v2"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	DagPoolLeveldbPath      = "POOL_LEVELDB_PATH"
	DagPooListenAddr        = "POOL_ADDR"
	DagNodeConfigPath       = "NODE_CONFIG_PATH"
	DagPoolImporterBatchNum = "POOL_IMPORTER_BATCH_NUM"
)
const (
	defaultPoolDB           = "/tmp/leveldb2/pool.db"
	defaultPoolListenAddr   = "localhost:50001"
	defaultNodeConfig       = "dag/config/node_config.json"
	defaultImporterBatchNum = "4"
)
const (
	defaultUser = "pool"
	defaultPass = "pool123"
)

var log = logging.Logger("pool-main")
var startCmd = &cli.Command{
	Name:  "run",
	Usage: "Start a dag pool process",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "pool-db-path",
			Usage: "set db path default:`/tmp/leveldb2/pool.db`",
			Value: defaultPoolDB,
		},
		&cli.StringFlag{
			Name:  "listen-addr",
			Usage: "set listen addr default:`localhost:50001`",
			Value: defaultPoolListenAddr,
		},
		&cli.StringFlag{
			Name:  "node-config-path",
			Usage: "set node config path,default:`dag/config/node_config.json'",
			Value: defaultNodeConfig,
		},
		&cli.StringFlag{
			Name:  "importer-batch-num",
			Usage: "set importer batch num default:4",
			Value: defaultImporterBatchNum,
		},
	},
	Action: func(cctx *cli.Context) error {

		if cctx.String("pool-db-path") != "" {
			err := os.Setenv(DagPoolLeveldbPath, cctx.String("pool-db-path"))
			if err != nil {
				return err
			}
		}

		if cctx.String("listen-addr") != "" {
			err := os.Setenv(DagPooListenAddr, cctx.String("listen-addr"))
			if err != nil {
				return err
			}
		}
		if cctx.String("node-config-path") != "" {
			err := os.Setenv(DagNodeConfigPath, cctx.String("node-config-path"))
			if err != nil {
				return err
			}
		}
		if cctx.String("importer-batch-num") != "" {
			err := os.Setenv(DagPoolImporterBatchNum, cctx.String("importer-batch-num"))
			if err != nil {
				return err
			}
		}
		startDagPoolServer()
		return nil
	},
}

func startDagPoolServer() {
	// listen port
	lis, err := net.Listen("tcp", os.Getenv(DagPooListenAddr))
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
		log.Errorf("NewDagPoolService err:%v", err)
		return
	}
	//add default user
	err = service.Iam.AddUser(dagpooluser.DagPoolUser{
		Username: defaultUser,
		Password: defaultPass,
		Policy:   userpolicy.ReadWrite,
		Capacity: 0,
	})
	if err != nil {
		return
	}
	proto.RegisterDagPoolServer(s, &server.DagPoolService{DagPool: service})
	log.Infof("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Errorf("failed to serve: %v", err)
	}
}
func LoadPoolConfig() (config.PoolConfig, error) {
	p := os.Getenv(DagNodeConfigPath)
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
