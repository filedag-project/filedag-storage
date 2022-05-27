package main

import (
	"flag"
	"fmt"
	"github.com/filedag-project/filedag-storage/http/objectstore/iamapi"
	"github.com/filedag-project/filedag-storage/http/objectstore/s3api"
	"github.com/filedag-project/filedag-storage/http/objectstore/store"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	"github.com/filedag-project/filedag-storage/http/objectstore/utils"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

//go run -tags example main.go run --pool-user=pool --pool-user-pass=pool123
func main() {
	var leveldbPath, port, poolAddr, poolUser, poolPass string
	f := flag.NewFlagSet("run", flag.ExitOnError)
	f.StringVar(&leveldbPath, "db-path", "/tmp/leveldb2/fds.db", "set db path default:`/tmp/leveldb2/pool.db`")
	f.StringVar(&port, "port", ":9985", "set listen addr default:`localhost:50001`")
	f.StringVar(&poolAddr, "pool-addr", "localhost:50001", "set the pool addr you want connect")
	f.StringVar(&poolUser, "pool-user", "", "set pool user ")
	f.StringVar(&poolPass, "pool-user-pass", "", "set pool user pass")

	switch os.Args[1] {
	case "run":
		f.Parse(os.Args[2:])
		if poolUser == "" || poolPass == "" {
			fmt.Printf("db-path:%v, port:%v, pool-addr:%v, pool-user:%v, pool-user-pass:%v", leveldbPath, port, poolAddr, poolUser, poolPass)
			fmt.Println("please check your input\n " +
				"USAGE ERROR: go run -tags example main.go run --db-path= --port= --pool-addr= --pool-user= pool-user-pass=")
		} else {
			run(leveldbPath, port, poolAddr, poolUser, poolPass)
		}
	default:
		fmt.Println("expected 'str' subcommands")
		os.Exit(1)
	}
}
func run(leveldbPath, port, poolAddr, poolUser, poolPass string) {
	var err error
	uleveldb.DBClient, err = uleveldb.OpenDb(leveldbPath)
	if err != nil {
		fmt.Printf("OpenDb err:%v", err)
		return
	}
	defer uleveldb.DBClient.Close()
	os.Setenv(store.PoolAddr, poolAddr)
	os.Setenv(store.PoolUser, poolUser)
	os.Setenv(store.PoolPass, poolPass)
	router := mux.NewRouter()
	iamapi.NewIamApiServer(router)
	s := s3api.NewS3Server(router)
	if s == nil {
		fmt.Printf("may be pool addr not right,please check your pool-addr")
		return
	}
	defer s.Close()

	for _, ip := range utils.MustGetLocalIP4().ToSlice() {
		fmt.Printf("start sever at http://%v%v", ip, port)
	}
	err = http.ListenAndServe(port, router)
	if err != nil {
		fmt.Printf("Listen And Serve err%v", err)
		return
	}
}
