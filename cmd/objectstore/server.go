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

//startServer Start a IamServer
func startServer() {
	uleveldb.DBClient = uleveldb.OpenDb(os.Getenv("DBPATH"))
	router := mux.NewRouter()
	s3api.NewS3Server(router)
	iamapi.NewIamApiServer(router)
	for _, ip := range utils.MustGetLocalIP4().ToSlice() {
		log.Infof("start sever at http://%v:%v", ip, 9985)
	}
	err := http.ListenAndServe(":9985", router)
	if err != nil {
		log.Errorf("ListenAndServe err%v", err)

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
			Value: dbFILE,
		},
	},
	Action: func(cctx *cli.Context) error {

		if cctx.String("db-path") != "" {
			err := os.Setenv("DBPATH", cctx.String("db-path"))
			if err != nil {
				return err
			}
		} else {
			err := os.Setenv("DBPATH", cctx.String(dbFILE))
			if err != nil {
				return err
			}
		}
		startServer()
		return nil
	},
}
