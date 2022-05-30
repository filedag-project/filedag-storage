package main

import (
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
		Name:                 "file-dag-storage-dagpool",
		Usage:                "file-dag-storage-dagpool",
		Version:              "0.0.3",
		EnableBashCompletion: true,
		Commands:             local,
	}
	app.Setup()
	app.Run(os.Args)
}
