package wireguard

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"golang.zx2c4.com/wireguard/windows/conf"

	sentinelsdk "github.com/sentinel-official/sentinel-go-sdk/v1/types"
	"github.com/sentinel-official/sentinel-go-sdk/v1/utils"
)

// Ensure Client implements the sentinelsdk.ClientService interface.
var _ sentinelsdk.ClientService = (*Client)(nil)

// Client represents a WireGuard client with associated home directory and name.
type Client struct {
	homeDir string // Home directory for client files.
	name    string // Name of the interface.
}

// NewClient creates a new WireGuard client with the given home directory and name.
func NewClient(homeDir, name string) *Client {
	return &Client{
		homeDir: homeDir,
		name:    name,
	}
}

// configFilePath returns the file path of the client's configuration file.
func (c *Client) configFilePath() string {
	return filepath.Join(c.homeDir, fmt.Sprintf("%s.conf", c.name))
}

// Type returns the service type of the client.
func (c *Client) Type() sentinelsdk.ServiceType {
	return sentinelsdk.ServiceTypeWireGuard
}

// IsUp checks if the WireGuard interface is up.
func (c *Client) IsUp(ctx context.Context) (bool, error) {
	// Retrieves the interface name.
	iface, err := c.interfaceName()
	if err != nil {
		return false, err
	}

	// Executes the 'wg show' command to check the interface status.
	cmd := exec.CommandContext(
		ctx,
		c.execFile("wg"),
		strings.Fields(fmt.Sprintf("show %s", iface))...,
	)
	if err := cmd.Run(); err != nil {
		return false, err
	}

	return true, nil
}

// PreUp writes the configuration to the config file before starting the client process.
func (c *Client) PreUp(v interface{}) error {
	// Checks for valid parameter type.
	cfg, ok := v.(*conf.Config)
	if !ok {
		return fmt.Errorf("invalid parameter type %T", v)
	}

	// Writes configuration to file.
	return os.WriteFile(c.configFilePath(), []byte(cfg.ToWgQuick()), 0600)
}

// PostUp performs operations after the client process is started.
func (c *Client) PostUp() error {
	return nil
}

// PreDown performs operations before the client process is terminated.
func (c *Client) PreDown() error {
	return nil
}

// PostDown performs cleanup operations after the client process is terminated.
func (c *Client) PostDown() error {
	// Removes configuration file.
	if err := utils.RemoveFile(c.configFilePath()); err != nil {
		return err
	}

	return nil
}

// Statistics returns the download and upload statistics for the WireGuard interface.
func (c *Client) Statistics(ctx context.Context) (int64, int64, error) {
	// Retrieves the interface name.
	iface, err := c.interfaceName()
	if err != nil {
		return 0, 0, err
	}

	// Executes the 'wg show [interface] transfer' command to get transfer statistics.
	cmd := exec.CommandContext(
		ctx,
		c.execFile("wg"),
		strings.Fields(fmt.Sprintf("show %s transfer", iface))...,
	)

	output, err := cmd.Output()
	if err != nil {
		return 0, 0, err
	}

	// Parses the output to retrieve download and upload statistics.
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		columns := strings.Split(line, "\t")
		if len(columns) != 3 {
			continue
		}

		download, err := strconv.ParseInt(columns[1], 10, 64)
		if err != nil {
			return 0, 0, err
		}

		upload, err := strconv.ParseInt(columns[2], 10, 64)
		if err != nil {
			return 0, 0, err
		}

		return download, upload, nil
	}

	return 0, 0, nil // Return 0 statistics if no data found.
}
