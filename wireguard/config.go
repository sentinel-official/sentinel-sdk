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
	InInterface  string `mapstructure:"in_interface"`
	IPv4Addr     string `mapstructure:"ipv4_addr"`
	IPv6Addr     string `mapstructure:"ipv6_addr"`
	OutInterface string `mapstructure:"out_interface"`
	Port         string `mapstructure:"port"`
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

func (c *ServerConfig) ListenPort() uint16 {
	v, err := types.NewPortFromString(c.Port)
	if err != nil {
		panic(err)
	}

	return v.InFrom
}

// Validate checks that the ServerConfig fields have valid values.
func (c *ServerConfig) Validate() error {
	if c.InInterface == "" {
		return errors.New("in_interface cannot be empty")
	}
	if c.IPv4Addr == "" && c.IPv6Addr == "" {
		return errors.New("either ipv4_addr or ipv6_addr is required")
	}
	if c.IPv4Addr != "" {
		if _, err := types.NewNetPrefixFromString(c.IPv4Addr); err != nil {
			return fmt.Errorf("invalid ipv4_addr: %w", err)
		}
	}
	if c.IPv6Addr != "" {
		if _, err := types.NewNetPrefixFromString(c.IPv6Addr); err != nil {
			return fmt.Errorf("invalid ipv6_addr: %w", err)
		}
	}
	if c.OutInterface == "" {
		return errors.New("out_interface cannot be empty")
	}
	if c.Port == "" {
		return errors.New("port cannot be empty")
	}
	if _, err := types.NewPortFromString(c.Port); err != nil {
		return fmt.Errorf("invalid port: %w", err)
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

func (c *ServerConfig) IPv4Addrs() []netip.Addr {
	prefix, err := types.NewNetPrefixFromString(c.IPv4Addr)
	if err != nil {
		panic(err)
	}

	addrs, err := prefix.Addrs()
	if err != nil {
		panic(err)
	}

	return addrs
}

func (c *ServerConfig) IPv6Addrs() []netip.Addr {
	prefix, err := types.NewNetPrefixFromString(c.IPv6Addr)
	if err != nil {
		panic(err)
	}

	addrs, err := prefix.Addrs()
	if err != nil {
		panic(err)
	}

	return addrs
}

func DefaultServerConfig() ServerConfig {
	pk, err := NewPrivateKey()
	if err != nil {
		panic(err)
	}

	return ServerConfig{
		InInterface:  "wg0",
		IPv4Addr:     "10.8.0.1/24",
		IPv6Addr:     "",
		OutInterface: "eth0",
		Port:         fmt.Sprintf("%d", utils.RandomPort()),
		PrivateKey:   pk.String(),
	}
}
