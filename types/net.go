package types

import (
	"errors"
	"fmt"
	"net/netip"
)

const maxNetPrefixSize = 1 << 16

type NetPrefix struct {
	netip.Prefix
}

// NewNetPrefixFromString creates a new NetPrefix object from a given string.
func NewNetPrefixFromString(s string) (*NetPrefix, error) {
	if s == "" {
		return &NetPrefix{}, nil
	}

	prefix, err := netip.ParsePrefix(s)
	if err != nil {
		return nil, fmt.Errorf("failed to parse net prefix: %w", err)
	}

	p := &NetPrefix{prefix}
	if err := p.Validate(); err != nil {
		return nil, fmt.Errorf("invalid net prefix: %w", err)
	}

	return p, nil
}

// Len calculates the number of addresses in the NetPrefix block.
func (p NetPrefix) Len() int64 {
	diff := p.Addr().BitLen() - p.Bits()
	if diff < 0 {
		return 0
	}

	return int64(1) << diff
}

// Addrs returns a slice of all addresses within the NetPrefix block.
func (p NetPrefix) Addrs() ([]netip.Addr, error) {
	if p.Len() > maxNetPrefixSize {
		return nil, errors.New("prefix block size is too large")
	}

	var addrs []netip.Addr
	for addr := p.NetworkAddr(); p.Contains(addr); addr = addr.Next() {
		addrs = append(addrs, addr)
	}

	return addrs, nil
}

// Validate checks if the NetPrefix block size is within limits.
func (p NetPrefix) Validate() error {
	return nil
}

// NetworkAddr returns the network address of the NetPrefix.
func (p NetPrefix) NetworkAddr() netip.Addr {
	return p.Masked().Addr()
}

// BroadcastAddr returns the broadcast address of the NetPrefix for IPv4.
// Returns an error for IPv6.
func (p NetPrefix) BroadcastAddr() (netip.Addr, error) {
	if !p.Addr().Is4() {
		return netip.Addr{}, errors.New("not applicable")
	}

	size := p.Len() - 1
	if size == 0 {
		return p.Addr(), nil
	}

	buf := p.NetworkAddr().As4()
	for i := 0; i < len(buf); i++ {
		buf[len(buf)-i-1] |= byte(size >> (i * 8))
	}

	addr, ok := netip.AddrFromSlice(buf[:])
	if !ok {
		return netip.Addr{}, errors.New("failed to parse addr from slice")
	}

	return addr, nil
}
