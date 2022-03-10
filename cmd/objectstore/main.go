package main

import (
	logging "github.com/ipfs/go-log/v2"
	"github.com/urfave/cli/v2"
	"os"
)

const (
	dbFILE = "/tmp/leveldb2/fds.db"
)

func main() {
	logging.SetLogLevel("*", "INFO")
	local := []*cli.Command{
		startCmd,
	}
	app := &cli.App{
		Name:                 "file-dag-storage",
		Usage:                "file-dag-storage",
		Version:              "0.0.11",
		EnableBashCompletion: true,
		Commands:             local,
	}
	app.Setup()
	app.Run(os.Args)

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
