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
		authCmd,
	}
	app := &cli.App{
		Name:                 "dagpool",
		Usage:                "dagpool",
		Version:              "0.0.3",
		EnableBashCompletion: true,
		Commands:             local,
	}
	app.Setup()
	app.Run(os.Args)
}
