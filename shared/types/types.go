package types

import (
	"time"

	rhpv2 "go.sia.tech/core/rhp/v2"
	rhpv3 "go.sia.tech/core/rhp/v3"
	rhpv4 "go.sia.tech/core/rhp/v4"
	"go.sia.tech/coreutils/chain"

	"go.sia.tech/core/types"
)

type Host struct {
	PublicKey      types.PublicKey    `json:"publicKey"`
	V2             bool               `json:"v2"`
	NetAddress     string             `json:"netAddress"`
	V2NetAddresses []chain.NetAddress `json:"v2NetAddresses,omitempty"`

	// Location geoip.Location `json:"location"`

	KnownSince             time.Time `json:"knownSince"`
	LastScan               time.Time `json:"lastScan"`
	NextScan               time.Time `json:"nextScan"`
	LastScanSuccessful     bool      `json:"lastScanSuccessful"`
	LastAnnouncement       time.Time `json:"lastAnnouncement"`
	TotalScans             uint64    `json:"totalScans"`
	SuccessfulInteractions uint64    `json:"successfulInteractions"`
	FailedInteractions     uint64    `json:"failedInteractions"`

	Settings   rhpv2.HostSettings   `json:"settings"`
	PriceTable rhpv3.HostPriceTable `json:"priceTable"`

	RHPV4Settings rhpv4.HostSettings `json:"rhpV4Settings"`
}

// HostScan represents the results of a host scan.
type HostScan struct {
	// PublicKey types.PublicKey `json:"publicKey"`
	// // Location  geoip.Location  `json:"location"`
	// Success   bool      `json:"success"`
	// Timestamp time.Time `json:"timestamp"`
	// NextScan  time.Time `json:"nextScan"`

	// Settings   rhpv2.HostSettings   `json:"settings"`
	// PriceTable rhpv3.HostPriceTable `json:"priceTable"`

	// RHPV4Settings rhpv4.HostSettings `json:"rhpV4Settings"`
	PublicKey          string             `json:"publicKey"`
	V2                 bool               `json:"v2"`
	V2NetAddresses     []chain.NetAddress `json:"v2NetAddresses,omitempty"`
	NetAddress         string             `json:"netAddress"`
	Success            bool               `json:"success"`
	Timestamp          time.Time          `json:"timestamp"`
	NextScan           time.Time          `json:"nextScan"`
	AcceptingContracts bool               `json:"acceptingContracts"`

	Error        string    `json:"error"`
	OnlineSince  time.Time `json:"onlineSince"`
	OfflineSince time.Time `json:"offlineSince"`

	TotalStorage     uint64 `json:"totalStorage"`
	RemainingStorage uint64 `json:"remainingStorage"`

	Settings   rhpv2.HostSettings   `json:"settings"`
	PriceTable rhpv3.HostPriceTable `json:"priceTable"`

	RHPV4Settings rhpv4.HostSettings `json:"rhpV4Settings"`
}
