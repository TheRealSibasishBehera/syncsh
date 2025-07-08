package machine

import (
	"net"

	"fmt"
	"github.com/TheRealSibasishBehera/syncsh/internal/config"
	"github.com/TheRealSibasishBehera/syncsh/internal/network"
	"github.com/TheRealSibasishBehera/syncsh/internal/watcher"
	"log/slog"
)

// Macine represents the syncsh machine configuration
type Machine struct {
	Interface   net.Interface
	HistoryFile string
	Peer        string

	Watcher *watcher.Watcher
}

func NewMachine(config *config.Config) error {
	_, err := network.GetOrCreateWireGuardInterface(config.InterfaceName)
	if err != nil {
		return fmt.Errorf("failed to create WireGuard interface: %w", err)
	}

	slog.Info("Generating WireGuard keys")
	privKey, pubKey, err := network.NewMachineKeys()
	if err != nil {
		return fmt.Errorf("failed to generate WireGuard keys: %w", err)
	}

	slog.Info("Saving WireGuard keys to config")
	config.PrivateKey = string(privKey)
	config.PublicKey = string(pubKey)
	config.Save()

	return nil
}
