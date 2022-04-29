package node

type Config struct {
	Path    string
	CaskNum uint32
}

func defaultConfig() *Config {
	return &Config{
		CaskNum: 256,
	}
}

type Option func(cfg *Config)

func CaskNumConf(caskNum int) Option {
	return func(cfg *Config) {
		cfg.CaskNum = uint32(caskNum)
	}
}

func PathConf(dir string) Option {
	return func(cfg *Config) {
		cfg.Path = dir
	}
}
