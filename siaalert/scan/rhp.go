package scan

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"time"

	rhpv2 "go.sia.tech/core/rhp/v2"
	rhpv3 "go.sia.tech/core/rhp/v3"
	rhpv4 "go.sia.tech/core/rhp/v4"
	"go.sia.tech/core/types"
	"go.sia.tech/coreutils/chain"

	crhpv4 "go.sia.tech/coreutils/rhp/v4"
)

// HostScan represents the results of a host scan.
type HostScan struct {
	PublicKey types.PublicKey `json:"publicKey"`
	// Location  geoip.Location  `json:"location"`
	Success   bool      `json:"success"`
	Timestamp time.Time `json:"timestamp"`
	NextScan  time.Time `json:"nextScan"`

	Settings   rhpv2.HostSettings   `json:"settings"`
	PriceTable rhpv3.HostPriceTable `json:"priceTable"`

	RHPV4Settings rhpv4.HostSettings `json:"rhpV4Settings"`
}

// UnscannedHost represents the metadata needed to scan a host.
type UnscannedHost struct {
	PublicKey                types.PublicKey    `json:"publicKey"`
	V2                       bool               `json:"v2"`
	NetAddress               string             `json:"netAddress"`
	V2NetAddresses           []chain.NetAddress `json:"v2NetAddresses,omitempty"`
	FailedInteractionsStreak uint64             `json:"failedInteractionsStreak"`
}

// V2SiamuxAddr returns the `Address` of the first TCP siamux `NetAddress` it
// finds in the host's list of net addresses.  The protocol for this address is
// ProtocolTCPSiaMux.
func (h UnscannedHost) V2SiamuxAddr() (string, bool) {
	for _, netAddr := range h.V2NetAddresses {
		if netAddr.Protocol == crhpv4.ProtocolTCPSiaMux {
			return netAddr.Address, true
		}
	}
	return "", false
}

// IsV2 returns whether a host supports V2 or not.
func (h UnscannedHost) IsV2() bool {
	return len(h.V2NetAddresses) > 0
}

type (
	// A Session is an RHP3 session with the host
	Session struct {
		hostKey types.PublicKey
		cm      ChainManager
		w       Wallet
		t       *rhpv3.Transport

		pt rhpv3.HostPriceTable
	}
)

func RPCSettings(ctx context.Context, t *rhpv2.Transport) (settings rhpv2.HostSettings, err error) {
	var resp rhpv2.RPCSettingsResponse
	if err := t.Call(rhpv2.RPCSettingsID, nil, &resp); err != nil {
		return rhpv2.HostSettings{}, err
	} else if err := json.Unmarshal(resp.Settings, &settings); err != nil {
		return rhpv2.HostSettings{}, fmt.Errorf("couldn't unmarshal json: %w", err)
	}
	return settings, nil
}

// NewSession creates a new session with a host
func NewSession(ctx context.Context, hostKey types.PublicKey, hostAddr string, cm ChainManager, w Wallet) (*Session, error) {
	conn, err := (&net.Dialer{}).DialContext(ctx, "tcp", hostAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to dial host: %w", err)
	}
	t, err := rhpv3.NewRenterTransport(conn, hostKey)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to create transport: %w", err)
	}

	return &Session{
		hostKey: hostKey,
		t:       t,
		w:       w,
		cm:      cm,
	}, nil
}

func ScanPriceTable(v3Session *Session) (rhpv3.HostPriceTable, error) {
	stream := v3Session.t.DialStream()
	defer stream.Close()

	if err := stream.WriteRequest(rhpv3.RPCUpdatePriceTableID, nil); err != nil {
		return rhpv3.HostPriceTable{}, fmt.Errorf("failed to write request: %w", err)
	}
	var resp rhpv3.RPCUpdatePriceTableResponse
	if err := stream.ReadResponse(&resp, 4096); err != nil {
		return rhpv3.HostPriceTable{}, fmt.Errorf("failed to read response: %w", err)
	}

	var pt rhpv3.HostPriceTable
	if err := json.Unmarshal(resp.PriceTableJSON, &pt); err != nil {
		return rhpv3.HostPriceTable{}, fmt.Errorf("failed to unmarshal price table: %w", err)
	}
	return pt, nil
}

func (nc *Checker) ScanV1Host(host UnscannedHost) (HostScan, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*5)
	defer cancel()

	dialer := (&net.Dialer{})

	conn, err := dialer.DialContext(ctx, "tcp", host.NetAddress)
	if err != nil {
		return HostScan{}, fmt.Errorf("scanHost: failed to connect to host: %w", err)
		// return HostScan{}, fmt.Errorf("failed to connect to host")
	}
	defer conn.Close()

	transport, err := rhpv2.NewRenterTransport(conn, host.PublicKey)
	if err != nil {
		return HostScan{}, fmt.Errorf("scanHost: failed to establish v2 transport: %w", err)
		// return HostScan{}, fmt.Errorf("failed to establish v2 transport")
	}
	defer transport.Close()

	settings, err := RPCSettings(ctx, transport)
	if err != nil {
		return HostScan{}, fmt.Errorf("scanHost: failed to get host settings: %w", err)
		// return HostScan{}, fmt.Errorf("failed to get host settings")
	}

	hostIP, _, err := net.SplitHostPort(settings.NetAddress)
	if err != nil {
		return HostScan{}, fmt.Errorf("scanHost: failed to parse net address: %w", err)
		// return HostScan{}, fmt.Errorf("failed to parse net address")
	}

	// resolved, err := net.ResolveIPAddr("ip", hostIP)
	// if err != nil {
	// 	return HostScan{}, fmt.Errorf("scanHost: failed to resolve host address: %w", err)
	// }

	// location, err := locator.Locate(resolved)
	// if err != nil {
	// 	e.log.Debug("Failed to resolve IP geolocation, not setting country code", zap.String("addr", host.NetAddress))
	// }

	v3Addr := net.JoinHostPort(hostIP, settings.SiaMuxPort)
	v3Session, err := NewSession(ctx, host.PublicKey, v3Addr, nc.cm, nil)
	if err != nil {
		return HostScan{}, fmt.Errorf("scanHost: failed to establish v3 transport: %w", err)
		// return HostScan{}, fmt.Errorf("failed to establish v3 transport")
	}

	table, err := ScanPriceTable(v3Session)
	if err != nil {
		return HostScan{}, fmt.Errorf("scanHost: failed to scan price table: %w", err)
		// return HostScan{}, fmt.Errorf("failed to scan price table")
	}

	return HostScan{
		PublicKey: host.PublicKey,
		// Location:  location,
		Success:   true,
		Timestamp: types.CurrentTimestamp(),

		Settings:   settings,
		PriceTable: table,
	}, nil
}

func (nc *Checker) ScanV2Host(host UnscannedHost) (HostScan, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*5)
	defer cancel()

	addr, ok := host.V2SiamuxAddr()
	if !ok {
		return HostScan{}, fmt.Errorf("host has no v2 siamux address")
	}

	transport, err := crhpv4.DialSiaMux(ctx, addr, host.PublicKey)
	if err != nil {
		return HostScan{}, fmt.Errorf("failed to dial host: %w", err)
		// return HostScan{}, fmt.Errorf("failed to dial host")
	}
	defer transport.Close()

	settings, err := crhpv4.RPCSettings(ctx, transport)
	if err != nil {
		return HostScan{}, fmt.Errorf("failed to get host settings: %w", err)
		// return HostScan{}, fmt.Errorf("failed to get host settings")
	}

	// hostIP, _, err := net.SplitHostPort(addr)
	// if err != nil {
	// 	return HostScan{}, fmt.Errorf("scanHost: failed to parse net address: %w", err)
	// }

	// resolved, err := net.ResolveIPAddr("ip", hostIP)
	// if err != nil {
	// 	return HostScan{}, fmt.Errorf("scanHost: failed to resolve host address: %w", err)
	// }

	// location, err := locator.Locate(resolved)
	// if err != nil {
	// 	e.log.Debug("Failed to resolve IP geolocation, not setting country code", zap.String("addr", host.NetAddress))
	// }

	return HostScan{
		PublicKey: host.PublicKey,
		// Location:  location,
		Success:   true,
		Timestamp: types.CurrentTimestamp(),

		RHPV4Settings: settings,
	}, nil
}
