package network

import (
	"time"
)

// Simple peer status for P2P connections
const (
	PeerStatusConnected    = "connected"
	PeerStatusDisconnected = "disconnected"
)

// SimplePeerStatus represents basic peer connection state
type SimplePeerStatus struct {
	IP               string
	PublicKey        string
	Status           string
	LastHandshake    time.Time
	BytesReceived    int64
	BytesTransmitted int64
}

// IsConnected returns true if peer had a handshake in the last 3 minutes
func (s SimplePeerStatus) IsConnected() bool {
	return time.Since(s.LastHandshake) < 3*time.Minute
}
