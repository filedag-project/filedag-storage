package diskv

type Config struct {
	Dir            string
	MaxLinkDagSize int
	Shard          ShardFun
	Parallel       int
	MaxCacheDags   int
}

type Option func(cfg *Config)

func DirConf(dir string) Option {
	return func(cfg *Config) {
		cfg.Dir = dir
	}
}

func MaxLinkDagSizeConf(size int) Option {
	return func(cfg *Config) {
		cfg.MaxLinkDagSize = size
	}
}

func ParallelConf(n int) Option {
	return func(cfg *Config) {
		cfg.Parallel = n
	}
}

func MaxCacheDagsConf(n int) Option {
	return func(cfg *Config) {
		cfg.MaxCacheDags = n
	}
}

func ShardFunConf(sha ShardFun) Option {
	return func(cfg *Config) {
		cfg.Shard = sha
	}
}
