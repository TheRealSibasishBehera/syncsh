package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

type Config struct {
	SQLitePath string `yaml:"sql_path"` // path to the SQLite database file

	//optional path
	HistoryPath   string `yaml:"history"`   // path for the history file
	InterfaceName string `yaml:"interface"` // name for wireguard interface

	//internal
	path string // config is read from this path
}

func NewConfig(opts ...ConfigOption) *Config {
	//set default values first
	cfg := &Config{
		SQLitePath:    "syncsh.db",
		HistoryPath:   "",
		InterfaceName: "syncsh0",
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}

func (c *Config) Path() string {
	return c.path
}

func NewFromFile(path string) (*Config, error) {
	_, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("check file permissions '%s': %w", path, err)
	}
	c := &Config{
		path: path,
	}
	if os.IsNotExist(err) {
		return c, nil
	}

	if err = c.Read(); err != nil {
		return nil, err
	}
	return c, nil
}

// Read reads the configuration from the file specified in c.path.
func (c *Config) Read() error {
	data, err := os.ReadFile(c.path)
	if err != nil {
		return fmt.Errorf("read config file '%s': %w", c.path, err)
	}
	if err = yaml.Unmarshal(data, c); err != nil {
		return fmt.Errorf("parse config file '%s': %s", c.path, yaml.FormatError(err, true, true))
	}

	return nil
}

// Save writes the configuration to the file specified in c.path.
func (c *Config) Save() error {
	dir, _ := filepath.Split(c.path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("create config directory '%s': %w", dir, err)
	}

	f, err := os.OpenFile(c.path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("write config file '%s': %w", c.path, err)
	}

	encoder := yaml.NewEncoder(f, yaml.Indent(2), yaml.IndentSequence(true))
	if err = encoder.Encode(c); err != nil {
		_ = f.Close()
		return fmt.Errorf("encode config file '%s': %w", c.path, err)
	}
	return f.Close()
}
