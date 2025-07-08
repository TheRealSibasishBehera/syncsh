package network

import (
	"fmt"
	"net/netip"

	"github.com/TheRealSibasishBehera/syncsh/internal/network/tunnel"
	"github.com/TheRealSibasishBehera/syncsh/internal/secret"
)

// CreateP2PConnection creates a simple point-to-point WireGuard connection
func CreateP2PConnection(isInitiator bool, remoteEndpoint string, localPrivKey, remotePubKey secret.Secret) (*tunnel.Tunnel, error) {
	var localIP, remoteIP netip.Addr

	// Assign fixed IPs for P2P
	if isInitiator {
		localIP = netip.MustParseAddr(MachineAIP)
		remoteIP = netip.MustParseAddr(MachineBIP)
	} else {
		localIP = netip.MustParseAddr(MachineBIP)
		remoteIP = netip.MustParseAddr(MachineAIP)
	}

	endpoint, err := netip.ParseAddrPort(remoteEndpoint)
	if err != nil {
		return nil, fmt.Errorf("parse remote endpoint: %w", err)
	}

	config := &tunnel.Config{
		LocalAddress:    localIP,
		LocalPrivateKey: localPrivKey,
		Endpoint:        endpoint,
		RemotePublicKey: remotePubKey,
		RemoteNetwork:   netip.PrefixFrom(remoteIP, 32), // Single host
	}

	return tunnel.Connect(config)
}

// CreateP2PConfig creates a P2PConfig for the given parameters
func CreateP2PConfig(isInitiator bool, remoteEndpoint string, localPrivKey, remotePubKey secret.Secret) (P2PConfig, error) {
	var localIP, remoteIP netip.Addr

	if isInitiator {
		localIP = netip.MustParseAddr(MachineAIP)
		remoteIP = netip.MustParseAddr(MachineBIP)
	} else {
		localIP = netip.MustParseAddr(MachineBIP)
		remoteIP = netip.MustParseAddr(MachineAIP)
	}

	endpoint, err := netip.ParseAddrPort(remoteEndpoint)
	if err != nil {
		return P2PConfig{}, fmt.Errorf("parse remote endpoint: %w", err)
	}

	return P2PConfig{
		LocalIP:         localIP,
		RemoteIP:        remoteIP,
		LocalPrivateKey: localPrivKey,
		RemotePublicKey: remotePubKey,
		RemoteEndpoint:  endpoint,
	}, nil
}
