package sdk

import (
	"time"

	"github.com/appwrite/sdk-for-go/models"
)

type Host struct {
	PublicKey              string `json:"publicKey"`
	V2                     bool   `json:"v2,omitempty"`
	NetAddress             string `json:"netAddress,omitempty"`
	V2NetAddresses         string `json:"v2NetAddresses,omitempty"`
	V2NetAddressesProto    string `json:"v2NetAddressesProto,omitempty"`
	CountryCode            string `json:"countryCode,omitempty"`
	KnownSince             string `json:"knownSince,omitempty"`
	LastScan               string `json:"lastScan,omitempty"`
	LastScanSuccessful     bool   `json:"lastScanSuccessful,omitempty"`
	LastAnnouncement       string `json:"lastAnnouncement,omitempty"`
	TotalScans             uint64 `json:"totalScans,omitempty"`
	SuccessfulInteractions uint64 `json:"successfulInteractions,omitempty"`
	FailedInteractions     uint64 `json:"failedInteractions,omitempty"`

	Error        string `json:"error"`
	Online       bool   `json:"online"`
	OnlineSince  string `json:"onlineSince"`
	OfflineSince string `json:"offlineSince"`
}

type HostDocument struct {
	models.Document
	Host
}

type HostList struct {
	Documents []HostDocument `json:"documents"`
	Total     uint64         `json:"total"`
}

type Status struct {
	Height uint64 `json:"height"`
}

type StatusDocument struct {
	models.Document
	Status
}

type StatusList struct {
	Documents []StatusDocument `json:"documents"`
	Total     uint64           `json:"total"`
}

type Alert struct {
	HostId string `json:"hostId"`
	Type   string `json:"type"`
	Sender string `json:"sender"`
}

type AlertDocument struct {
	models.Document
	Alert
}

type AlertList struct {
	Documents []AlertDocument `json:"documents"`
	Total     uint64          `json:"total"`
}

type Check struct {
	HostId             string  `json:"hostId"`
	V4Addr             string  `json:"v4Addr"`
	V6Addr             string  `json:"v6Addr"`
	Rhp2Port           string  `json:"rhp2Port"`
	Rhp2V4Delay        float64 `json:"rhp2V4Delay"`
	Rhp2V6Delay        float64 `json:"rhp2V6Delay"`
	Rhp2V4             bool    `json:"rhp2V4"`
	Rhp2V6             bool    `json:"rhp2V6"`
	Rhp3Port           string  `json:"rhp3Port"`
	Rhp3V4             bool    `json:"rhp3V4"`
	Rhp3V6             bool    `json:"rhp3V6"`
	Rhp3V4Delay        float64 `json:"rhp3V4Delay"`
	Rhp3V6Delay        float64 `json:"rhp3V6Delay"`
	Rhp4Port           string  `json:"rhp4Port"`
	Rhp4V4             bool    `json:"rhp4V4"`
	Rhp4V6             bool    `json:"rhp4V6"`
	Rhp4V4Delay        float64 `json:"rhp4V4Delay"`
	Rhp4V6Delay        float64 `json:"rhp4V6Delay"`
	AcceptingContracts bool    `json:"acceptingContracts"`
	Release            string  `json:"release"`
}

type CheckDocument struct {
	models.Document
	Check
}

type CheckList struct {
	Documents []CheckDocument `json:"documents"`
	Total     uint64          `json:"total"`
}

type CheckParams struct {
	Rhp2v4, Rhp2v6, Rhp3v4, Rhp3v6, Rhp4v4, Rhp4v6        bool
	HasARecord, HasAAAARecord, AcceptingContracts         bool
	Rhp2v4Delay, Rhp2v6Delay, Rhp3v4Delay, Rhp3v6Delay    time.Duration
	Rhp4v4Delay, Rhp4v6Delay                              time.Duration
	V4, V6, HostId, Rhp2Port, Rhp3Port, Rhp4Port, Release string
}

type Rhp2 struct {
	AcceptingContracts   bool   `json:"acceptingcontracts"`
	MaxDownloadBatchSize uint64 `json:"maxdownloadbatchsize"`
	MaxDuration          uint64 `json:"maxduration"`
	MaxReviseBatchSize   uint64 `json:"maxrevisebatchsize"`
	RemainingStorage     uint64 `json:"remainingstorage"`
	TotalStorage         uint64 `json:"totalstorage"`
	RevisionNumber       uint64 `json:"revisionnumber"`
	Version              string `json:"version"`
	Release              string `json:"release"`
	SiaMuxPort           string `json:"siamuxport"`
	HostId               string `json:"hostId"`
}

type Rhp3 struct {
	HostBlockHeight uint64 `json:"hostblockheight"`
	HostId          string `json:"hostId"`
}

type Rhp4 struct {
	HostId string `json:"hostId"`
}
