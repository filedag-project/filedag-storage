package main

import (
	"encoding/json"
	"fmt"
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
	"strconv"
	"strings"
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
		&cli.StringFlag{
			Name:  "pool-user",
			Usage: "set root user default:pool",
			Value: defaultUser,
		},
		&cli.StringFlag{
			Name:  "pool-pass",
			Usage: "set root user pass default:pool123",
			Value: defaultPass,
		},
	},
	Action: func(cctx *cli.Context) error {
		var (
			dbpath           = defaultPoolDB
			addr             = defaultPoolListenAddr
			nodeConfigPath   = defaultNodeConfig
			importerBatchNum = defaultImporterBatchNum
			poolUser         = defaultUser
			poolPass         = defaultPass
		)
		if cctx.String("pool-db-path") != "" {
			dbpath = cctx.String("pool-db-path")
		} else {
			fmt.Println("use default pool db path:", defaultPoolDB)
		}

		if cctx.String("listen-addr") != "" {
			addr = cctx.String("listen-addr")
		} else {
			fmt.Println("use default listen addr:", defaultPoolListenAddr)
		}
		if cctx.String("node-config-path") != "" {
			nodeConfigPath = cctx.String("node-config-path")
		} else {
			fmt.Println("use default node config path:", defaultNodeConfig)
		}
		if cctx.String("importer-batch-num") != "" {
			importerBatchNum = cctx.String("importer-batch-num")
		} else {
			fmt.Println("use default importer batch num:", defaultImporterBatchNum)
		}
		if cctx.String("pool-user") != "" {
			poolUser = cctx.String("pool-user")
		} else {
			fmt.Println("use default pool user:", defaultUser)
		}
		if cctx.String("pool-pass") != "" {
			poolPass = cctx.String("pool-pass")
		} else {
			fmt.Println("use default pool pass:", defaultPass)
		}
		startDagPoolServer(dbpath, addr, nodeConfigPath, importerBatchNum, poolUser, poolPass)
		return nil
	},
}

func startDagPoolServer(dbpath, addr, nodeConfigPath, importerBatchNum, poolUser, poolPass string) {
	// listen port
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Errorf("failed to listen: %v", err)
	}
	// new server
	s := grpc.NewServer()
	con, err := LoadPoolConfig(dbpath, nodeConfigPath, importerBatchNum, poolUser, poolPass)
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
		Username: poolUser,
		Password: poolPass,
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
func LoadPoolConfig(dbpath, nodeConfigPath, importerBatchNum, poolUser, poolPass string) (config.PoolConfig, error) {
	i, _ := strconv.Atoi(importerBatchNum)
	var nodeConfigs []config.NodeConfig
	for _, path := range strings.Split(nodeConfigPath, ",") {
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
		LeveldbPath:      dbpath,
		ImporterBatchNum: i,
		DefaultUser:      poolUser,
		DefaultPass:      poolPass,
	}, nil
}
