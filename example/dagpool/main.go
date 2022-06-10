package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/config"
	"github.com/filedag-project/filedag-storage/dag/pool"
	"github.com/filedag-project/filedag-storage/dag/pool/server"
	"github.com/filedag-project/filedag-storage/dag/proto"
	"google.golang.org/grpc"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

//go run -tags example main.go daemon --pool-db-path=/tmp/dagpool-db --listen-addr=localhost:50001 --node-config-path=node_config.json --root-user=dagpool --root-pass=dagpool
func main() {
	var leveldbPath, listenAddr, nodeConfigPath, user, pass string
	f := flag.NewFlagSet("daemon", flag.ExitOnError)
	f.StringVar(&leveldbPath, "pool-db-path", "/tmp/leveldb2/pool.db", "set db path default:`/tmp/leveldb2/pool.db`")
	f.StringVar(&listenAddr, "listen-addr", "localhost:50001", "set listen addr default:`localhost:50001`")
	f.StringVar(&nodeConfigPath, "node-config-path", "node_config.json", "set node config path,default:`dag/config/node_config.json'")
	f.StringVar(&user, "root-user", "pool", "set root user default:pool")
	f.StringVar(&pass, "root-pass", "pool123", "set root user pass default:pool123")

	switch os.Args[1] {
	case "daemon":
		f.Parse(os.Args[2:])
		if leveldbPath == "" || listenAddr == "" || nodeConfigPath == "" || user == "" || pass == "" {
			fmt.Printf("leveldbPath:%v, listenAddr:%v, nodeConfigPath:%v,user:%v,pass:%v", leveldbPath, listenAddr, nodeConfigPath, user, pass)
			fmt.Println("please check your input\n " +
				"USAGE ERROR: go run -tags example main.go daemon --pool-db-path= --listen-addr= --node-config-path= --pool-user= --pool-pass=")
		} else {
			run(leveldbPath, listenAddr, nodeConfigPath, user, pass)
		}
	default:
		fmt.Println("expected 'daemon' subcommands")
		os.Exit(1)
	}

}

func run(leveldbPath, listenAddr, nodeConfigPath, user, pass string) {
	// listen port
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		fmt.Printf("failed to listen: %v\n", err)
	}
	// new server
	s := grpc.NewServer()
	var nodeConfigs []config.DagNodeConfig
	for _, path := range strings.Split(nodeConfigPath, ",") {
		var nc config.DagNodeConfig
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
		DagNodeConfig: nodeConfigs,
		LeveldbPath:   leveldbPath,
		RootUser:      user,
		RootPassword:  pass,
	}
	service, err := pool.NewDagPoolService(cfg)
	if err != nil {
		fmt.Printf("NewDagPoolService err:%v\n", err)
		return
	}
	proto.RegisterDagPoolServer(s, &server.DagPoolService{DagPool: service})
	fmt.Printf("server listening at %v\n", lis.Addr())
	if err := s.Serve(lis); err != nil {
		fmt.Printf("failed to serve: %v\n", err)
	}
}
