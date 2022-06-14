package config

//PoolConfig is the configuration for the dag pool
type PoolConfig struct {
	Listen        string          `json:"listen"`
	DagNodeConfig []DagNodeConfig `json:"dag_node"`
	LeveldbPath   string          `json:"leveldb_path"`
	RootUser      string          `json:"root_user"`
	RootPassword  string          `json:"root_password"`
	GcPeriod      string          `json:"gc_period"`
}

//DagNodeConfig is the configuration for a dag node
type DagNodeConfig struct {
	Nodes        []dataNodeConfig `json:"nodes"`
	DataBlocks   int              `json:"data_blocks"`
	ParityBlocks int              `json:"parity_blocks"`
	LevelDbPath  string           `json:"level_db_path"`
}

type dataNodeConfig struct {
	Ip        string `json:"ip"`
	Port      string `json:"port"`
	HeartAddr string `json:"heart_addr"`
}
