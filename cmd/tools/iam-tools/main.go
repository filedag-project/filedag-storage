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
		addUserCmd,
		getUserCmd,
		delUserCmd,
		changePassCmd,
		setStatusCmd,
	}
	app := &cli.App{
		Name:                 "iam-tool",
		Usage:                "test some interface",
		Version:              "0.0.1",
		EnableBashCompletion: true,
		Commands:             local,
	}
	app.Setup()
	if err := app.Run(os.Args); err != nil {
		fmt.Println("Error: ", err)
	}
}
