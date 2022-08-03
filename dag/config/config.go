package config

type PoolConfig struct {
	Listen        string          `json:"listen"`
	DagNodeConfig []DagNodeConfig `json:"dag_node"`
	LeveldbPath   string          `json:"leveldb_path"`
	RootUser      string          `json:"root_user"`
	RootPassword  string          `json:"root_password"`
}

//DagNodeConfig is the configuration for a dag node
type DagNodeConfig struct {
	Nodes        []DataNodeConfig `json:"nodes"`
	DataBlocks   int              `json:"data_blocks"`
	ParityBlocks int              `json:"parity_blocks"`
	LevelDbPath  string           `json:"level_db_path"`
}

//DataNodeConfig is the configuration for a datanode
type DataNodeConfig struct {
	Ip   string `json:"ip"`
	Port string `json:"port"`
}
