package wireguard

import (
	"embed"
	"errors"
	"fmt"
	"strings"

	"github.com/sentinel-official/sentinel-go-sdk/types"
	"github.com/sentinel-official/sentinel-go-sdk/utils"
)

//go:embed *.tmpl
var fs embed.FS

// ClientConfig represents the WireGuard client configuration.
type ClientConfig struct{}

// Validate verifies the ClientConfig. (No-op for now.)
func (c *ClientConfig) Validate() error {
	return nil
}

// WriteToFile writes the template to a file using the ClientConfig structure.
func (c *ClientConfig) WriteToFile(name string) error {
	text, err := fs.ReadFile("client.toml.tmpl")
	if err != nil {
		return err
	}

	return utils.ExecTemplateToFile(string(text), c, name)
}

// WriteBuiltToFile writes the built template to a file using the ClientConfig structure.
func (c *ClientConfig) WriteBuiltToFile(name string) error {
	text, err := fs.ReadFile("client_built.conf.tmpl")
	if err != nil {
		return err
	}

	return utils.ExecTemplateToFile(string(text), c, name)
}

// ServerConfig represents the WireGuard server configuration.
type ServerConfig struct {
	IPv4CIDR     string `mapstructure:"ipv4_cidr"`
	IPv6CIDR     string `mapstructure:"ipv6_cidr"`
	Interface    string `mapstructure:"interface"`
	ListenPort   uint16 `mapstructure:"listen_port"`
	OutInterface string `mapstructure:"out_interface"`
	PrivateKey   string `mapstructure:"private_key"`
}

// Address returns the combined IPv4 and IPv6 CIDRs, separated by a comma if both are present.
func (c *ServerConfig) Address() string {
	var addrs []string
	if c.IPv4CIDR != "" {
		addrs = append(addrs, c.IPv4CIDR)
	}
	if c.IPv6CIDR != "" {
		addrs = append(addrs, c.IPv6CIDR)
	}

	return strings.Join(addrs, ", ")
}

// Validate checks that the ServerConfig fields have valid values.
func (c *ServerConfig) Validate() error {
	if c.IPv4CIDR == "" && c.IPv6CIDR == "" {
		return errors.New("either ipv4_cidr or ipv6_cidr is required")
	}
	if c.IPv4CIDR != "" {
		cidr, err := types.NewCIDR(c.IPv4CIDR)
		if err != nil {
			return fmt.Errorf("invalid ipv4_cidr: %w", err)
		}
		if cidr.Len() > 256 {
			return errors.New("ipv4_cidr is too large")
		}
	}
	if c.IPv6CIDR != "" {
		cidr, err := types.NewCIDR(c.IPv6CIDR)
		if err != nil {
			return fmt.Errorf("invalid ipv6_cidr: %w", err)
		}
		if cidr.Len() > 256 {
			return errors.New("ipv6_cidr is too large")
		}
	}
	if c.Interface == "" {
		return errors.New("interface cannot be empty")
	}
	if c.ListenPort == 0 {
		return errors.New("listen_port cannot be zero")
	}
	if c.OutInterface == "" {
		return errors.New("out_interface cannot be empty")
	}
	if c.PrivateKey == "" {
		return errors.New("private_key cannot be empty")
	}
	if _, err := NewKeyFromString(c.PrivateKey); err != nil {
		return fmt.Errorf("invalid private_key: %w", err)
	}

	return nil
}

// WriteToFile writes the template to a file using the ServerConfig structure.
func (c *ServerConfig) WriteToFile(name string) error {
	text, err := fs.ReadFile("server.toml.tmpl")
	if err != nil {
		return err
	}

	return utils.ExecTemplateToFile(string(text), c, name)
}

// WriteBuiltToFile writes the built template to a file using the ServerConfig structure.
func (c *ServerConfig) WriteBuiltToFile(name string) error {
	text, err := fs.ReadFile("server_built.conf.tmpl")
	if err != nil {
		return err
	}

	return utils.ExecTemplateToFile(string(text), c, name)
}

func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		IPv4CIDR:     "10.8.0.1/24",
		IPv6CIDR:     "",
		Interface:    "wg0",
		ListenPort:   51820,
		OutInterface: "eth0",
		PrivateKey:   "",
	}
}
