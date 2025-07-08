//go:build !linux

package network

import (
	"errors"
	"github.com/vishvananda/netlink"
)

// GetOrCreateWireGuardInterface is a stub for non-Linux systems
func GetOrCreateWireGuardInterface(name string) (netlink.Link, error) {
	return nil, errors.New("WireGuard interface creation is only supported on Linux")
}