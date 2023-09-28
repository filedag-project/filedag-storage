package main

import (
	"fmt"
	"github.com/filedag-project/filedag-storage/objectservice/utils"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	utils.SetupLogLevels()
	local := []*cli.Command{
		startCmd,
		authCmd,
		clusterCmd,
	}
	app := &cli.App{
		Name:                 "dagpool",
		Usage:                "dag pool cluster",
		Version:              "0.1.0",
		EnableBashCompletion: true,
		Commands:             local,
	}
	app.Setup()
	if err := app.Run(os.Args); err != nil {
		fmt.Println("Error: ", err)
	}
}
