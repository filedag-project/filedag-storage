package config

import (
	"github.com/filedag-project/filedag-storage/dag/slotsmgr"
	"time"
)

//PoolConfig is the configuration for the dag pool
type PoolConfig struct {
	Listen       string        `json:"listen"`
	LeveldbPath  string        `json:"leveldb_path"`
	RootUser     string        `json:"root_user"`
	RootPassword string        `json:"root_password"`
	GcPeriod     time.Duration `json:"gc_period"`
}

//ClusterConfig is the configuration for a cluster
type ClusterConfig struct {
	Version int           `json:"version"`
	Cluster []DagNodeInfo `json:"cluster"`
}

//DagNodeConfig is the configuration for a dag node
type DagNodeConfig struct {
	Name         string   `json:"name"`
	Nodes        []string `json:"nodes"`         // rpc address list of datanodes
	DataBlocks   int      `json:"data_blocks"`   // Number of data shards
	ParityBlocks int      `json:"parity_blocks"` // Number of parity shards
}

type DagNodeInfo struct {
	Config    DagNodeConfig       `json:"config"`
	SlotPairs []slotsmgr.SlotPair `json:"slot_pairs"`
}
