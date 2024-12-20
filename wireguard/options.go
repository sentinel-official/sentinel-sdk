package wireguard

import (
	"errors"
	"fmt"
	"net/netip"
	"os"

	"github.com/pelletier/go-toml/v2"
)

// ClientOptions represents the WireGuard client configuration options.
type ClientOptions struct{}

// WriteToFile writes the ClientOptions configuration to a TOML file.
func (co *ClientOptions) WriteToFile(filepath string) error {
	data, err := toml.Marshal(co)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath, data, 0644)
}

// WriteConfigToFile writes the WireGuard configuration to a file in a format recognized by WireGuard.
func (co *ClientOptions) WriteConfigToFile(filepath string) error {
	data, err := co.ToConfig()
	if err != nil {
		return err
	}

	return os.WriteFile(filepath, []byte(data), 0644)
}

// ServerOptions represents the WireGuard server configuration options.
type ServerOptions struct {
	Addresses    []string `json:"addresses"`
	EnableIPv4   bool     `json:"enable_ipv4"`
	EnableIPv6   bool     `json:"enable_ipv6"`
	Interface    string   `json:"interface"`
	ListenPort   uint16   `json:"listen_port"`
	OutInterface string   `json:"out_interface"`
	PrivateKey   string   `json:"private_key"`
}

// WriteToFile writes the ServerOptions configuration to a TOML file.
func (so *ServerOptions) WriteToFile(filepath string) error {
	data, err := toml.Marshal(so)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath, data, 0644)
}

// WriteConfigToFile writes the WireGuard server configuration to a file in a format recognized by WireGuard.
func (so *ServerOptions) WriteConfigToFile(filepath string) error {
	data, err := so.ToConfig()
	if err != nil {
		return err
	}

	return os.WriteFile(filepath, []byte(data), 0644)
}

// Validate checks that the ServerOptions fields have valid values.
func (so *ServerOptions) Validate() error {
	if len(so.Addresses) == 0 {
		return errors.New("addresses cannot be empty")
	}
	for _, item := range so.Addresses {
		if _, err := netip.ParsePrefix(item); err != nil {
			return fmt.Errorf("invalid address: %w", err)
		}
	}

	if so.Interface == "" {
		return errors.New("interface cannot be empty")
	}
	if so.ListenPort == 0 {
		return errors.New("listen_port cannot be zero")
	}
	if so.OutInterface == "" {
		return errors.New("out_interface cannot be empty")
	}

	if _, err := NewKeyFromString(so.PrivateKey); err != nil {
		return fmt.Errorf("invalid private_key: %w", err)
	}

	return nil
}
