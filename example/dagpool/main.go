package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/config"
	"github.com/filedag-project/filedag-storage/dag/pool/poolservice"
	"github.com/filedag-project/filedag-storage/dag/pool/server"
	"github.com/filedag-project/filedag-storage/dag/proto"
	"google.golang.org/grpc"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

// go run -tags example main.go daemon --datadir=/tmp/dagpool-db --listen=localhost:50001 --config=node_config.json --root-user=dagpool --root-password=dagpool
func main() {
	var leveldbPath, listenAddr, nodeConfigPath, user, pass string
	f := flag.NewFlagSet("daemon", flag.ExitOnError)
	f.StringVar(&leveldbPath, "datadir", "/tmp/leveldb2/pool.db", "directory to store data in")
	f.StringVar(&listenAddr, "listen", "localhost:50001", "set server listen")
	f.StringVar(&nodeConfigPath, "config", "node_config.json", "set config path")
	f.StringVar(&user, "root-user", "dagpool", "set root user")
	f.StringVar(&pass, "root-password", "dagpool", "set root password")

	switch os.Args[1] {
	case "daemon":
		f.Parse(os.Args[2:])
		if leveldbPath == "" || listenAddr == "" || nodeConfigPath == "" || user == "" || pass == "" {
			fmt.Printf("leveldbPath:%v, listenAddr:%v, nodeConfigPath:%v,user:%v,pass:%v", leveldbPath, listenAddr, nodeConfigPath, user, pass)
			fmt.Println("please check your input\n " +
				"USAGE ERROR: go run -tags example main.go daemon --datadir=/tmp/dagpool-db --listen=localhost:50001 --config=node_config.json --root-user=dagpool --root-password=dagpool")
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
		LeveldbPath:  leveldbPath,
		RootUser:     user,
		RootPassword: pass,
	}
	service, err := poolservice.NewDagPoolService(context.TODO(), cfg)
	if err != nil {
		fmt.Printf("NewDagPoolService err:%v\n", err)
		return
	}
	for _, nd := range nodeConfigs {
		err = service.AddDagNode(&nd)
		if err != nil {
			fmt.Printf("AddDagNode err:%v\n", err)
			return
		}
	}
	if err = service.BalanceSlots(); err != nil {
		fmt.Printf("BalanceSlots err:%v\n", err)
		return
	}
	proto.RegisterDagPoolServer(s, &server.DagPoolServer{DagPool: service})
	fmt.Printf("server listening at %v\n", lis.Addr())
	if err := s.Serve(lis); err != nil {
		fmt.Printf("failed to serve: %v\n", err)
	}
}
