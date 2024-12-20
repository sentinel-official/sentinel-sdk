package v2ray

import (
	"errors"
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

// ClientOptions represents the V2Ray client configuration options.
type ClientOptions struct{}

// WriteToFile writes the ClientOptions configuration to a TOML file.
func (co *ClientOptions) WriteToFile(filepath string) error {
	data, err := toml.Marshal(co)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath, data, 0644)
}

// WriteConfigToFile writes the ClientOptions configuration to a file in JSON format.
func (co *ClientOptions) WriteConfigToFile(filepath string) error {
	data, err := co.ToConfig()
	if err != nil {
		return err
	}

	return os.WriteFile(filepath, []byte(data), 0644)
}

// InboundServerOptions represents the V2Ray inbound server configuration options.
type InboundServerOptions struct {
	Network     string `json:"network"`
	Port        uint16 `json:"port"`
	Protocol    string `json:"protocol"`
	Security    string `json:"security"`
	TLSCertPath string `json:"tls_cert_path"`
	TLSKeyPath  string `json:"tls_key_path"`
}

// Tag creates a Tag instance based on the InboundServerOptions configuration.
func (so *InboundServerOptions) Tag() *Tag {
	protocol := NewProtocolFromString(so.Protocol)
	network := NewNetworkFromString(so.Network)
	security := NewSecurityFromString(so.Security)

	return &Tag{
		p: protocol,
		n: network,
		s: security,
	}
}

// Validate validates the InboundServerOptions fields.
func (so *InboundServerOptions) Validate() error {
	network := NewNetworkFromString(so.Network)
	if !network.IsValid() {
		return fmt.Errorf("invalid network value: %s", so.Network)
	}

	protocol := NewProtocolFromString(so.Protocol)
	if !protocol.IsValid() {
		return fmt.Errorf("invalid protocol value: %s", so.Protocol)
	}

	security := NewSecurityFromString(so.Security)
	if !security.IsValid() {
		return fmt.Errorf("invalid security value: %s", so.Security)
	}

	if security == SecurityTLS {
		if so.TLSCertPath == "" || so.TLSKeyPath == "" {
			return errors.New("TLS cert path and key path cannot be empty when security is 'tls'")
		}
	}

	return nil
}

// ServerOptions represents the V2Ray server configuration options.
type ServerOptions struct {
	Inbounds []*InboundServerOptions `json:"inbounds"`
}

// Validate validates the ServerOptions fields.
func (so *ServerOptions) Validate() error {
	portSet := make(map[uint16]bool)
	tagSet := make(map[string]bool)

	for _, inbound := range so.Inbounds {
		if err := inbound.Validate(); err != nil {
			return err
		}

		if inbound.Port <= 1024 {
			return fmt.Errorf("port must be greater than 1024, got: %d", inbound.Port)
		}
		if portSet[inbound.Port] {
			return fmt.Errorf("port collision detected for port: %d", inbound.Port)
		}
		portSet[inbound.Port] = true

		tag := inbound.Tag().String()
		if tagSet[tag] {
			return fmt.Errorf("duplicate tag detected: %s", tag)
		}
		tagSet[tag] = true
	}

	return nil
}

// WriteToFile writes the ServerOptions configuration to a TOML file.
func (so *ServerOptions) WriteToFile(filepath string) error {
	data, err := toml.Marshal(so)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath, data, 0644)
}

// WriteConfigToFile writes the ServerOptions configuration to a file in JSON format.
func (so *ServerOptions) WriteConfigToFile(filepath string) error {
	data, err := so.ToConfig()
	if err != nil {
		return err
	}

	return os.WriteFile(filepath, []byte(data), 0644)
}
