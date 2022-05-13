package config

type PoolConfig struct {
	DagNodeConfig    []NodeConfig `json:"dag_node"`
	LeveldbPath      string       `json:"leveldb_path"`
	ImporterBatchNum int          `json:"importer_batch_num"`
}
type NodeConfig struct {
	Casks        []CaskConfig `json:"casks"`
	DataBlocks   int          `json:"data_blocks"`
	ParityBlocks int          `json:"parity_blocks"`
	LevelDbPath  string       `json:"level_db_path"`
}
type CaskConfig struct {
	Path    string `json:"path"`
	CaskNum uint32 `json:"cask_num"`
}

func DefaultConfig() *CaskConfig {
	return &CaskConfig{
		CaskNum: 256,
	}
}

type Option func(cfg *CaskConfig)

func CaskNumConf(caskNum int) Option {
	return func(cfg *CaskConfig) {
		cfg.CaskNum = uint32(caskNum)
	}
}

func PathConf(dir string) Option {
	return func(cfg *CaskConfig) {
		cfg.Path = dir
	}
}