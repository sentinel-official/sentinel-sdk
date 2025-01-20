package types

import (
	"errors"
	"fmt"
	"net/netip"
	"sync"
)

// IPPool manages a pool of IP addresses, including assigned, reserved, and unassigned addresses.
// It ensures thread-safe operations and manages address allocation and deallocation.
type IPPool struct {
	assigned   map[netip.Addr]bool // Tracks IPs that are currently assigned.
	reserved   map[netip.Addr]bool // Tracks IPs that are reserved.
	unassigned []netip.Addr        // List of unassigned IPs available for allocation.

	addr   netip.Addr // Current IP address in the pool.
	prefix *NetPrefix // The network prefix associated with the pool.

	m *sync.Mutex // Mutex to ensure thread-safe access to the pool.
}

// NewIPPoolFromString creates a new IPPool using a given network prefix string.
// It reserves the network address and, if applicable, the broadcast address for the prefix.
func NewIPPoolFromString(s string) (*IPPool, error) {
	prefix, err := NewNetPrefixFromString(s)
	if err != nil {
		return nil, fmt.Errorf("failed to get net prefix: %w", err)
	}

	p := &IPPool{
		assigned:   make(map[netip.Addr]bool),
		reserved:   make(map[netip.Addr]bool),
		unassigned: []netip.Addr{},
		addr:       prefix.NetworkAddr(),
		prefix:     prefix,
		m:          &sync.Mutex{},
	}

	// Reserve the network and prefix address.
	_ = p.Reserve(prefix.Addr())
	_ = p.Reserve(prefix.NetworkAddr())

	// For IPv4, reserve the broadcast address if it exists.
	if p.addr.Is4() {
		broadcast, err := prefix.BroadcastAddr()
		if err != nil {
			return nil, fmt.Errorf("failed to get broadcast addr: %w", err)
		}

		_ = p.Reserve(broadcast)
	}

	return p, nil
}

// Reserve marks an IP address as reserved, ensuring it cannot be assigned.
// Returns an error if the address is outside the prefix or already assigned/reserved.
func (p *IPPool) Reserve(addr netip.Addr) error {
	p.m.Lock()
	defer p.m.Unlock()

	if !p.prefix.Contains(addr) {
		return errors.New("addr is outside of prefix")
	}
	if p.assigned[addr] || p.reserved[addr] {
		return errors.New("addr is already assigned or reserved")
	}

	p.reserved[addr] = true
	return nil
}

// Get fetches an available IP address from the pool.
// If there are no unassigned addresses, it increments the current address until one is found.
func (p *IPPool) Get() (addr netip.Addr, err error) {
	p.m.Lock()
	defer p.m.Unlock()

	// Check if there are preloaded unassigned addresses.
	if len(p.unassigned) > 0 {
		addr, p.unassigned = p.unassigned[0], p.unassigned[1:]
	} else {
		// Increment through addresses within the prefix until an available one is found.
		for {
			if !p.prefix.Contains(p.addr) {
				return netip.Addr{}, errors.New("pool is empty")
			}

			addr, p.addr = p.addr, p.addr.Next()
			if !p.reserved[addr] {
				break
			}
		}
	}

	p.assigned[addr] = true
	return addr, nil
}

// Put returns an IP address to the pool, making it available for future allocations.
// Returns an error if the address was not previously assigned.
func (p *IPPool) Put(addr netip.Addr) error {
	p.m.Lock()
	defer p.m.Unlock()

	if !p.assigned[addr] {
		return errors.New("addr is not assigned")
	}

	// Remove from assigned list and add back to unassigned.
	delete(p.assigned, addr)
	p.unassigned = append(p.unassigned, addr)

	return nil
}
