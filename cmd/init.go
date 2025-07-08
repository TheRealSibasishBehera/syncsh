package cmd

import (
	"fmt"
	"github.com/TheRealSibasishBehera/syncsh/internal/config"
	"github.com/TheRealSibasishBehera/syncsh/internal/machine"
	"github.com/TheRealSibasishBehera/syncsh/pkg/utils"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

func NewInitCommand() *cobra.Command {
	var historyPath string
	var interfaceName string
	var shellKind string

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize syncsh on this machine",
		Long:  `This command sets up the necessary configuration and keys for syncsh on this machine.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			utils.ShowInitBanner()

			if err := config.ValidateShellKindAndHistoryPath(shellKind, historyPath); err != nil {
				return fmt.Errorf("validation failed: %w", err)
			}

			shell := config.ShellKind(shellKind)

			configDir := filepath.Join(os.Getenv("HOME"), ".config", "syncsh")
			if err := os.MkdirAll(configDir, 0700); err != nil {
				return fmt.Errorf("failed to create config directory: %w", err)
			}

			var configOpts []config.ConfigOption

			if shell != "" {
				configOpts = append(configOpts, config.WithShellKind(shell))
			}
			if interfaceName != "" {
				configOpts = append(configOpts, config.WithInterfaceName(interfaceName))
			}
			if historyPath != "" {
				configOpts = append(configOpts, config.WithHistoryPath(historyPath))
			}

			dbPath := filepath.Join(configDir, "syncsh.db")
			configOpts = append(configOpts, config.WithSQLitePath(dbPath))

			cfg, err := config.NewConfigWithOpts(configOpts...)
			if err != nil {
				return fmt.Errorf("failed to create configuration: %w", err)
			}
			configPath := filepath.Join(configDir, "config.yaml")
			cfg.SetPath(configPath)
			if err := cfg.Save(); err != nil {
				return fmt.Errorf("failed to save configuration: %w", err)
			}

			err = machine.NewMachine(cfg)
			if err != nil {
				return fmt.Errorf("failed to create machine: %w", err)
			}

			return nil
		},
	}

	initCmd.Flags().StringVar(&shellKind, "shell", "zsh", "Shell type (bash, zsh, fish)")
	initCmd.Flags().StringVar(&historyPath, "history-path", "", "Custom path to shell history file (default: auto-detect based on shell)")
	initCmd.Flags().StringVar(&interfaceName, "interface", "syncsh0", "WireGuard interface name")

	return initCmd
}
