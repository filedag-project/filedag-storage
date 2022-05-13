package main

import (
	logging "github.com/ipfs/go-log/v2"
	"os"
)

const (
	defaultNodeConfig       = "dag/config/node_config2.json"
	defaultImporterBatchNum = "4"
	defaultPoolDB           = "/tmp/leveldb2/pool.db"
	defaultPoolAddr         = "localhost:50001"
)

func main() {
	logging.SetLogLevel("*", "INFO")
	os.Setenv(DagPoolLeveldbPath, defaultPoolDB)

	os.Setenv(DagNodeConfig, defaultNodeConfig)

	os.Setenv(DagPoolImporterBatchNum, defaultImporterBatchNum)
	os.Setenv(DagPoolAddr, defaultPoolAddr)

	startDagPoolServer()
}
