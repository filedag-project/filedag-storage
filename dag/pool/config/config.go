package config

import "github.com/filedag-project/filedag-storage/dag/node"

type SimplePoolConfig struct {
	NodesConfig []node.Config
	//todo more
	LeveldbPath string
	StorePath   string // the path for kv db
	BatchNum    int    // use for blockstore batch task and importer batch task config
	CaskNum     int    // the number of vlog files in mutcask kv
}
