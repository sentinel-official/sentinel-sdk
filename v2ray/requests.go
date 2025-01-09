package v2ray

import (
	"encoding/base64"

	"github.com/v2fly/v2ray-core/v5/common/uuid"
)

// AddPeerRequest represents a request to add a peer.
type AddPeerRequest struct {
	UUID uuid.UUID `json:"uuid"`
}

// Bytes returns the byte representation of the UUID.
func (r *AddPeerRequest) Bytes() []byte {
	return r.UUID.Bytes()
}

// Key returns the base64-encoded byte representation of the UUID.
func (r *AddPeerRequest) Key() string {
	buf := r.Bytes()
	return base64.StdEncoding.EncodeToString(buf)
}

// Validate ensures the request is valid.
func (r *AddPeerRequest) Validate() error {
	return nil
}

// NewAddPeerRequestFromBytes creates an AddPeerRequest from bytes.
func NewAddPeerRequestFromBytes(data []byte) (*AddPeerRequest, error) {
	buf, err := uuid.ParseBytes(data)
	if err != nil {
		return nil, err
	}

	return &AddPeerRequest{
		UUID: buf,
	}, nil
}

// NewAddPeerRequestFromKey creates an AddPeerRequest from a base64-encoded key.
func NewAddPeerRequestFromKey(s string) (*AddPeerRequest, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}

	return NewAddPeerRequestFromBytes(data)
}

// HasPeerRequest represents a request to check if a peer exists.
type HasPeerRequest struct {
	UUID uuid.UUID `json:"uuid"`
}

// Bytes returns the byte representation of the UUID.
func (r *HasPeerRequest) Bytes() []byte {
	return r.UUID.Bytes()
}

// Key returns the base64-encoded byte representation of the UUID.
func (r *HasPeerRequest) Key() string {
	buf := r.Bytes()
	return base64.StdEncoding.EncodeToString(buf)
}

// Validate ensures the request is valid.
func (r *HasPeerRequest) Validate() error {
	return nil
}

// NewHasPeerRequestFromBytes creates a HasPeerRequest from bytes.
func NewHasPeerRequestFromBytes(data []byte) (*HasPeerRequest, error) {
	buf, err := uuid.ParseBytes(data)
	if err != nil {
		return nil, err
	}

	return &HasPeerRequest{
		UUID: buf,
	}, nil
}

// NewHasPeerRequestFromKey creates a HasPeerRequest from a base64-encoded key.
func NewHasPeerRequestFromKey(s string) (*HasPeerRequest, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}

	return NewHasPeerRequestFromBytes(data)
}

// RemovePeerRequest represents a request to remove a peer.
type RemovePeerRequest struct {
	UUID uuid.UUID `json:"uuid"`
}

// Bytes returns the byte representation of the UUID.
func (r *RemovePeerRequest) Bytes() []byte {
	return r.UUID.Bytes()
}

// Key returns the base64-encoded byte representation of the UUID.
func (r *RemovePeerRequest) Key() string {
	buf := r.Bytes()
	return base64.StdEncoding.EncodeToString(buf)
}

// Validate ensures the request is valid.
func (r *RemovePeerRequest) Validate() error {
	return nil
}

// NewRemovePeerRequestFromBytes creates a RemovePeerRequest from bytes.
func NewRemovePeerRequestFromBytes(data []byte) (*RemovePeerRequest, error) {
	buf, err := uuid.ParseBytes(data)
	if err != nil {
		return nil, err
	}

	return &RemovePeerRequest{
		UUID: buf,
	}, nil
}

// NewRemovePeerRequestFromKey creates a RemovePeerRequest from a base64-encoded key.
func NewRemovePeerRequestFromKey(s string) (*RemovePeerRequest, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}

	return NewRemovePeerRequestFromBytes(data)
}
