package wireguard

import (
	"fmt"
	"net/netip"
	"sync"
)

// Peer represents a network peer with identity and IP addresses.
type Peer struct {
	Identity string     // Identity of the peer
	IPv4Addr netip.Addr // IPv4 address of the peer
	IPv6Addr netip.Addr // IPv6 address of the peer
}

// Key returns the identity of the peer as the key.
func (p *Peer) Key() string {
	return p.Identity
}

// PeerManager manages a collection of Peers and their associated IP addresses.
type PeerManager struct {
	*sync.RWMutex                  // Read-write mutex for thread-safe access
	IPv4Addrs     []netip.Addr     // Available IPv4 addresses
	IPv6Addrs     []netip.Addr     // Available IPv6 addresses
	m             map[string]*Peer // Map of identities to Peers
}

// NewPeerManager creates a new instance of PeerManager.
func NewPeerManager(ipv4Addrs, ipv6Addrs []netip.Addr) *PeerManager {
	return &PeerManager{
		RWMutex:   &sync.RWMutex{},
		IPv4Addrs: ipv4Addrs,
		IPv6Addrs: ipv6Addrs,
		m:         make(map[string]*Peer),
	}
}

// Get retrieves a Peer from the PeerManager by its identity.
func (pm *PeerManager) Get(v string) *Peer {
	pm.RLock()
	defer pm.RUnlock()

	return pm.m[v]
}

// Put adds a new Peer with the given identity to the PeerManager.
// It assigns available IPv4 and IPv6 addresses to the Peer.
func (pm *PeerManager) Put(v string) (ipv4Addr, ipv6Addr netip.Addr, err error) {
	pm.Lock()
	defer pm.Unlock()

	// Check if the Peer already exists
	if _, ok := pm.m[v]; ok {
		return netip.Addr{}, netip.Addr{}, fmt.Errorf("peer %s already exists", v)
	}

	// Check if there are available IP addresses
	if len(pm.IPv4Addrs) == 0 || len(pm.IPv6Addrs) == 0 {
		return netip.Addr{}, netip.Addr{}, fmt.Errorf("no available IP addresses")
	}

	// Assign the first available IPv4 and IPv6 addresses
	ipv4Addr = pm.IPv4Addrs[0]
	ipv6Addr = pm.IPv6Addrs[0]

	// Remove assigned IP addresses from the available list
	pm.IPv4Addrs = pm.IPv4Addrs[1:]
	pm.IPv6Addrs = pm.IPv6Addrs[1:]

	// Create and store the new Peer
	pm.m[v] = &Peer{
		Identity: v,
		IPv4Addr: ipv4Addr,
		IPv6Addr: ipv6Addr,
	}

	return ipv4Addr, ipv6Addr, nil
}

// Delete removes a Peer from the PeerManager by its identity.
func (pm *PeerManager) Delete(v string) {
	pm.Lock()
	defer pm.Unlock()

	// Retrieve the Peer and its IP addresses
	item, ok := pm.m[v]
	if !ok {
		return
	}

	// Add the IPv4 and IPv6 addresses back to the available list
	pm.IPv4Addrs = append(pm.IPv4Addrs, item.IPv4Addr)
	pm.IPv6Addrs = append(pm.IPv6Addrs, item.IPv6Addr)

	// Remove the Peer from the PeerManager
	delete(pm.m, v)
}

// Len returns the number of Peers in the PeerManager.
func (pm *PeerManager) Len() int {
	pm.RLock()
	defer pm.RUnlock()

	return len(pm.m)
}

// Iterate iterates over each Peer in the PeerManager and applies the provided function.
// If the function returns true, the iteration stops.
// If the function returns an error, the iteration stops and the error is returned.
func (pm *PeerManager) Iterate(fn func(key string, value *Peer) (bool, error)) error {
	pm.RLock()
	defer pm.RUnlock()

	for key, value := range pm.m {
		stop, err := fn(key, value)
		if err != nil {
			return err
		}

		if stop {
			return nil
		}
	}

	return nil
}
