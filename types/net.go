package types

import (
	"errors"
	"net/netip"
)

const MaxNetPrefixSize = 256

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
		return nil, err
	}

	p := &NetPrefix{prefix}
	if err := p.Validate(); err != nil {
		return nil, err
	}

	return p, nil
}

// Len calculates the number of addresses in the NetPrefix block.
func (p NetPrefix) Len() int64 {
	bitDiff := p.Addr().BitLen() - p.Bits()
	if bitDiff < 0 {
		return 0
	}

	return int64(1) << bitDiff
}

// Addrs returns a slice of all addresses within the NetPrefix block.
func (p NetPrefix) Addrs() ([]netip.Addr, error) {
	if p.Len() > MaxNetPrefixSize {
		return nil, errors.New("prefix block size is too large")
	}

	var addrs []netip.Addr
	for addr := p.Addr(); p.Contains(addr); addr = addr.Next() {
		addrs = append(addrs, addr)
	}

	return addrs, nil
}

func (p NetPrefix) Validate() error {
	if p.Len() > MaxNetPrefixSize {
		return errors.New("prefix block size is too large")
	}

	return nil
}
