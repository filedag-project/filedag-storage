package mutcask

type Config struct {
	CaskNum uint32
}

func defaultConfig() *Config {
	return &Config{
		CaskNum: 256,
	}
}

type Option func(cfg *Config)

func ConfCaskNum(caskNum int) Option {
	return func(cfg *Config) {
		cfg.CaskNum = uint32(caskNum)
	}
}
