package main

import (
	command "github.com/TheRealSibasishBehera/syncsh/cmd"
	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use:   "syncsh",
		Short: "syncsh is a tool for synchronizing shells across networks",
		Long: `
			syncsh is a powerful command-line tool designed to synchronize shell sessions across multiple machines in a network.
	`,
	}
	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {

		return nil
	}
	cmd.AddCommand(
		command.NewRootCommand(),
	)
	cobra.CheckErr(cmd.Execute())
}
