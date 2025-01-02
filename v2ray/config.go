package v2ray

import (
	"embed"
	"errors"
	"fmt"

	"github.com/sentinel-official/sentinel-go-sdk/types"
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
	Port        string `mapstructure:"port"`
	Protocol    string `mapstructure:"protocol"`
	Security    string `mapstructure:"security"`
	TLSCertPath string `mapstructure:"tls_cert_path"`
	TLSKeyPath  string `mapstructure:"tls_key_path"`
}

func (c *InboundServerConfig) ListenPort() string {
	v, err := types.NewPortFromString(c.Port)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%d-%d", v.InFrom, v.InTo)
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
	if c.Port == "" {
		return errors.New("port cannot be empty")
	}
	if _, err := types.NewPortFromString(c.Port); err != nil {
		return fmt.Errorf("invalid port: %w", err)
	}
	if v := NewNetworkFromString(c.Network); !v.IsValid() {
		return fmt.Errorf("invalid network %s", v)
	}
	if v := NewProtocolFromString(c.Protocol); !v.IsValid() {
		return fmt.Errorf("invalid protocol %s", v)
	}

	security := NewSecurityFromString(c.Security)
	if !security.IsValid() {
		return fmt.Errorf("invalid security %s", security)
	}

	if security == SecurityTLS {
		if c.TLSCertPath == "" {
			return errors.New("tls_cert_path cannot be empty")
		}
		if c.TLSKeyPath == "" {
			return errors.New("tls_key_path cannot be empty")
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

	inPortSet := make(map[uint16]bool)
	outPortSet := make(map[uint16]bool)
	tagSet := make(map[string]bool)

	for _, inbound := range c.Inbounds {
		if err := inbound.Validate(); err != nil {
			return fmt.Errorf("invalid inbound: %w", err)
		}

		port, err := types.NewPortFromString(inbound.Port)
		if err != nil {
			panic(err)
		}

		for p := port.InFrom; p <= port.InTo; p++ {
			if inPortSet[p] {
				return fmt.Errorf("duplicate in port %d", p)
			}
			inPortSet[p] = true
		}
		for p := port.OutFrom; p <= port.OutTo; p++ {
			if outPortSet[p] {
				return fmt.Errorf("duplicate out port %d", p)
			}
			outPortSet[p] = true
		}

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
				Port:        fmt.Sprintf("%d", utils.RandomPort()),
				Protocol:    "vmess",
				Security:    "none",
				TLSCertPath: "",
				TLSKeyPath:  "",
			},
			{
				Network:     "tcp",
				Port:        fmt.Sprintf("%d", utils.RandomPort()),
				Protocol:    "vmess",
				Security:    "none",
				TLSCertPath: "",
				TLSKeyPath:  "",
			},
		},
	}
}
