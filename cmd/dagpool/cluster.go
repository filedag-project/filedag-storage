package main

import (
	"encoding/json"
	"github.com/filedag-project/filedag-storage/dag/config"
	"github.com/filedag-project/filedag-storage/dag/pool/poolservice"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"os"
	"path"
)

var clusterCmd = &cli.Command{
	Name:  "cluster",
	Usage: "Manage dagpool cluster nodes",
	Subcommands: []*cli.Command{
		initSlots,
		//slots,
		//addSlots,
		//addSlotsRange,
		//delSlots,
		//delSlotsRange,
	},
}

var initSlots = &cli.Command{
	Name:  "init",
	Usage: "Init slots of the dag pool cluster",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "datadir",
			Usage: "directory to store data in",
			Value: "./dp-data",
		},
		&cli.StringFlag{
			Name:  "config",
			Usage: "set config path",
			Value: "./conf/node_config.json",
		},
	},
	Action: func(cctx *cli.Context) error {
		var cfg config.PoolConfig
		datadir := cctx.String("datadir")
		if err := os.MkdirAll(datadir, 0777); err != nil {
			return err
		}
		cfg.LeveldbPath = path.Join(datadir, "leveldb")
		configPath := cctx.String("config")

		var clusterConfig config.ClusterConfig
		file, err := ioutil.ReadFile(configPath)
		if err != nil {
			log.Errorf("ReadFile err:%v", err)
			return err
		}
		err = json.Unmarshal(file, &clusterConfig)
		if err != nil {
			log.Errorf("Unmarshal err:%v", err)
			return err
		}
		cfg.ClusterConfig = clusterConfig

		return poolservice.AllocateSlotsEvenly(cfg)
	},
}
