package config

type ConfigOption func(*Config)

func WithHistoryPath(path string) ConfigOption {
	return func(c *Config) {
		c.HistoryPath = path
	}
}

func WithInterfaceName(name string) ConfigOption {
	return func(c *Config) {
		c.InterfaceName = name
	}
}

func WithSQLitePath(path string) ConfigOption {
	return func(c *Config) {
		c.SQLitePath = path
	}
}
