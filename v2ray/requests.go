package v2ray

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"

	"github.com/v2fly/v2ray-core/v5/common/uuid"
	"google.golang.org/protobuf/types/known/anypb"
)

// AddPeerRequest represents a request to add a peer to the V2Ray server.
type AddPeerRequest struct {
	Protocol Protocol  `json:"protocol"`
	Network  Network   `json:"network"`
	Security Security  `json:"security"`
	UUID     uuid.UUID `json:"uuid"`
}

// Account returns the account information associated with this peer request.
func (r *AddPeerRequest) Account() *anypb.Any {
	tag := r.Tag()
	return tag.Account(r.UUID)
}

// Bytes encodes the AddPeerRequest into a byte slice using gob encoding.
func (r *AddPeerRequest) Bytes() []byte {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(r); err != nil {
		panic(err)
	}

	return buf.Bytes()
}

// Key returns a base64-encoded string of the AddPeerRequest's byte representation.
func (r *AddPeerRequest) Key() string {
	buf := r.Bytes()
	return base64.StdEncoding.EncodeToString(buf)
}

// Tag returns a Tag object based on the protocol, network, and security settings.
func (r *AddPeerRequest) Tag() *Tag {
	return &Tag{
		p: r.Protocol,
		n: r.Network,
		s: r.Security,
	}
}

// Validate checks if the AddPeerRequest contains valid protocol, network, and security settings.
func (r *AddPeerRequest) Validate() error {
	if !r.Protocol.IsValid() {
		return fmt.Errorf("invalid protocol: %s", r.Protocol)
	}
	if !r.Network.IsValid() {
		return fmt.Errorf("invalid network: %s", r.Network)
	}
	if !r.Security.IsValid() {
		return fmt.Errorf("invalid security: %s", r.Security)
	}

	return nil
}

// HasPeerRequest represents a request to check if a peer exists in the V2Ray server.
type HasPeerRequest struct {
	Protocol Protocol  `json:"protocol"`
	Network  Network   `json:"network"`
	Security Security  `json:"security"`
	UUID     uuid.UUID `json:"uuid"`
}

// Bytes encodes the HasPeerRequest into a byte slice using gob encoding.
func (r *HasPeerRequest) Bytes() []byte {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(r); err != nil {
		panic(err)
	}

	return buf.Bytes()
}

// Key returns a base64-encoded string of the HasPeerRequest's byte representation.
func (r *HasPeerRequest) Key() string {
	buf := r.Bytes()
	return base64.StdEncoding.EncodeToString(buf)
}

// Validate checks if the HasPeerRequest contains valid protocol, network, and security settings.
func (r *HasPeerRequest) Validate() error {
	if !r.Protocol.IsValid() {
		return fmt.Errorf("invalid protocol: %s", r.Protocol)
	}
	if !r.Network.IsValid() {
		return fmt.Errorf("invalid network: %s", r.Network)
	}
	if !r.Security.IsValid() {
		return fmt.Errorf("invalid security: %s", r.Security)
	}

	return nil
}

// RemovePeerRequest represents a request to remove a peer from the V2Ray server.
type RemovePeerRequest struct {
	Protocol Protocol  `json:"protocol"`
	Network  Network   `json:"network"`
	Security Security  `json:"security"`
	UUID     uuid.UUID `json:"uuid"`
}

// Bytes encodes the RemovePeerRequest into a byte slice using gob encoding.
func (r *RemovePeerRequest) Bytes() []byte {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(r); err != nil {
		panic(err)
	}

	return buf.Bytes()
}

// Key returns a base64-encoded string of the RemovePeerRequest's byte representation.
func (r *RemovePeerRequest) Key() string {
	buf := r.Bytes()
	return base64.StdEncoding.EncodeToString(buf)
}

// Tag returns a Tag object based on the protocol, network, and security settings.
func (r *RemovePeerRequest) Tag() *Tag {
	return &Tag{
		p: r.Protocol,
		n: r.Network,
		s: r.Security,
	}
}

// Validate checks if the RemovePeerRequest contains valid protocol, network, and security settings.
func (r *RemovePeerRequest) Validate() error {
	if !r.Protocol.IsValid() {
		return fmt.Errorf("invalid protocol: %s", r.Protocol)
	}
	if !r.Network.IsValid() {
		return fmt.Errorf("invalid network: %s", r.Network)
	}
	if !r.Security.IsValid() {
		return fmt.Errorf("invalid security: %s", r.Security)
	}

	return nil
}