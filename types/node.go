package types

import (
	"cosmossdk.io/math"
	sentinelhub "github.com/sentinel-official/hub/v12/types"
	"github.com/sentinel-official/hub/v12/types/v1"

	"github.com/sentinel-official/sentinel-go-sdk/libs/geoip"
)

type (
	BandwidthInfo struct {
		Down math.Int `json:"down"`
		Up   math.Int `json:"up"`
	}
	HandshakeInfo struct {
		Enable bool  `json:"enable"`
		Peers  int64 `json:"peers"`
	}
	QOSInfo struct {
		MaxPeers int `json:"max_peers"`
	}
)

type NodeInfo struct {
	Addr           sentinelhub.NodeAddress `json:"addr"`
	Bandwidth      *BandwidthInfo          `json:"bandwidth"`
	GigabytePrices v1.Prices               `json:"gigabyte_prices"`
	Handshake      *HandshakeInfo          `json:"handshake"`
	HourlyPrices   v1.Prices               `json:"hourly_prices"`
	Location       *geoip.Location         `json:"location"`
	Moniker        string                  `json:"moniker"`
	Peers          int                     `json:"peers"`
	QOS            *QOSInfo                `json:"qos"`
	Type           ServiceType             `json:"type"`
	Version        string                  `json:"version"`
}
