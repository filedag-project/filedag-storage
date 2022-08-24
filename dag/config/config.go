package config

import "time"

//PoolConfig is the configuration for the dag pool
type PoolConfig struct {
	Listen        string          `json:"listen"`
	DagNodeConfig []DagNodeConfig `json:"dag_node"`
	LeveldbPath   string          `json:"leveldb_path"`
	RootUser      string          `json:"root_user"`
	RootPassword  string          `json:"root_password"`
	GcPeriod      time.Duration   `json:"gc_period"`
	CacheTimeout  time.Duration   `json:"cache_timeout"`
}

//DagNodeConfig is the configuration for a dag node
type DagNodeConfig struct {
	Nodes        []DataNodeConfig `json:"nodes"`
	DataBlocks   int              `json:"data_blocks"`
	ParityBlocks int              `json:"parity_blocks"`
}

//DataNodeConfig is the configuration for a datanode
type DataNodeConfig struct {
	Ip   string `json:"ip"`
	Port string `json:"port"`
}
