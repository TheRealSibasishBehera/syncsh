package network

import (
	"fmt"
	"net"
	"net/netip"

	"github.com/TheRealSibasishBehera/syncsh/internal/secret"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// Simple P2P WireGuard configuration
type P2PConfig struct {
	LocalIP         netip.Addr
	RemoteIP        netip.Addr
	LocalPrivateKey secret.Secret
	RemotePublicKey secret.Secret
	RemoteEndpoint  netip.AddrPort
}

// SimplePeer represents a single peer in P2P connection
type SimplePeer struct {
	IP        netip.Addr
	PublicKey secret.Secret
	Endpoint  netip.AddrPort
}

// IsConfigured returns true if the P2P configuration is complete
func (c P2PConfig) IsConfigured() bool {
	return c.LocalIP.IsValid() && c.RemoteIP.IsValid() &&
		c.LocalPrivateKey != nil && c.RemotePublicKey != nil &&
		c.RemoteEndpoint.IsValid()
}

// CreateP2PDeviceConfig creates a simple WireGuard device config for P2P connection
func (c P2PConfig) CreateP2PDeviceConfig() (wgtypes.Config, error) {
	privateKey, err := wgtypes.NewKey(c.LocalPrivateKey)
	if err != nil {
		return wgtypes.Config{}, fmt.Errorf("parse private key: %w", err)
	}
	
	publicKey, err := wgtypes.NewKey(c.RemotePublicKey)
	if err != nil {
		return wgtypes.Config{}, fmt.Errorf("parse remote public key: %w", err)
	}
	
	listenPort := WireGuardPort
	keepalive := WireGuardKeepaliveInterval
	
	// Simple P2P peer configuration
	peerConfig := wgtypes.PeerConfig{
		PublicKey: publicKey,
		Endpoint: &net.UDPAddr{
			IP:   c.RemoteEndpoint.Addr().AsSlice(),
			Port: int(c.RemoteEndpoint.Port()),
		},
		AllowedIPs: []net.IPNet{
			{
				IP:   c.RemoteIP.AsSlice(),
				Mask: net.CIDRMask(32, 32), // Single host
			},
		},
		PersistentKeepaliveInterval: &keepalive,
		ReplaceAllowedIPs:           true,
	}
	
	return wgtypes.Config{
		PrivateKey:   &privateKey,
		ListenPort:   &listenPort,
		ReplacePeers: true, // Simple replacement for P2P
		Peers:        []wgtypes.PeerConfig{peerConfig},
	}, nil
}
