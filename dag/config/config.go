package config

import "time"

//PoolConfig is the configuration for the dag pool
type PoolConfig struct {
	Listen        string        `json:"listen"`
	ClusterConfig ClusterConfig `json:"cluster"`
	LeveldbPath   string        `json:"leveldb_path"`
	RootUser      string        `json:"root_user"`
	RootPassword  string        `json:"root_password"`
	GcPeriod      time.Duration `json:"gc_period"`
}

//ClusterConfig is the configuration for a cluster
type ClusterConfig struct {
	Cluster []DagNodeConfig `json:"cluster"`
}

//DagNodeConfig is the configuration for a dag node
type DagNodeConfig struct {
	Name         string           `json:"name"`
	Nodes        []DataNodeConfig `json:"nodes"`
	DataBlocks   int              `json:"data_blocks"`   // Number of data shards
	ParityBlocks int              `json:"parity_blocks"` // Number of parity shards
}

//DataNodeConfig is the configuration for a datanode
type DataNodeConfig struct {
	SetIndex   int    `json:"index"`
	RpcAddress string `json:"rpc_address"`
}

type DataNodeConfigs []DataNodeConfig

func (n DataNodeConfigs) Len() int {
	return len(n)
}

func (n DataNodeConfigs) Less(i int, j int) bool {
	return n[i].SetIndex < n[j].SetIndex
}

func (n DataNodeConfigs) Swap(i int, j int) {
	n[i], n[j] = n[j], n[i]
}
