package main

import (
	"context"
	"flag"
	"fmt"
	dagpoolcli "github.com/filedag-project/filedag-storage/dag/pool/client"
	"github.com/filedag-project/filedag-storage/objectservice/iam"
	"github.com/filedag-project/filedag-storage/objectservice/iam/auth"
	"github.com/filedag-project/filedag-storage/objectservice/iamapi"
	"github.com/filedag-project/filedag-storage/objectservice/s3api"
	"github.com/filedag-project/filedag-storage/objectservice/uleveldb"
	"github.com/filedag-project/filedag-storage/objectservice/utils"
	"github.com/gorilla/mux"
	"github.com/ipfs/go-blockservice"
	offline "github.com/ipfs/go-ipfs-exchange-offline"
	"github.com/ipfs/go-merkledag"
	"log"
	"net/http"
	"os"
)

//go run -tags example main.go daemon --pool-user=pool --pool-user-pass=pool123
func main() {
	var leveldbPath, port, poolAddr, poolUser, poolPass string
	f := flag.NewFlagSet("daemon", flag.ExitOnError)
	f.StringVar(&leveldbPath, "db-path", "/tmp/leveldb2/fds.db", "set db path default:`/tmp/leveldb2/pool.db`")
	f.StringVar(&port, "port", ":9985", "set listen addr default:`localhost:50001`")
	f.StringVar(&poolAddr, "pool-addr", "localhost:50001", "set the pool addr you want connect")
	f.StringVar(&poolUser, "pool-user", "", "set pool user ")
	f.StringVar(&poolPass, "pool-user-pass", "", "set pool user pass")

	switch os.Args[1] {
	case "daemon":
		f.Parse(os.Args[2:])
		if poolUser == "" || poolPass == "" {
			fmt.Printf("db-path:%v, port:%v, pool-addr:%v, pool-user:%v, pool-user-pass:%v", leveldbPath, port, poolAddr, poolUser, poolPass)
			fmt.Println("please check your input\n " +
				"USAGE ERROR: go daemon -tags example main.go run daemon --db-path= --port= --pool-addr= --pool-user= pool-user-pass=")
		} else {
			run(leveldbPath, port, poolAddr, poolUser, poolPass)
		}
	default:
		fmt.Println("expected 'daemon' subcommands")
		os.Exit(1)
	}
}
func run(leveldbPath, port, poolAddr, poolUser, poolPass string) {
	db, err := uleveldb.OpenDb(leveldbPath)
	if err != nil {
		fmt.Printf("OpenDb err:%v", err)
		return
	}
	defer db.Close()
	cred, err := auth.CreateCredentials(auth.DefaultAccessKey, auth.DefaultSecretKey)
	if err != nil {
		println(err)
		return
	}
	authSys := iam.NewAuthSys(db, cred)
	router := mux.NewRouter()
	iamapi.NewIamApiServer(router, authSys)
	poolClient, err := dagpoolcli.NewPoolClient(poolAddr, poolUser, poolPass, true)
	if err != nil {
		log.Fatalf("connect dagpool server err: %v", err)
	}
	defer poolClient.Close(context.TODO())
	dagServ := merkledag.NewDAGService(blockservice.New(poolClient, offline.Exchange(poolClient)))
	s3api.NewS3Server(context.TODO(), router, dagServ, authSys, db)

	for _, ip := range utils.MustGetLocalIP4().ToSlice() {
		fmt.Printf("start sever at http://%v%v", ip, port)
	}
	err = http.ListenAndServe(port, router)
	if err != nil {
		fmt.Printf("Listen And Serve err%v", err)
		return
	}
}
