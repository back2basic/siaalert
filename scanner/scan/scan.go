package scan

import (
	"context"
	"net"
	"regexp"
	"strconv"
	"time"

	"github.com/back2basic/siaalert/scanner/db"
	"github.com/back2basic/siaalert/scanner/explored"
	"github.com/back2basic/siaalert/scanner/logger"
	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson"
	"go.sia.tech/core/types"
)

type Checker struct {
	cm ChainManager
}

type NetworkChecker interface {
	CheckPortOpen(address string, port string) (bool, time.Duration)
	CheckDNSRecords(hostname string) (hasA, hasAAAA bool, v4, v6 []net.IP)
	ClassifyNetAddress(address string) string
	SplitAddressPort(address string) (host, port string, err error)
	PortScan(hostId string, scanned Scan)
	ScanV1Host(host UnscannedHost) (HostScan, error)
	ScanV2Host(host UnscannedHost) (HostScan, error)
}

func (c *Check) ToBSON() bson.M {
	return bson.M{
		"createdAt":   time.Now(),
		"publicKey":   c.PublicKey,
		"v4addr":      c.V4Addr,
		"v6addr":      c.V6Addr,
		"rhp2port":    c.Rhp2Port,
		"rhp2v4delay": c.Rhp2V4Delay,
		"rhp2v6delay": c.Rhp2V6Delay,
		"rhp2v4":      c.Rhp2V4,
		"rhp2v6":      c.Rhp2V6,
		"rhp3port":    c.Rhp3Port,
		"rhp3v4":      c.Rhp3V4,
		"rhp3v6":      c.Rhp3V6,
		"rhp3v4delay": c.Rhp3V4Delay,
		"rhp3v6delay": c.Rhp3V6Delay,
		"rhp4port":    c.Rhp4Port,
		"rhp4v4":      c.Rhp4V4,
		"rhp4v6":      c.Rhp4V6,
		"rhp4v4delay": c.Rhp4V4Delay,
		"rhp4v6delay": c.Rhp4V6Delay,
	}
}

func (h *HostScan) ToBSON() bson.M {
	return bson.M{
		"publicKey":          h.PublicKey,
		"netAddress":         h.NetAddress,
		"v2":                 h.V2,
		"v2NetAddresses":     h.V2NetAddresses,
		"success":            h.Success,
		"timestamp":          h.Timestamp,
		"nextScan":           h.NextScan,
		"totalStorage":       h.TotalStorage,
		"remainingStorage":   h.RemainingStorage,
		"acceptingContracts": h.AcceptingContracts,
		"error":              h.Error,
		"onlineSince":        h.OnlineSince,
		"offlineSince":       h.OfflineSince,

		// "settings":      h.Settings,
		// "priceTable":    h.PriceTable,

		// "rhpV4Settings": h.RHPV4Settings,
	}
}

func (nc *Checker) CheckPortOpen(address, port string) (bool, float64) {
	timeout := 3 * time.Second
	// measure delay
	start := time.Now()
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(address, port), timeout)
	if err != nil {
		return false, 0
	}
	conn.Close()
	return true, float64(time.Since(start).Milliseconds())
}

func (nc *Checker) CheckDNSRecords(hostname string) (bool, bool, []net.IP, []net.IP) {
	resolver := net.Resolver{
		PreferGo: true,
		// Timeout:  5 * time.Second,
	}

	// Check for A record
	v4, err := resolver.LookupIP(context.Background(), "ip4", hostname)
	hasARecord := err == nil

	// Check for AAAA record
	v6, err := resolver.LookupIP(context.Background(), "ip6", hostname)
	hasAAAARecord := err == nil

	return hasARecord, hasAAAARecord, v4, v6
}

func (nc *Checker) ClassifyNetAddress(address string) string {
	ipv4Pattern := regexp.MustCompile(`^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$`)
	ipv6Pattern := regexp.MustCompile(`^([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}$`)

	if ipv4Pattern.MatchString(address) {
		return "IPv4"
	} else if ipv6Pattern.MatchString(address) {
		return "IPv6"
	} else {
		return "Hostname"
	}
}

func (nc *Checker) SplitAddressPort(address string) (string, string, error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return "", "", err
	}
	return host, port, nil
}

func (nc *Checker) PortScan(hostId types.PublicKey, scanned HostScan, mongoDB *db.MongoDB) {
	// fmt.Println("PortScan", scanned.Settings.NetAddress)
	log := logger.GetLogger()
	defer logger.Sync()

	var netAddress string
	var rhp2, rhp3, rhp4 int

	if scanned.V2 {
		address, port, err := nc.SplitAddressPort(scanned.V2NetAddresses[0].Address)
		if err != nil {
			log.Error("PortScan", zap.Error(err))
		}
		netAddress = address
		rhp4, _ = strconv.Atoi(port)
		rhp3 = rhp4 - 1
		rhp2 = rhp3 - 1
	} else {
		address, port, err := nc.SplitAddressPort(scanned.NetAddress)
		if err != nil {
			log.Error("PortScan", zap.Error(err))
		}
		netAddress = address
		rhp2, _ = strconv.Atoi(port)
		rhp3 = rhp2 + 1
		rhp4 = rhp3 + 1
	}
	strRhp2 := strconv.Itoa(rhp2)
	strRhp3 := strconv.Itoa(rhp3)
	strRhp4 := strconv.Itoa(rhp4)

	// clasify netaddress
	var v4, v6 []net.IP = nil, nil
	var hasARecord, hasAAAARecord bool

	check := Check{}
	// check.HostId = hostId
	check.PublicKey = scanned.PublicKey
	check.Rhp2Port = strRhp2
	check.Rhp3Port = strRhp3
	check.Rhp4Port = strRhp4

	classify := nc.ClassifyNetAddress(netAddress)
	switch classify {
	case "Hostname":
		hasARecord, hasAAAARecord, v4, v6 = nc.CheckDNSRecords(netAddress)
		if hasARecord {
			check.Rhp2V4, check.Rhp2V4Delay = nc.CheckPortOpen(v4[0].String(), strRhp2)
			check.Rhp3V4, check.Rhp3V4Delay = nc.CheckPortOpen(v4[0].String(), strRhp3)
			check.Rhp4V4, check.Rhp4V4Delay = nc.CheckPortOpen(v4[0].String(), strRhp4)
		}
		if hasAAAARecord {
			check.Rhp2V6, check.Rhp2V6Delay = nc.CheckPortOpen(v6[0].String(), strRhp2)
			check.Rhp3V6, check.Rhp3V6Delay = nc.CheckPortOpen(v6[0].String(), strRhp3)
			check.Rhp4V6, check.Rhp4V6Delay = nc.CheckPortOpen(v6[0].String(), strRhp4)
		}

		if hasARecord || hasAAAARecord {
			if len(v4) > 0 {
				check.V4Addr = v4[0].String()
			} else {
				check.V4Addr = ""
			}
			if len(v6) > 0 {
				check.V6Addr = v6[0].String()
			} else {
				check.V6Addr = ""
			}

			mongoDB.InsertScan(check.ToBSON())
		}

	case "IPv4":
		check.V4Addr = netAddress
		check.V6Addr = ""
		check.Rhp2V4, check.Rhp2V4Delay = nc.CheckPortOpen(netAddress, strRhp2)
		check.Rhp3V4, check.Rhp3V4Delay = nc.CheckPortOpen(netAddress, strRhp3)
		check.Rhp4V4, check.Rhp4V4Delay = nc.CheckPortOpen(netAddress, strRhp4)

		mongoDB.InsertScan(check.ToBSON())

	case "IPv6":
		check.V4Addr = ""
		check.V6Addr = netAddress
		check.Rhp2V6, check.Rhp2V6Delay = nc.CheckPortOpen(netAddress, strRhp2)
		check.Rhp3V6, check.Rhp3V6Delay = nc.CheckPortOpen(netAddress, strRhp3)
		check.Rhp4V6, check.Rhp4V6Delay = nc.CheckPortOpen(netAddress, strRhp4)

		mongoDB.InsertScan(check.ToBSON())
	}
}

func (nc *Checker) CheckVersion(publicKey string) (string, error) {
	host, err := explored.GetHostByPublicKey(publicKey)
	if err != nil {
		return "", err
	}
	return host.Settings.Version, nil
}
