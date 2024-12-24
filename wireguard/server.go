package wireguard

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/sentinel-official/sentinel-go-sdk/types"
	"github.com/sentinel-official/sentinel-go-sdk/utils"
)

// Ensure Server implements types.ServerService interface.
var _ types.ServerService = (*Server)(nil)

// Server represents the WireGuard server instance.
type Server struct {
	homeDir string       // Home directory of the WireGuard server.
	info    []byte       // Information about the server instance.
	name    string       // Name of the server instance.
	pm      *PeerManager // Peer manager for handling peer information.
}

// NewServer creates a new Server instance.
func NewServer() *Server {
	return &Server{}
}

// WithHomeDir sets the home directory for the server and returns the updated Server instance.
func (s *Server) WithHomeDir(homeDir string) *Server {
	s.homeDir = homeDir
	return s
}

// WithName sets the name for the server and returns the updated Server instance.
func (s *Server) WithName(name string) *Server {
	s.name = name
	return s
}

// WithPeerManager sets the PeerManager for the server and returns the updated Server instance.
func (s *Server) WithPeerManager(pm *PeerManager) *Server {
	s.pm = pm
	return s
}

// Info returns the server's information.
func (s *Server) Info() []byte {
	return s.info
}

// configFilePath returns the file path of the server's configuration file.
func (s *Server) configFilePath() string {
	return filepath.Join(s.homeDir, fmt.Sprintf("%s.conf", s.name))
}

// Type returns the service type of the server.
func (s *Server) Type() types.ServiceType {
	return types.ServiceTypeWireGuard
}

// IsUp checks if the WireGuard server process is running.
func (s *Server) IsUp(ctx context.Context) (bool, error) {
	// Retrieves the interface name.
	iface, err := s.interfaceName()
	if err != nil {
		return false, err
	}

	// Executes the 'wg show' command to check the interface status.
	cmd := exec.CommandContext(
		ctx,
		s.execFile("wg"),
		strings.Fields(fmt.Sprintf("show %s", iface))...,
	)
	if err := cmd.Run(); err != nil {
		return false, err
	}

	return true, nil
}

// PreUp writes the configuration to the config file before starting the server process.
func (s *Server) PreUp(v interface{}) error {
	// Checks for valid parameter type.
	cfg, ok := v.(*ServerConfig)
	if !ok {
		return fmt.Errorf("invalid parameter type %T", v)
	}

	// Writes configuration to file.
	return cfg.WriteBuiltToFile(s.configFilePath())
}

// PostUp performs operations after the server process is started.
func (s *Server) PostUp() error {
	return nil
}

// PreDown performs operations before the server process is terminated.
func (s *Server) PreDown() error {
	return nil
}

// PostDown performs cleanup operations after the server process is terminated.
func (s *Server) PostDown() error {
	// Removes configuration file.
	if err := utils.RemoveFile(s.configFilePath()); err != nil {
		return err
	}

	return nil
}

// AddPeer adds a new peer to the WireGuard server.
func (s *Server) AddPeer(ctx context.Context, req interface{}) (res []byte, err error) {
	// Cast the request to AddPeerRequest type.
	r, ok := req.(*AddPeerRequest)
	if !ok {
		return nil, fmt.Errorf("invalid request type: %T", req)
	}
	if err := r.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Retrieve the identity from the request.
	identity := r.Key()

	// Add peer to the peer manager and retrieve assigned IP addresses.
	ipv4Addr, ipv6Addr, err := s.pm.Put(identity)
	if err != nil {
		return nil, err
	}

	// Executes the 'wg set' command to add the peer to the WireGuard interface.
	cmd := exec.CommandContext(
		ctx,
		s.execFile("wg"),
		strings.Fields(fmt.Sprintf("set %s peer %s allowed-ips %s/32,%s/128", s.name, identity, ipv4Addr, ipv6Addr))...,
	)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	// Run the command and check for errors.
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	// Append IP addresses to the response.
	res = append(res, ipv4Addr...)
	res = append(res, ipv6Addr...)
	return res, nil
}

// HasPeer checks if a peer exists in the WireGuard server's peer list.
func (s *Server) HasPeer(_ context.Context, req interface{}) (bool, error) {
	// Cast the request to HasPeerRequest type.
	r, ok := req.(*HasPeerRequest)
	if !ok {
		return false, fmt.Errorf("invalid request type: %T", req)
	}
	if err := r.Validate(); err != nil {
		return false, fmt.Errorf("invalid request: %w", err)
	}

	// Retrieve the identity from the request.
	identity := r.Key()
	peer := s.pm.Get(identity)

	// Return true if the peer exists, otherwise false.
	return peer != nil, nil
}

// RemovePeer removes a peer from the WireGuard server.
func (s *Server) RemovePeer(ctx context.Context, req interface{}) error {
	// Cast the request to RemovePeerRequest type.
	r, ok := req.(*RemovePeerRequest)
	if !ok {
		return fmt.Errorf("invalid request type: %T", req)
	}
	if err := r.Validate(); err != nil {
		return fmt.Errorf("invalid request: %w", err)
	}

	// Retrieve the identity from the request.
	identity := r.Key()

	// Executes the 'wg set' command to remove the peer from the WireGuard interface.
	cmd := exec.CommandContext(
		ctx,
		s.execFile("wg"),
		strings.Fields(fmt.Sprintf(`set %s peer %s remove`, s.name, identity))...,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command and check for errors.
	if err := cmd.Run(); err != nil {
		return err
	}

	// Remove the peer information from the local collection.
	s.pm.Delete(identity)
	return nil
}

// PeerCount returns the number of peers connected to the WireGuard server.
func (s *Server) PeerCount() int {
	return s.pm.Len()
}

// PeerStatistics retrieves statistics for each peer connected to the WireGuard server.
func (s *Server) PeerStatistics(ctx context.Context) (items []*types.PeerStatistic, err error) {
	// Retrieves the interface name.
	iface, err := s.interfaceName()
	if err != nil {
		return nil, err
	}

	// Executes the 'wg show' command to get transfer statistics.
	output, err := exec.CommandContext(
		ctx,
		s.execFile("wg"),
		strings.Fields(fmt.Sprintf("show %s transfer", iface))...,
	).Output()
	if err != nil {
		return nil, err
	}

	// Split the command output into lines and process each line.
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		columns := strings.Split(line, "\t")
		if len(columns) != 3 {
			continue
		}

		// Parse upload traffic stats.
		uploadBytes, err := strconv.ParseInt(columns[1], 10, 64)
		if err != nil {
			return nil, err
		}

		// Parse download traffic stats.
		downloadBytes, err := strconv.ParseInt(columns[2], 10, 64)
		if err != nil {
			return nil, err
		}

		// Append peer statistics to the result collection.
		items = append(
			items,
			&types.PeerStatistic{
				Key:           columns[0],
				DownloadBytes: downloadBytes,
				UploadBytes:   uploadBytes,
			},
		)
	}

	// Return the constructed collection of peer statistics.
	return items, nil
}
