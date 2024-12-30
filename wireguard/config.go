package wireguard

import (
	"embed"
	"errors"
	"fmt"
	"net/netip"
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
	text, err := fs.ReadFile("client.conf.tmpl")
	if err != nil {
		return err
	}

	return utils.ExecTemplateToFile(string(text), c, name)
}

// ServerConfig represents the WireGuard server configuration.
type ServerConfig struct {
	IPv4Addr     string `mapstructure:"ipv4_addr"`
	IPv6Addr     string `mapstructure:"ipv6_addr"`
	Interface    string `mapstructure:"interface"`
	ListenPort   uint16 `mapstructure:"listen_port"`
	OutInterface string `mapstructure:"out_interface"`
	PrivateKey   string `mapstructure:"private_key"`
}

// Address returns the combined IPv4 and IPv6 Addrs, separated by a comma if both are present.
func (c *ServerConfig) Address() string {
	var addrs []string
	if c.IPv4Addr != "" {
		addrs = append(addrs, c.IPv4Addr)
	}
	if c.IPv6Addr != "" {
		addrs = append(addrs, c.IPv6Addr)
	}

	return strings.Join(addrs, ", ")
}

// Validate checks that the ServerConfig fields have valid values.
func (c *ServerConfig) Validate() error {
	if c.IPv4Addr == "" && c.IPv6Addr == "" {
		return errors.New("either ipv4_addr or ipv6_addr is required")
	}
	if c.IPv4Addr != "" {
		prefix, err := types.NewNetPrefix(c.IPv4Addr)
		if err != nil {
			return fmt.Errorf("invalid ipv4_addr: %w", err)
		}
		if prefix.Len() > 256 {
			return errors.New("ipv4_addr prefix block is too large")
		}
	}
	if c.IPv6Addr != "" {
		prefix, err := types.NewNetPrefix(c.IPv6Addr)
		if err != nil {
			return fmt.Errorf("invalid ipv6_addr: %w", err)
		}
		if prefix.Len() > 256 {
			return errors.New("ipv6_addr prefix block is too large")
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
	text, err := fs.ReadFile("server.conf.tmpl")
	if err != nil {
		return err
	}

	return utils.ExecTemplateToFile(string(text), c, name)
}

func (c *ServerConfig) IPv4Addrs() ([]netip.Addr, error) {
	prefix, err := types.NewNetPrefix(c.IPv4Addr)
	if err != nil {
		return nil, err
	}

	return prefix.Addrs()
}

func (c *ServerConfig) IPv6Addrs() ([]netip.Addr, error) {
	prefix, err := types.NewNetPrefix(c.IPv6Addr)
	if err != nil {
		return nil, err
	}

	return prefix.Addrs()
}

func DefaultServerConfig() ServerConfig {
	pk, err := NewPrivateKey()
	if err != nil {
		panic(err)
	}

	return ServerConfig{
		IPv4Addr:     "10.8.0.1/24",
		IPv6Addr:     "",
		Interface:    "wg0",
		ListenPort:   utils.RandomPort(),
		OutInterface: "eth0",
		PrivateKey:   pk.String(),
	}
}
