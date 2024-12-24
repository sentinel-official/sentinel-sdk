package types

import (
	"errors"
	"net/netip"
)

type CIDR struct {
	netip.Prefix
}

// NewCIDR creates a new CIDR object from a given CIDR string.
func NewCIDR(s string) (*CIDR, error) {
	prefix, err := netip.ParsePrefix(s)
	if err != nil {
		return nil, err
	}

	return &CIDR{prefix}, nil
}

// Len calculates the number of addresses in the CIDR block.
func (p CIDR) Len() int64 {
	bitDiff := p.Addr().BitLen() - p.Bits()
	if bitDiff < 0 {
		return 0
	}

	return int64(1) << bitDiff
}

// Adds returns a slice of all addresses within the CIDR block.
func (p CIDR) Adds() ([]netip.Addr, error) {
	if p.Len() > 256 {
		return nil, errors.New("CIDR block is too large to enumerate addresses")
	}

	var addrs []netip.Addr
	for addr := p.Addr(); p.Contains(addr); addr = addr.Next() {
		addrs = append(addrs, addr)
	}

	return addrs, nil
}
