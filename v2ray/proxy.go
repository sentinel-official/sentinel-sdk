package v2ray

// ProxyProtocol is a custom type used to represent different proxy protocols.
type ProxyProtocol byte

// Constants for ProxyProtocol type with automatic incrementation for each protocol method.
const (
	ProxyProtocolUnspecified ProxyProtocol = iota // Default value for unspecified protocol
	ProxyProtocolVLess                            // ProxyProtocolVLess represents the VLess protocol
	ProxyProtocolVMess                            // ProxyProtocolVMess represents the VMess protocol
)

// String returns a string representation of the ProxyProtocol type.
func (p ProxyProtocol) String() string {
	switch p {
	case ProxyProtocolVLess:
		return "vless"
	case ProxyProtocolVMess:
		return "vmess"
	default:
		return "" // Return empty string for unspecified or unknown protocols
	}
}

// IsValid checks if the ProxyProtocol value is valid.
func (p ProxyProtocol) IsValid() bool {
	return p.String() != ""
}

// NewProxyProtocolFromString converts a string to a ProxyProtocol type.
func NewProxyProtocolFromString(v string) ProxyProtocol {
	switch v {
	case "vless":
		return ProxyProtocolVLess
	case "vmess":
		return ProxyProtocolVMess
	default:
		return ProxyProtocolUnspecified // Returns the default protocol if no match is found
	}
}
