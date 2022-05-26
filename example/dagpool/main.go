package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/config"
	"github.com/filedag-project/filedag-storage/dag/pool"
	"github.com/filedag-project/filedag-storage/dag/pool/dagpooluser"
	"github.com/filedag-project/filedag-storage/dag/pool/server"
	"github.com/filedag-project/filedag-storage/dag/pool/userpolicy"
	"github.com/filedag-project/filedag-storage/dag/proto"
	"google.golang.org/grpc"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
)

//go run -tags example main.go --pool-db-path=/tmp/leveldb2/pool.db --listen-addr=localhost:50001 --node-config-path=node_config.json --importer-batch-num=4
func main() {
	var leveldbPath, listenAddr, nodeConfigPath, importerBatchNum string
	flag.StringVar(&leveldbPath, "pool-db-path", "", "set db path default:`/tmp/leveldb2/pool.db`")
	flag.StringVar(&listenAddr, "listen-addr", "", "set listen addr default:`localhost:50001`")
	flag.StringVar(&nodeConfigPath, "node-config-path", "", "set node config path,default:`dag/config/node_config.json'")
	flag.StringVar(&importerBatchNum, "importer-batch-num", "", "set importer batch num default:4")
	flag.Parse()
	if leveldbPath == "" || listenAddr == "" || nodeConfigPath == "" || importerBatchNum == "" {
		fmt.Printf("leveldbPath:%v, listenAddr:%v, nodeConfigPath:%v, importerBatchNum:%v", leveldbPath, listenAddr, nodeConfigPath, importerBatchNum)
		fmt.Println("please check your input\n " +
			"USAGE ERROR: go run -tags example main.go --pool-db-path= --listen-addr= --node-config-path= --importer-batch-num=")
	}

	run(leveldbPath, listenAddr, nodeConfigPath, importerBatchNum)
}

func run(leveldbPath, listenAddr, nodeConfigPath, importerBatchNum string) {
	// listen port
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		fmt.Printf("failed to listen: %v\n", err)
	}
	// new server
	s := grpc.NewServer()
	a, _ := strconv.Atoi(importerBatchNum)
	var nodeConfigs []config.NodeConfig
	for _, path := range strings.Split(nodeConfigPath, ",") {
		var nc config.NodeConfig
		file, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Printf("ReadFile err:%v\n", err)
			return
		}
		err = json.Unmarshal(file, &nc)
		if err != nil {
			fmt.Printf("Unmarshal err:%v\n", err)
			return
		}
		nodeConfigs = append(nodeConfigs, nc)
	}
	cfg := config.PoolConfig{
		DagNodeConfig:    nodeConfigs,
		LeveldbPath:      leveldbPath,
		ImporterBatchNum: a,
	}
	service, err := pool.NewDagPoolService(cfg)
	if err != nil {
		fmt.Printf("NewDagPoolService err:%v\n", err)
		return
	}
	//add default user
	err = service.Iam.AddUser(dagpooluser.DagPoolUser{
		Username: "pool",
		Password: "pool123",
		Policy:   userpolicy.ReadWrite,
		Capacity: 0,
	})
	if err != nil {
		return
	}
	proto.RegisterDagPoolServer(s, &server.DagPoolService{DagPool: service})
	fmt.Printf("server listening at %v\n", lis.Addr())
	if err := s.Serve(lis); err != nil {
		fmt.Printf("failed to serve: %v\n", err)
	}
}
