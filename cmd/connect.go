package cmd

import (
	// "github.com/TheRealSibasishBehera/syncsh/internal/machine"
	"github.com/spf13/cobra"
)

func NewConnectCommand() *cobra.Command {
	connectCmd := &cobra.Command{
		Use:   "connect",
		Short: "Connect to a remote syncsh machine",
		Long:  `This command connects to a remote syncsh machine and starts synchronizing shell sessions.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Call the function to connect to the remote machine
			// if err := machine.ConnectToRemoteMachine(); err != nil {
			// 	return err
			// }
			return nil
		},
	}

	return connectCmd
}
