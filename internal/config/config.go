package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/TheRealSibasishBehera/syncsh/pkg/utils"
	"github.com/goccy/go-yaml"
)

type ShellKind string

const (
	ShellBash ShellKind = "bash"
	ShellZsh  ShellKind = "zsh"
	ShellFish ShellKind = "fish"
)

// GetDefaultHistoryPath returns the default history path for each shell
func (s ShellKind) GetDefaultHistoryPath() string {
	home := os.Getenv("HOME")
	if home == "" {
		return ""
	}

	switch s {
	case ShellBash:
		return filepath.Join(home, ".bash_history")
	case ShellZsh:
		if histFile := os.Getenv("HISTFILE"); histFile != "" {
			return histFile
		}
		return filepath.Join(home, ".zsh_history")
	case ShellFish:
		return filepath.Join(home, ".local", "share", "fish", "fish_history")
	default:
		return ""
	}
}

type Config struct {
	SQLitePath string    `yaml:"sql_path"` // path to the SQLite database file
	kind       ShellKind `yaml:"shell"`    // kind of shell to use, e.g., ["bash", "zsh"]

	// WireGuard keys
	PrivateKey string `yaml:"private_key"` // WireGuard private key (base64)
	PublicKey  string `yaml:"public_key"`  // WireGuard public key (base64)

	//optional path
	HistoryPath   string `yaml:"history"`   // path for the history file
	InterfaceName string `yaml:"interface"` // name for wireguard interface

	//internal
	path string // config is read from this path
}

func NewConfigWithOpts(opts ...ConfigOption) (*Config, error) {
	cfg := &Config{}

	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.SQLitePath == "" {
		cfg.SQLitePath = "store/syncsh.db"
	}
	if cfg.InterfaceName == "" {
		cfg.InterfaceName = "syncsh0"
	}
	if cfg.kind == "" {
		shellKind, err := utils.GetShellKind()
		if err != nil {
			return nil, err
		}
		cfg.kind = ShellKind(shellKind)
	}
	if cfg.HistoryPath == "" {
		historyPath, err := utils.GetDefaultHistoryPath(string(cfg.kind))
		if err != nil {
			return nil, err
		}
		cfg.HistoryPath = historyPath
	}

	return cfg, nil
}

func (c *Config) Path() string {
	return c.path
}

// SetPath sets the config file path
func (c *Config) SetPath(path string) {
	c.path = path
}

// GetResolvedHistoryPath returns the history path, using default if not set
func (c *Config) GetResolvedHistoryPath() string {
	if c.HistoryPath != "" {
		return c.HistoryPath
	}
	return c.kind.GetDefaultHistoryPath()
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
