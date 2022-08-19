package main

import (
	"fmt"
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
