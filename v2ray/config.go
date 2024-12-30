package v2ray

import (
	"embed"
	"errors"
	"fmt"

	"github.com/sentinel-official/sentinel-go-sdk/utils"
)

//go:embed *.tmpl
var fs embed.FS

// ClientConfig represents the V2Ray client configuration options.
type ClientConfig struct{}

func (c *ClientConfig) WriteToFile(name string) error {
	text, err := fs.ReadFile("client.json.tmpl")
	if err != nil {
		return err
	}

	return utils.ExecTemplateToFile(string(text), c, name)
}

// InboundServerConfig represents the V2Ray inbound server configuration options.
type InboundServerConfig struct {
	Network     string `mapstructure:"network"`
	Port        uint16 `mapstructure:"port"`
	Protocol    string `mapstructure:"protocol"`
	Security    string `mapstructure:"security"`
	TLSCertPath string `mapstructure:"tls_cert_path"`
	TLSKeyPath  string `mapstructure:"tls_key_path"`
}

// Tag creates a Tag instance based on the InboundServerConfig configuration.
func (c *InboundServerConfig) Tag() *Tag {
	protocol := NewProtocolFromString(c.Protocol)
	network := NewNetworkFromString(c.Network)
	security := NewSecurityFromString(c.Security)

	return &Tag{
		p: protocol,
		n: network,
		s: security,
	}
}

// Validate validates the InboundServerConfig fields.
func (c *InboundServerConfig) Validate() error {
	network := NewNetworkFromString(c.Network)
	if !network.IsValid() {
		return fmt.Errorf("invalid network %s", c.Network)
	}

	protocol := NewProtocolFromString(c.Protocol)
	if !protocol.IsValid() {
		return fmt.Errorf("invalid protocol %s", c.Protocol)
	}

	security := NewSecurityFromString(c.Security)
	if !security.IsValid() {
		return fmt.Errorf("invalid security %s", c.Security)
	}

	if security == SecurityTLS {
		if c.TLSCertPath == "" || c.TLSKeyPath == "" {
			return errors.New("TLS cert path and key path cannot be empty")
		}
	}

	return nil
}

// ServerConfig represents the V2Ray server configuration options.
type ServerConfig struct {
	Inbounds []InboundServerConfig `mapstructure:"inbounds"`
}

// Validate validates the ServerConfig fields.
func (c *ServerConfig) Validate() error {
	if len(c.Inbounds) == 0 {
		return errors.New("inbounds cannot be empty")
	}

	portSet := make(map[uint16]bool)
	tagSet := make(map[string]bool)

	for _, inbound := range c.Inbounds {
		if err := inbound.Validate(); err != nil {
			return err
		}

		if inbound.Port <= 1024 {
			return errors.New("port must be greater than 1024")
		}
		if portSet[inbound.Port] {
			return fmt.Errorf("duplicate port %d", inbound.Port)
		}
		portSet[inbound.Port] = true

		tag := inbound.Tag().String()
		if tagSet[tag] {
			return fmt.Errorf("duplicate tag %s", tag)
		}
		tagSet[tag] = true
	}

	return nil
}

func (c *ServerConfig) WriteToFile(name string) error {
	text, err := fs.ReadFile("server.json.tmpl")
	if err != nil {
		return err
	}

	return utils.ExecTemplateToFile(string(text), c, name)
}

func DefaultServerConfig() ServerConfig {
	return ServerConfig{
		Inbounds: []InboundServerConfig{
			{
				Network:     "grpc",
				Port:        utils.RandomPort(),
				Protocol:    "vmess",
				Security:    "none",
				TLSCertPath: "",
				TLSKeyPath:  "",
			},
			{
				Network:     "tcp",
				Port:        utils.RandomPort(),
				Protocol:    "vmess",
				Security:    "none",
				TLSCertPath: "",
				TLSKeyPath:  "",
			},
		},
	}
}
