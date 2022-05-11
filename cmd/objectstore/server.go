package main

import (
	"github.com/filedag-project/filedag-storage/dag/pool"
	"github.com/filedag-project/filedag-storage/http/objectstore/iamapi"
	"github.com/filedag-project/filedag-storage/http/objectstore/s3api"
	"github.com/filedag-project/filedag-storage/http/objectstore/store"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	"github.com/filedag-project/filedag-storage/http/objectstore/utils"
	"github.com/gorilla/mux"
	logging "github.com/ipfs/go-log/v2"
	"github.com/urfave/cli/v2"
	"net/http"
	"os"
)

var log = logging.Logger("sever")

const (
	deFaultDBFILE        = "/tmp/leveldb2/fds.db"
	deFaultPoolDBFILE    = "/tmp/leveldb2/pool.db"
	defaultPort          = ":9985"
	fileDagStoragePort   = "FILE_DAG_STORAGE_PORT"
	dbPath               = "DBPATH"
	defaultPoolStorePath = "./dag/node/config.json"
	defaultPoolBatchNum  = "4"
	defaultPoolCaskNum   = "2"
)

//startServer Start a IamServer
func startServer() {
	var err error
	uleveldb.DBClient, err = uleveldb.OpenDb(os.Getenv(dbPath))
	if err != nil {
		return
	}
	defer uleveldb.DBClient.Close()
	router := mux.NewRouter()
	iamapi.NewIamApiServer(router)
	s3api.NewS3Server(router)

	for _, ip := range utils.MustGetLocalIP4().ToSlice() {
		log.Infof("start sever at http://%v%v", ip, os.Getenv(fileDagStoragePort))
	}
	err = http.ListenAndServe(os.Getenv(fileDagStoragePort), router)
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
			Name:  "pool-db-path",
			Usage: "set pool db path,if you need a local pool , use it",
		},
		&cli.StringFlag{
			Name:  "pool-ip-path",
			Usage: "set node path  if you need a local node , use it",
		},
		&cli.StringFlag{
			Name:  "pool-batch-num",
			Usage: "set pool batch num if you need a local pool , use it",
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

		if cctx.String("db-path") != "" {
			err := os.Setenv(dbPath, cctx.String("db-path"))
			if err != nil {
				return err
			}
		}

		if cctx.String("port") != "" {
			err := os.Setenv(fileDagStoragePort, cctx.String("port"))
			if err != nil {
				return err
			}
		}
		if cctx.String("pool-db-path") != "" {
			os.Setenv(pool.DagPoolLeveldbPath, cctx.String("pool-db-path"))
		}
		if cctx.String("pool-ip-path") != "" {
			os.Setenv(pool.DagNodeIpOrPath, cctx.String("pool-ip-path"))
		}
		if cctx.String("pool-batch-num") != "" {
			os.Setenv(pool.DagPoolImporterBatchNum, cctx.String("pool-batch-num"))
		}
		if cctx.String("pool-user") != "" {
			os.Setenv(store.PoolUser, cctx.String("pool-user"))
		}
		if cctx.String("pool-user-pass") != "" {
			os.Setenv(store.PoolPass, cctx.String("pool-user-pass"))
		}
		startServer()
		return nil
	},
}
