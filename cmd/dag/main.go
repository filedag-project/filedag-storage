package main

import (
	logging "github.com/ipfs/go-log/v2"
	"os"
)

func main() {
	logging.SetLogLevel("*", "INFO")
	os.Setenv(DagPoolLeveldbPath, "/tmp/leveldb2/pool")

	os.Setenv(DagNodeConfig, "dag/config/node_config.json")

	os.Setenv(DagPoolImporterBatchNum, "4")
	os.Setenv(DagPoolAddr, "localhost:50001")
	startDagPoolServer()
}
