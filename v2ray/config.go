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
		return fmt.Errorf("failed to read template: %w", err)
	}

	if err := utils.ExecTemplateToFile(string(text), c, name); err != nil {
		return fmt.Errorf("failed to execute template to file: %w", err)
	}

	return nil
}

// InboundServerConfig represents the V2Ray inbound server configuration options.
type InboundServerConfig struct {
	Port        string `mapstructure:"port"`
	Proxy       string `mapstructure:"proxy"`
	Security    string `mapstructure:"security"`
	TLSCertPath string `mapstructure:"tls_cert_path"`
	TLSKeyPath  string `mapstructure:"tls_key_path"`
	Transport   string `mapstructure:"transport"`
}

func (c *InboundServerConfig) InPort() string {
	v, err := types.NewPortFromString(c.Port)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%d-%d", v.InFrom, v.InTo)
}

func (c *InboundServerConfig) OutPort() string {
	v, err := types.NewPortFromString(c.Port)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%d-%d", v.OutFrom, v.OutTo)
}

// Tag creates a Tag instance based on the InboundServerConfig configuration.
func (c *InboundServerConfig) Tag() *Tag {
	proxy := NewProxyProtocolFromString(c.Proxy)
	security := NewTransportSecurityFromString(c.Security)
	transport := NewTransportProtocolFromString(c.Transport)

	return &Tag{
		Proxy:     proxy,
		Security:  security,
		Transport: transport,
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
	if v := NewProxyProtocolFromString(c.Proxy); !v.IsValid() {
		return fmt.Errorf("invalid proxy %s", v)
	}

	security := NewTransportSecurityFromString(c.Security)
	if !security.IsValid() {
		return fmt.Errorf("invalid security %s", security)
	}
	if security == TransportSecurityTLS {
		if c.TLSCertPath == "" {
			return errors.New("tls_cert_path cannot be empty")
		}
		if c.TLSKeyPath == "" {
			return errors.New("tls_key_path cannot be empty")
		}
	}

	if v := NewTransportProtocolFromString(c.Transport); !v.IsValid() {
		return fmt.Errorf("invalid transport %s", v)
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
		return fmt.Errorf("failed to read template: %w", err)
	}

	if err := utils.ExecTemplateToFile(string(text), c, name); err != nil {
		return fmt.Errorf("failed to execute template to file: %w", err)
	}

	return nil
}

func DefaultServerConfig() ServerConfig {
	return ServerConfig{
		Inbounds: []InboundServerConfig{
			{
				Port:        fmt.Sprintf("%d", utils.RandomPort()),
				Proxy:       "vmess",
				Security:    "none",
				TLSCertPath: "",
				TLSKeyPath:  "",
				Transport:   "grpc",
			},
			{
				Port:        fmt.Sprintf("%d", utils.RandomPort()),
				Proxy:       "vmess",
				Security:    "none",
				TLSCertPath: "",
				TLSKeyPath:  "",
				Transport:   "tcp",
			},
		},
	}
}
