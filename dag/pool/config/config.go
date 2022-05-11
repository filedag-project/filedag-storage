package config

import "github.com/filedag-project/filedag-storage/dag/node"

type SimplePoolConfig struct {
	NodesConfig []node.Config
	//todo more
	LeveldbPath      string
	ImporterBatchNum int
}
type PoolConfig struct {
	IpOrPath []string `json:"ip_or_path"`
	//todo more
	LeveldbPath      string `json:"leveldb_path"`
	ImporterBatchNum int    `json:"importer_batch_num"`
}
