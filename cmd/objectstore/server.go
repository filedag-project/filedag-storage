package main

import (
	"github.com/filedag-project/filedag-storage/http/objectstore/iamapi"
	"github.com/filedag-project/filedag-storage/http/objectstore/s3api"
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
	deFaultDBFILE  = "/tmp/leveldb2/fds.db"
	defaultPort    = ":9985"
	fileDagStorage = "FILE_DAG_STORAGE_PORT"
	dbPath         = "DBPATH"
)

//startServer Start a IamServer
func startServer() {
	uleveldb.DBClient = uleveldb.OpenDb(os.Getenv(dbPath))
	router := mux.NewRouter()
	s3api.NewS3Server(router)
	iamapi.NewIamApiServer(router)
	for _, ip := range utils.MustGetLocalIP4().ToSlice() {
		log.Infof("start sever at http://%v:%v", ip, 9985)
	}
	err := http.ListenAndServe(os.Getenv(fileDagStorage), router)
	if err != nil {
		log.Errorf("ListenAndServe err%v", err)
		return
	}
	defer uleveldb.DBClient.Close()
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
			Usage: "set port",
			Value: defaultPort,
		},
	},
	Action: func(cctx *cli.Context) error {

		if cctx.String("db-path") != "" {
			err := os.Setenv(dbPath, cctx.String("db-path"))
			if err != nil {
				return err
			}
		} else {
			err := os.Setenv(deFaultDBFILE, cctx.String(deFaultDBFILE))
			if err != nil {
				return err
			}
		}
		if cctx.String("port") != "" {
			err := os.Setenv(fileDagStorage, cctx.String("port"))
			if err != nil {
				return err
			}
		} else {
			err := os.Setenv(fileDagStorage, cctx.String(defaultPort))
			if err != nil {
				return err
			}
		}
		startServer()
		return nil
	},
}
