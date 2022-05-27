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
	"os"
	"strconv"
	"strings"
)

//go run -tags example main.go run --pool-db-path=/tmp/leveldb2/pool.db --listen-addr=localhost:50001 --node-config-path=node_config.json --importer-batch-num=4
func main() {
	var leveldbPath, listenAddr, nodeConfigPath, importerBatchNum string
	f := flag.NewFlagSet("run", flag.ExitOnError)
	f.StringVar(&leveldbPath, "pool-db-path", "/tmp/leveldb2/pool.db", "set db path default:`/tmp/leveldb2/pool.db`")
	f.StringVar(&listenAddr, "listen-addr", "localhost:50001", "set listen addr default:`localhost:50001`")
	f.StringVar(&nodeConfigPath, "node-config-path", "node_config.json", "set node config path,default:`dag/config/node_config.json'")
	f.StringVar(&importerBatchNum, "importer-batch-num", "4", "set importer batch num default:4")

	switch os.Args[1] {
	case "run":
		f.Parse(os.Args[2:])
		if leveldbPath == "" || listenAddr == "" || nodeConfigPath == "" || importerBatchNum == "" {
			fmt.Printf("leveldbPath:%v, listenAddr:%v, nodeConfigPath:%v, importerBatchNum:%v", leveldbPath, listenAddr, nodeConfigPath, importerBatchNum)
			fmt.Println("please check your input\n " +
				"USAGE ERROR: go run -tags example main.go --pool-db-path= --listen-addr= --node-config-path= --importer-batch-num=")
		} else {
			run(leveldbPath, listenAddr, nodeConfigPath, importerBatchNum)
		}
	default:
		fmt.Println("expected 'str' subcommands")
		os.Exit(1)
	}

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
