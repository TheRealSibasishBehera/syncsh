package cmd

import (
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "syncsh",
		Short: "syncsh is a tool for synchronizing shells across networks",
		Long:  `syncsh is a powerful command-line tool designed to synchronize shell sessions across multiple machines in a network.`,
	}

	rootCmd.AddCommand(
		NewInitCommand(),
		NewConnectCommand(),
	)

	return rootCmd
}
