package geoip

import (
	"time"

	"github.com/sentinel-official/sentinel-go-sdk/utils"
)

// Location represents geographical location information associated with an IP address.
type Location struct {
	City      string  `json:"city,omitempty"`      // City where the IP address is located.
	Country   string  `json:"country,omitempty"`   // Country where the IP address is located.
	IP        string  `json:"ip,omitempty"`        // IP address that was resolved.
	Latitude  float64 `json:"latitude,omitempty"`  // Latitude of the location.
	Longitude float64 `json:"longitude,omitempty"` // Longitude of the location.
}

func (l *Location) String() string {
	return string(utils.MustMarshalJSON(l))
}

// Client is an interface for resolving IP addresses into location data.
type Client interface {
	Get(ip string) (*Location, error)
}

// NewDefaultClient creates a new default Client instance using the default IPAPIClient.
func NewDefaultClient() Client {
	return NewIPAPIClient(15 * time.Second)
}
