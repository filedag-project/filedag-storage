package main

import (
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
	defaultPort          = ":9985"
	fileDagStoragePort   = "FILE_DAG_STORAGE_PORT"
	dbPath               = "DBPATH"
	defaultPoolStorePath = "./"
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
			Name:  "pool-path",
			Usage: "set pool path  eg.`.`",
			Value: defaultPoolStorePath,
		},
		&cli.StringFlag{
			Name:  "pool-batch-num",
			Usage: "set pool batch num eg.10",
			Value: defaultPoolBatchNum,
		},
		&cli.StringFlag{
			Name:  "pool-cask-num",
			Usage: "set pool cask num.:10",
			Value: defaultPoolCaskNum,
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
		if cctx.String("pool-path") != "" {
			err := os.Setenv(store.PoolStorePath, cctx.String("pool-path"))
			if err != nil {
				return err
			}
		}
		if cctx.String("pool-batch-num") != "" {
			err := os.Setenv(store.PoolBatchNum, cctx.String("pool-batch-num"))
			if err != nil {
				return err
			}
		}
		if cctx.String("pool-cask-num") != "" {
			err := os.Setenv(store.PoolCaskNum, cctx.String("pool-cask-num"))
			if err != nil {
				return err
			}
		}
		startServer()
		return nil
	},
}
