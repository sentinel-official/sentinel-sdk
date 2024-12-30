package v2ray

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/v4/process"

	"github.com/sentinel-official/sentinel-go-sdk/types"
	"github.com/sentinel-official/sentinel-go-sdk/utils"
)

// Ensure Client implements the types.ClientService interface.
var _ types.ClientService = (*Client)(nil)

// Client represents a V2Ray client with associated command, home directory, and name.
type Client struct {
	cmd     *exec.Cmd // Command for running the V2Ray client.
	homeDir string    // Home directory for client files.
	name    string    // Name of the interface.
}

// NewClient creates a new Client instance.
func NewClient() *Client {
	return &Client{}
}

// configFilePath returns the file path of the client's configuration file.
func (c *Client) configFilePath() string {
	return filepath.Join(c.homeDir, fmt.Sprintf("%s.json", c.name))
}

// pidFilePath returns the file path of the client's PID file.
func (c *Client) pidFilePath() string {
	return filepath.Join(c.homeDir, fmt.Sprintf("%s.pid", c.name))
}

// readPIDFromFile reads the PID from the client's PID file.
func (c *Client) readPIDFromFile() (int32, error) {
	// Reads PID from the PID file.
	data, err := os.ReadFile(c.pidFilePath())
	if err != nil {
		return 0, err
	}

	// Converts PID data to integer.
	pid, err := strconv.ParseInt(string(data), 10, 32)
	if err != nil {
		return 0, err
	}

	return int32(pid), nil
}

// writePIDToFile writes the given PID to the client's PID file.
func (c *Client) writePIDToFile(pid int) error {
	// Converts PID to byte slice.
	data := []byte(strconv.Itoa(pid))

	// Writes PID to file with appropriate permissions.
	if err := os.WriteFile(c.pidFilePath(), data, 0644); err != nil {
		return err
	}

	return nil
}

// Type returns the service type of the client.
func (c *Client) Type() types.ServiceType {
	return types.ServiceTypeV2Ray
}

// IsUp checks if the V2Ray client process is running.
func (c *Client) IsUp(ctx context.Context) (bool, error) {
	// Reads PID from file.
	pid, err := c.readPIDFromFile()
	if err != nil {
		return false, err
	}

	// Retrieves process with the given PID.
	proc, err := process.NewProcessWithContext(ctx, pid)
	if err != nil {
		return false, err
	}

	// Checks if the process is running.
	ok, err := proc.IsRunningWithContext(ctx)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}

	// Retrieves the name of the process.
	name, err := proc.NameWithContext(ctx)
	if err != nil {
		return false, err
	}

	// Checks if the process name matches constant v2ray.
	if name != v2ray {
		return false, nil
	}

	return true, nil
}

// PreUp writes the configuration to the config file before starting the client process.
func (c *Client) PreUp(v interface{}) error {
	// Checks for valid parameter type.
	cfg, ok := v.(*ClientConfig)
	if !ok {
		return fmt.Errorf("invalid parameter type %T", v)
	}

	// Writes configuration to file.
	return cfg.WriteToFile(c.configFilePath())
}

// Up starts the V2Ray client process.
func (c *Client) Up(ctx context.Context) error {
	// Constructs the command to start the V2Ray client.
	c.cmd = exec.CommandContext(
		ctx,
		c.execFile(v2ray),
		strings.Fields(fmt.Sprintf("run --config %s", c.configFilePath()))...,
	)
	c.cmd.Stdout = os.Stdout
	c.cmd.Stderr = os.Stderr

	// Starts the V2Ray client process.
	return c.cmd.Start()
}

// PostUp performs operations after the client process is started.
func (c *Client) PostUp() error {
	// Checks if command or process is nil.
	if c.cmd == nil || c.cmd.Process == nil {
		return fmt.Errorf("nil command or process")
	}

	// Writes PID to file.
	if err := c.writePIDToFile(c.cmd.Process.Pid); err != nil {
		return err
	}

	return nil
}

// PreDown performs operations before the client process is terminated.
func (c *Client) PreDown() error {
	return nil
}

// Down terminates the V2Ray client process.
func (c *Client) Down(ctx context.Context) error {
	// Reads PID from file.
	pid, err := c.readPIDFromFile()
	if err != nil {
		return err
	}

	// Retrieves process with the given PID.
	proc, err := process.NewProcessWithContext(ctx, pid)
	if err != nil {
		return err
	}

	// Terminates the process.
	if err := proc.TerminateWithContext(ctx); err != nil {
		return err
	}

	// Resets the command.
	c.cmd = nil
	return nil
}

// PostDown performs cleanup operations after the client process is terminated.
func (c *Client) PostDown() error {
	// Removes configuration file.
	if err := utils.RemoveFile(c.configFilePath()); err != nil {
		return err
	}

	// Removes PID file.
	if err := utils.RemoveFile(c.pidFilePath()); err != nil {
		return err
	}

	return nil
}

// Statistics returns dummy statistics for now (to be implemented).
func (c *Client) Statistics(_ context.Context) (int64, int64, error) {
	return 0, 0, nil
}
