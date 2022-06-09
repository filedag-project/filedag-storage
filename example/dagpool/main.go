package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/config"
	"github.com/filedag-project/filedag-storage/dag/pool/dagpooluser"
	"github.com/filedag-project/filedag-storage/dag/pool/poolservice"
	"github.com/filedag-project/filedag-storage/dag/pool/server"
	"github.com/filedag-project/filedag-storage/dag/pool/userpolicy"
	"github.com/filedag-project/filedag-storage/dag/proto"
	"google.golang.org/grpc"
	"io/ioutil"
	"net"
	"os"
	"path"
	"strings"
)

//go run -tags example main.go daemon --pool-db-path=/tmp/leveldb2/pool.db --listen-addr=localhost:50001 --node-config-path=node_config.json --pool-user=pool --pool-pass=pool123
func main() {
	var leveldbPath, listenAddr, nodeConfigPath, user, pass string
	f := flag.NewFlagSet("daemon", flag.ExitOnError)
	f.StringVar(&leveldbPath, "pool-db-path", "/tmp/leveldb2/", "set db path default:`/tmp/leveldb2/pool.db`")
	f.StringVar(&listenAddr, "listen-addr", "localhost:50001", "set listen addr default:`localhost:50001`")
	f.StringVar(&nodeConfigPath, "node-config-path", "node_config.json", "set node config path,default:`dag/config/node_config.json'")
	f.StringVar(&user, "pool-user", "pool", "set root user default:pool")
	f.StringVar(&pass, "pool-pass", "pool123", "set root user pass default:pool123")
	leveldbPath = path.Join(leveldbPath, "leveldb")
	switch os.Args[1] {
	case "daemon":
		f.Parse(os.Args[2:])
		if leveldbPath == "" || listenAddr == "" || nodeConfigPath == "" || user == "" || pass == "" {
			fmt.Printf("leveldbPath:%v, listenAddr:%v, nodeConfigPath:%v,user:%v,pass:%v", leveldbPath, listenAddr, nodeConfigPath, user, pass)
			fmt.Println("please check your input\n " +
				"USAGE ERROR: go run -tags example main.go daemon --pool-db-path= --listen-addr= --node-config-path= --pool-user= --pool-pass= --datastore-path=")
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
	var nodeConfigs []config.NodeConfig
	for _, p := range strings.Split(nodeConfigPath, ",") {
		var nc config.NodeConfig
		file, err := ioutil.ReadFile(p)
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
		DefaultUser:   user,
		DefaultPass:   pass,
	}
	service, err := poolservice.NewDagPoolService(cfg)
	if err != nil {
		fmt.Printf("NewDagPoolService err:%v\n", err)
		return
	}
	//add default user
	err = service.AddUser(dagpooluser.DagPoolUser{
		Username: user,
		Password: pass,
		Policy:   userpolicy.ReadWrite,
		Capacity: 0,
	})
	if err != nil {
		return
	}
	proto.RegisterDagPoolServer(s, &server.DagPoolServer{DagPool: service})
	fmt.Printf("server listening at %v\n", lis.Addr())
	if err := s.Serve(lis); err != nil {
		fmt.Printf("failed to serve: %v\n", err)
	}
}
