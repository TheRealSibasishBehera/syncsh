package cmd

import (
	// "github.com/TheRealSibasishBehera/syncsh/internal/machine"
	"github.com/spf13/cobra"
)

func NewInitCommand() *cobra.Command {
	var historyPath string
	var interfaceName string

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize syncsh on this machine",
		Long:  `This command sets up the necessary configuration and keys for syncsh on this machine.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// create a new config
			// store the config
			// create a machine
			return nil
		},
	}

	initCmd.Flags().StringVar(&historyPath, "history-path", "", "Custom path to shell history file (default: auto-detect)")
	initCmd.Flags().StringVar(&interfaceName, "interface", "syncsh0", "WireGuard interface name")

	return initCmd
}
