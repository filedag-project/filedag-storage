package main

import (
	"fmt"
	"github.com/filedag-project/filedag-storage/kv/mutcask"
	logging "github.com/ipfs/go-log/v2"
	"github.com/urfave/cli/v2"
	"os"
)

const (
	defaultHost = "localhost"
	defaultPort = "9010"
	defaultPath = "/tmp/dag/data1"
)

func main() {
	logging.SetLogLevel("*", "INFO")
	local := []*cli.Command{
		startCmd,
	}
	app := &cli.App{
		Name:     "mut-cask",
		Usage:    "store data",
		Version:  "0.0.1",
		Commands: local,
	}
	app.Setup()
	if err := app.Run(os.Args); err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}

var startCmd = &cli.Command{
	Name:  "run",
	Usage: "Start a data node process",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "host",
			Usage: "set host eg.:127.0.0.1",
			Value: defaultHost,
		},
		&cli.StringFlag{
			Name:  "port",
			Usage: "set port eg.:9010",
			Value: defaultPort,
		},
		&cli.StringFlag{
			Name:  "path",
			Usage: "set data node path",
			Value: defaultPath,
		},
	},
	Action: func(c *cli.Context) error {
		if c.String("host") != "" {
			err := os.Setenv(defaultHost, c.String("host"))
			if err != nil {
				return err
			}
		}
		if c.String("port") != "" {
			err := os.Setenv(defaultPort, c.String("port"))
			if err != nil {
				return err
			}
		}
		if c.String("path") != "" {
			err := os.Setenv(defaultPath, c.String("path"))
			if err != nil {
				return err
			}
		}
		mutcask.MutServer(defaultHost, defaultPort, defaultPath)
		return nil
	},
}
