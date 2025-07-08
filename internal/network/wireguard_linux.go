//go:build linux

package network

import (
	"fmt"
	"github.com/vishvananda/netlink"
	"log/slog"
)

const (
	DefaultWireGuardMTU = 1420 // Default MTU for WireGuard interfaces
)

// we only manage 1 point to point connection
// this is the network interface for that

func GetOrCreateWireGuardInterface(name string) (netlink.Link, error) {
	link, err := netlink.LinkByName(name)
	if err == nil {
		slog.Info("Found existing WireGuard interface.", "name", name)
		if _, ok := link.(*netlink.Wireguard); !ok {
			return nil, fmt.Errorf("link %q is not a WireGuard interface", name)
		}
		return link, nil
	}
	if _, ok := err.(netlink.LinkNotFoundError); !ok {
		return nil, fmt.Errorf("find WireGuard link %q: %v", name, err)
	}
	if link, err = NewWireGuardInterface(name); err != nil {
		return nil, err
	}
	slog.Info("Created WireGuard interface.", "name", name)

	link, err = netlink.LinkByName(name)
	if err != nil {
		return nil, fmt.Errorf("find created WireGuard link %q: %v", name, err)
	}
	return link, nil
}

func NewWireGuardInterface(name string) (*netlink.Wireguard, error) {
	link := &netlink.Wireguard{
		LinkAttrs: netlink.LinkAttrs{
			Name: name,
			MTU:  DefaultWireGuardMTU,
		},
	}

	if err := netlink.LinkAdd(link); err != nil {
		return nil, fmt.Errorf("create WireGuard interface %s: %w", name, err)
	}

	return link, nil
}
