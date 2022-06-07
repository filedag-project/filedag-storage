package main

import (
	logging "github.com/ipfs/go-log/v2"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	logging.SetLogLevel("*", "INFO")
	local := []*cli.Command{
		addBlock,
		getBlock,
		removeBlock,

		addUser,
		removeUser,
		getUser,
		updateUser,
	}
	app := &cli.App{
		Name:                 "dagpool-client",
		Usage:                "sent rpc request to dagpool",
		Version:              "0.0.4",
		EnableBashCompletion: true,
		Commands:             local,
	}
	app.Setup()
	app.Run(os.Args)
}
