package v2ray

// ServerMetadata represents metadata for a V2Ray server's inbound connection.
type ServerMetadata struct {
	Tag  *Tag   `json:"tag"`  // Tag uniquely identifies the inbound connection.
	Port string `json:"port"` // Port specifies the listening port for the inbound connection.
}
