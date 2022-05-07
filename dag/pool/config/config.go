package config

import "github.com/filedag-project/filedag-storage/dag/node"

type SimplePoolConfig struct {
	NodesConfig []node.Config
	//todo more
	LeveldbPath      string
	ImporterBatchNum int
}
type PoolConfig struct {
	NodesConfig []node.Config
	//todo more
	LeveldbPath      string
	ImporterBatchNum int
}
