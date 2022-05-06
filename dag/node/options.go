package node

type CaskConfig struct {
	Path    string
	CaskNum uint32
}

func defaultConfig() *CaskConfig {
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
