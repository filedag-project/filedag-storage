package main

import (
	"fmt"
	"github.com/filedag-project/filedag-storage/http/objectstore/utils"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	utils.SetupLogLevels()
	local := []*cli.Command{
		startCmd,
		authCmd,
	}
	app := &cli.App{
		Name:                 "dag-pool",
		Usage:                "dag pool daemon",
		Version:              "0.0.1",
		EnableBashCompletion: true,
		Commands:             local,
	}
	app.Setup()
	if err := app.Run(os.Args); err != nil {
		fmt.Println("Error: ", err)
	}
}
