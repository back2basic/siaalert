package scan

import (
	"go.sia.tech/core/consensus"
	rhpv2 "go.sia.tech/core/rhp/v2"
	rhpv3 "go.sia.tech/core/rhp/v3"
	"go.sia.tech/core/types"
)

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
