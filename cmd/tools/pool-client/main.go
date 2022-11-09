package main

import (
	"github.com/filedag-project/filedag-storage/objectservice/utils"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	utils.SetupLogLevels()
	local := []*cli.Command{
		addBlock,
		getBlock,
		removeBlock,
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
