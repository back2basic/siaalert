package explored

import (
	"time"

	rhpv2 "go.sia.tech/core/rhp/v2"
	rhpv3 "go.sia.tech/core/rhp/v3"
	rhpv4 "go.sia.tech/core/rhp/v4"
)

type Consensus struct {
	Index Index `json:"index"`
}

type Index struct {
	Height uint64 `json:"height"`
	Id     string `json:"id"`
}

type Host struct {
	PublicKey      string     `json:"publicKey"`
	V2             bool       `json:"v2"`
	NetAddress     string     `json:"netAddress"`
	V2NetAddresses NetAddress `json:"v2NetAddresses,omitempty"`

	CountryCode string `json:"countryCode"`

	KnownSince             time.Time `json:"knownSince"`
	LastScan               time.Time `json:"lastScan"`
	LastScanSuccessful     bool      `json:"lastScanSuccessful"`
	LastAnnouncement       time.Time `json:"lastAnnouncement"`
	TotalScans             uint64    `json:"totalScans"`
	SuccessfulInteractions uint64    `json:"successfulInteractions"`
	FailedInteractions     uint64    `json:"failedInteractions"`

	Settings   Settings      `json:"settings"`
	PriceTable PriceTable    `json:"priceTable"`
	RHPV4Settings RHPV4Settings `json:"rhpV4Settings"`
}

type (
	Settings      rhpv2.HostSettings
	PriceTable    rhpv3.HostPriceTable
	RHPV4Settings rhpv4.HostSettings
)

type NetAddress struct {
	Protocol string `json:"protocol"`
	Address  string `json:"address"`
}
