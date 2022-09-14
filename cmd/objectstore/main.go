package main

import (
	"fmt"
	"github.com/filedag-project/filedag-storage/http/objectservice/utils"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	utils.SetupLogLevels()
	local := []*cli.Command{
		startCmd,
	}
	app := &cli.App{
		Name:                 "filedag-storage",
		Usage:                "filedag-storage",
		Version:              "0.0.11",
		EnableBashCompletion: true,
		Commands:             local,
	}
	app.Setup()
	if err := app.Run(os.Args); err != nil {
		fmt.Println("Error: ", err)
	}
}
