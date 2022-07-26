package main

import (
	"fmt"
	logging "github.com/ipfs/go-log/v2"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	logging.SetLogLevel("*", "DEBUG")
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
	if err := app.Run(os.Args); err != nil {
		fmt.Println("Error: ", err)
	}
}
