package machine

import (
	"github.com/TheRealSibasishBehera/syncsh/internal/watcher"
	"net"
)

// Macine represents the syncsh machine configuration
type Machine struct {
	Interface   net.Interface
	HistoryFile string
	Peer        string

	//not sure but lets start the watcher here
	Watcher *watcher.Watcher
}

// newmachine
// call the new machine keys
// make a new tunnel
func NewMachine() error {
	// setup interfaces
	// setup what files to watch
	// deafult files to watch are bash_history, zsh_history
	return nil
}
