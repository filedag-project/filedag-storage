package main

import (
	"errors"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/node"
	logging "github.com/ipfs/go-log/v2"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	logging.SetLogLevel("*", "INFO")
	local := []*cli.Command{
		startCmd,
	}
	app := &cli.App{
		Name:     "datanode",
		Usage:    "store data",
		Version:  "0.0.1",
		Commands: local,
	}
	app.Setup()
	if err := app.Run(os.Args); err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}

var startCmd = &cli.Command{
	Name:  "daemon",
	Usage: "Start a data node process",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "listen",
			Usage: "set server listen",
			Value: ":9010",
		},
		&cli.StringFlag{
			Name:  "datadir",
			Usage: "directory to store data in",
			Value: "./dn-data",
		},
		&cli.StringFlag{
			Name:  "kvdb",
			Usage: "choose kvdb, badger or mutcask",
			Value: "badger",
		},
	},
	Action: func(c *cli.Context) error {
		kvType := node.KVType(c.String("kvdb"))
		switch kvType {
		case node.KVBadge:
		case node.KVMutcask:
		default:
			return errors.New(fmt.Sprintf("not support this kvdb %s", kvType))
		}
		node.MutDataNodeServer(c.String("listen"), kvType, c.String("datadir"))
		return nil
	},
}
