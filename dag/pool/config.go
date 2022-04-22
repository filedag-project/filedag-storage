package pool

type SimplePoolConfig struct {
	StorePath string // the path for kv db
	BatchNum  int    // use for blockstore batch task and importer batch task config
	CaskNum   int    // the number of vlog files in mutcask kv
}
