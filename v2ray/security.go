package v2ray

// TransportSecurity is a custom type used to represent different transport security settings.
type TransportSecurity byte

// Constants for TransportSecurity type with automatic incrementation for each security setting.
const (
	TransportSecurityUnspecified TransportSecurity = iota // Default value for unspecified transport security
	TransportSecurityNone                                 // TransportSecurityNone represents no security
	TransportSecurityTLS                                  // TransportSecurityTLS represents TLS security
)

// String returns a string representation of the TransportSecurity type.
func (t TransportSecurity) String() string {
	switch t {
	case TransportSecurityNone:
		return "none"
	case TransportSecurityTLS:
		return "tls"
	default:
		return "" // Return empty string for unspecified or unknown security settings
	}
}

// IsValid checks if the TransportSecurity value is valid.
func (t TransportSecurity) IsValid() bool {
	return t.String() != ""
}

// NewTransportSecurityFromString converts a string to a TransportSecurity type.
func NewTransportSecurityFromString(v string) TransportSecurity {
	switch v {
	case "none":
		return TransportSecurityNone
	case "tls":
		return TransportSecurityTLS
	default:
		return TransportSecurityUnspecified // Returns the default security if no match is found
	}
}
