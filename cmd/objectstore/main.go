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
