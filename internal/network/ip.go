package network

import (
	"fmt"
	"github.com/vishvananda/netlink"
	"net"
	"net/netip"
)

const (
	// Fixed IPs for point-to-point
	MachineAIP  = "10.100.0.1"
	MachineBIP  = "10.100.0.2"
	NetworkCIDR = "10.100.0.0/30"
)

// since we are not using any distributed database
// we have only one point-to-point connection

type Peer struct {
	TunnelIP netip.Addr
}

// its clear that init shout now create or assign IP ,
// connect should
// so when you try to connect to a peer
// it should first assign itself an IP
// mark it and say it in the connection , its already assigned , use the next one

func GetIP(initiator bool) (netip.Addr, error) {
	if initiator {
		return netip.MustParseAddr(MachineAIP), nil
	}
	return netip.MustParseAddr(MachineBIP), nil
}

func ReserveIp(addr netip.Addr, link string) error {
	devLink, err := netlink.LinkByName(link)
	if err != nil {
		return fmt.Errorf("find link %q: %w", link, err)
	}
	addrNetlink := &netlink.Addr{}
	netlink.AddrAdd(devLink, addrNetlink)
	return nil
}

func addrToSingleIPPrefix(addr netip.Addr) (netip.Prefix, error) {
	if !addr.IsValid() {
		return netip.Prefix{}, fmt.Errorf("invalid IP address")
	}
	bits := 32
	if addr.Is6() {
		bits = 128
	}
	return addr.Prefix(bits)
}

func prefixToIPNet(prefix netip.Prefix) net.IPNet {
	return net.IPNet{
		IP:   prefix.Addr().AsSlice(),
		Mask: net.CIDRMask(prefix.Bits(), prefix.Addr().BitLen()),
	}
}
