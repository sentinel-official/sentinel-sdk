package wireguard

import (
	"net/netip"
)

// AddPeerResponse represents the response for adding a peer to the WireGuard server.
type AddPeerResponse struct {
	IPv4Addr netip.Addr      `json:"ipv4_addr"` // Assigned IPv4 address for the peer.
	IPv6Addr netip.Addr      `json:"ipv6_addr"` // Assigned IPv6 address for the peer.
	Metadata *ServerMetadata `json:"metadata"`  // Metadata about the server.
}
