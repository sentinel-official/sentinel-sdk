package v2ray

// TransportProtocol is a custom type used to represent different transport protocols.
type TransportProtocol byte

// Constants for TransportProtocol type with automatic incrementation for each transport method.
const (
	TransportProtocolUnspecified  TransportProtocol = iota // Default value for unspecified transport protocol
	TransportProtocolDomainSocket                          // TransportProtocolDomainSocket represents a UNIX domain socket
	TransportProtocolGUN                                   // TransportProtocolGUN represents the GUN protocol
	TransportProtocolGRPC                                  // TransportProtocolGRPC represents gRPC, a high-performance RPC framework
	TransportProtocolHTTP                                  // TransportProtocolHTTP represents the HTTP protocol
	TransportProtocolMKCP                                  // TransportProtocolMKCP represents the MKCP (modified KCP) protocol
	TransportProtocolQUIC                                  // TransportProtocolQUIC represents the QUIC protocol
	TransportProtocolTCP                                   // TransportProtocolTCP represents the TCP transport protocol
	TransportProtocolWebSocket                             // TransportProtocolWebSocket represents the WebSocket protocol
)

// String returns a string representation of the TransportProtocol type.
func (t TransportProtocol) String() string {
	switch t {
	case TransportProtocolDomainSocket:
		return "domainsocket"
	case TransportProtocolGUN:
		return "gun"
	case TransportProtocolGRPC:
		return "grpc"
	case TransportProtocolHTTP:
		return "http"
	case TransportProtocolMKCP:
		return "mkcp"
	case TransportProtocolQUIC:
		return "quic"
	case TransportProtocolTCP:
		return "tcp"
	case TransportProtocolWebSocket:
		return "websocket"
	default:
		return "" // Return empty string for unspecified or unknown transport protocol types
	}
}

// IsValid checks if the TransportProtocol value is valid.
func (t TransportProtocol) IsValid() bool {
	return t.String() != ""
}

// NewTransportProtocolFromString converts a string to a TransportProtocol type.
func NewTransportProtocolFromString(v string) TransportProtocol {
	switch v {
	case "domainsocket":
		return TransportProtocolDomainSocket
	case "gun":
		return TransportProtocolGUN
	case "grpc":
		return TransportProtocolGRPC
	case "http":
		return TransportProtocolHTTP
	case "mkcp":
		return TransportProtocolMKCP
	case "quic":
		return TransportProtocolQUIC
	case "tcp":
		return TransportProtocolTCP
	case "websocket", "ws":
		return TransportProtocolWebSocket
	default:
		return TransportProtocolUnspecified // Returns the default transport protocol type if no match is found
	}
}
