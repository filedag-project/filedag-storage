package main

import (
	"fmt"
	"github.com/filedag-project/filedag-storage/http/objectstore/iamapi"
	"github.com/filedag-project/filedag-storage/http/objectstore/s3api"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	"github.com/filedag-project/filedag-storage/http/objectstore/utils"
	"github.com/gorilla/mux"
	logging "github.com/ipfs/go-log/v2"
	"github.com/urfave/cli/v2"
	"golang.org/x/xerrors"
	"net/http"
)

var log = logging.Logger("sever")

const (
	deFaultDBFILE   = "/tmp/leveldb2/fds.db"
	defaultPort     = ":9985"
	defaultPoolAddr = "localhost:50001"
)

//startServer Start a IamServer
func startServer(dbPath, port, poolAddr, poolUser, poolPass string) {
	var err error
	uleveldb.DBClient, err = uleveldb.OpenDb(dbPath)
	if err != nil {
		return
	}
	defer uleveldb.DBClient.Close()
	router := mux.NewRouter()
	iamapi.NewIamApiServer(router)
	s := s3api.NewS3Server(router, poolAddr, poolUser, poolPass)
	if s == nil {
		log.Errorf("may be pool addr not right,please check your pool-addr")
		return
	}
	defer s.Close()
	for _, ip := range utils.MustGetLocalIP4().ToSlice() {
		log.Infof("start sever at http://%v%v", ip, port)
	}
	err = http.ListenAndServe(port, router)
	if err != nil {
		log.Errorf("Listen And Serve err%v", err)
		return
	}
}

var startCmd = &cli.Command{
	Name:  "run",
	Usage: "Start a file dag storage process",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "db-path",
			Usage: "set db path",
			Value: deFaultDBFILE,
		},
		&cli.StringFlag{
			Name:  "port",
			Usage: "set port eg.:9985",
			Value: defaultPort,
		},
		&cli.StringFlag{
			Name:  "pool-addr",
			Usage: "set the pool addr you want connect",
			Value: defaultPoolAddr,
		},
		&cli.StringFlag{
			Name:  "pool-user",
			Usage: "set pool user",
		},
		&cli.StringFlag{
			Name:  "pool-user-pass",
			Usage: "set pool user pass",
		},
	},
	Action: func(cctx *cli.Context) error {
		var poolUser, poolPass string
		var (
			dbPath   = deFaultDBFILE
			port     = defaultPort
			poolAddr = defaultPoolAddr
		)
		if cctx.String("db-path") != "" {
			dbPath = cctx.String("db-path")
		} else {
			fmt.Println("use default db path:", deFaultDBFILE)
		}

		if cctx.String("port") != "" {
			port = cctx.String("port")
		} else {
			fmt.Println("use default port:", defaultPort)
		}
		if cctx.String("pool-addr") != "" {
			poolAddr = cctx.String("pool-addr")
		} else {
			fmt.Println("use default pool addr:", defaultPoolAddr)
		}
		if cctx.String("pool-user") != "" {
			poolUser = cctx.String("pool-user")
		} else {
			fmt.Println("please set pool user ,eg:--pool-user=pool")
			return xerrors.Errorf("please set pool user ,eg:--pool-user=pool")
		}
		if cctx.String("pool-user-pass") != "" {
			poolPass = cctx.String("pool-user-pass")
		} else {
			fmt.Println("please set pool user pass,eg:--pool-user-pass=pool123")
			return xerrors.Errorf("please set pool user pass,eg:--pool-user-pass=pool123")
		}
		startServer(dbPath, port, poolAddr, poolUser, poolPass)
		return nil
	},
}
