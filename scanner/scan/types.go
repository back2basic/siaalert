package scan

import (
	"time"

	"go.sia.tech/core/consensus"
	rhpv2 "go.sia.tech/core/rhp/v2"
	rhpv3 "go.sia.tech/core/rhp/v3"
	rhpv4 "go.sia.tech/core/rhp/v4"
	"go.sia.tech/core/types"
	"go.sia.tech/coreutils/chain"
)

type TaskCheckDoc struct {
	ID      int
	Job     string
	CheckID string
	Check   Check
}

type Scan struct {
	Settings   rhpv2.HostSettings
	PriceTable rhpv3.HostPriceTable
}

type (
	// An accountPayment pays for usage using an ephemeral account
	accountPayment struct {
		Account    rhpv3.Account
		PrivateKey types.PrivateKey
	}

	// A contractPayment pays for usage using a contract
	contractPayment struct {
		Revision      *rhpv2.ContractRevision
		RefundAccount rhpv3.Account
		RenterKey     types.PrivateKey
	}

	// A PaymentMethod facilitates payments to the host using either a contract
	// or an ephemeral account
	PaymentMethod interface {
		// Pay(amount types.Currency, height uint64) (rhp3.PaymentMethod, bool)
	}

	// A Wallet funds and signs transactions
	Wallet interface {
		Address() types.Address
		FundTransaction(txn *types.Transaction, amount types.Currency) ([]types.Hash256, func(), error)
		SignTransaction(cs consensus.State, txn *types.Transaction, toSign []types.Hash256, cf types.CoveredFields) error
	}

	// A ChainManager is used to get the current consensus state
	ChainManager interface {
		TipState() consensus.State
	}
)

// A Check is a port scan of a host
type Check struct {
	CreatedAt   time.Time
	PublicKey   string  `json:"publicKey"`
	V4Addr      string  `json:"v4Addr"`
	V6Addr      string  `json:"v6Addr"`
	Rhp2Port    string  `json:"rhp2Port"`
	Rhp2V4Delay float64 `json:"rhp2V4Delay"`
	Rhp2V6Delay float64 `json:"rhp2V6Delay"`
	Rhp2V4      bool    `json:"rhp2V4"`
	Rhp2V6      bool    `json:"rhp2V6"`
	Rhp3Port    string  `json:"rhp3Port"`
	Rhp3V4      bool    `json:"rhp3V4"`
	Rhp3V6      bool    `json:"rhp3V6"`
	Rhp3V4Delay float64 `json:"rhp3V4Delay"`
	Rhp3V6Delay float64 `json:"rhp3V6Delay"`
	Rhp4Port    string  `json:"rhp4Port"`
	Rhp4V4      bool    `json:"rhp4V4"`
	Rhp4V6      bool    `json:"rhp4V6"`
	Rhp4V4Delay float64 `json:"rhp4V4Delay"`
	Rhp4V6Delay float64 `json:"rhp4V6Delay"`
}

// HostScan represents the results of a host scan.
type HostScan struct {
	// Location  geoip.Location  `json:"location"`
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
