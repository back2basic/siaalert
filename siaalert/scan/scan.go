package scan

import (
	"context"
	"fmt"
	"net"
	"regexp"
	"sync"
	"time"

	"github.com/back2basic/siadata/siaalert/bench"
	"github.com/back2basic/siadata/siaalert/explored"
	"github.com/back2basic/siadata/siaalert/sdk"
)

type Checker struct {
	cm ChainManager
}

type NetworkChecker interface {
	CheckPortOpen(address string, port string) (bool, time.Duration)
	CheckDNSRecords(hostname string) (hasA, hasAAAA bool, v4, v6 []net.IP)
	ClassifyNetAddress(address string) string
	SplitAddressPort(address string) (host, port string, err error)
	PortScan(hostId string, scanned bench.Scan)
	ScanV1Host(host UnscannedHost) (HostScan, error)
	ScanV2Host(host UnscannedHost) (HostScan, error)
}

func (nc *Checker) CheckPortOpen(address, port string) (bool, time.Duration) {
	timeout := 3 * time.Second
	// measure delay
	start := time.Now()
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(address, port), timeout)
	if err != nil {
		return false, time.Since(start)
	}
	conn.Close()
	return true, time.Since(start)
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

func (nc *Checker) PortScan(hostId string, scanned HostScan, wg *sync.WaitGroup, task chan sdk.TaskCheckDoc) {
	// fmt.Println("PortScan", scanned.Settings.NetAddress)
	netAddress, rhp2, err := nc.SplitAddressPort(scanned.Settings.NetAddress)
	if err != nil {
		fmt.Println("PortScan", err)
		return
	}

	rhp3 := scanned.Settings.SiaMuxPort
	// need to be changed to to the v2 address WIP!
	rhp4 := "9984"

	// clasify netaddress
	var v4, v6 []net.IP = nil, nil
	params := sdk.CheckParams{}
	params.HostId = hostId
	params.Rhp2Port = rhp2
	params.Rhp3Port = rhp3
	params.Rhp4Port = rhp4
	params.AcceptingContracts = scanned.Settings.AcceptingContracts
	classify := nc.ClassifyNetAddress(netAddress)
	switch classify {
	case "Hostname":
		params.HasARecord, params.HasAAAARecord, v4, v6 = nc.CheckDNSRecords(netAddress)
		if params.HasARecord {
			params.Rhp2v4, params.Rhp2v4Delay = nc.CheckPortOpen(v4[0].String(), rhp2)
			params.Rhp3v4, params.Rhp3v4Delay = nc.CheckPortOpen(v4[0].String(), rhp3)
			params.Rhp4v4, params.Rhp4v4Delay = nc.CheckPortOpen(v4[0].String(), rhp4)
		}
		if params.HasAAAARecord {
			params.Rhp2v6, params.Rhp2v6Delay = nc.CheckPortOpen(v6[0].String(), rhp2)
			params.Rhp3v6, params.Rhp3v6Delay = nc.CheckPortOpen(v6[0].String(), rhp3)
			params.Rhp4v6, params.Rhp4v6Delay = nc.CheckPortOpen(v6[0].String(), rhp4)
		}

		if len(v4) > 0 {
			params.V4 = v4[0].String()
		} else {
			params.V4 = ""
		}
		if len(v6) > 0 {
			params.V6 = v6[0].String()
		} else {
			params.V6 = ""
		}
		sdk.UpdateCheck(params, wg, task)
		break

	case "IPv4":
		params.V4 = netAddress
		params.V6 = ""
		params.Rhp2v4, params.Rhp2v4Delay = nc.CheckPortOpen(netAddress, rhp2)
		params.Rhp3v4, params.Rhp3v4Delay = nc.CheckPortOpen(netAddress, rhp3)
		params.Rhp4v4, params.Rhp4v4Delay = nc.CheckPortOpen(netAddress, rhp4)
		sdk.UpdateCheck(params, wg, task)
		break

	case "IPv6":
		params.V4 = ""
		params.V6 = netAddress
		params.Rhp2v6, params.Rhp2v6Delay = nc.CheckPortOpen(netAddress, rhp2)
		params.Rhp3v6, params.Rhp3v6Delay = nc.CheckPortOpen(netAddress, rhp3)
		params.Rhp4v6, params.Rhp4v6Delay = nc.CheckPortOpen(netAddress, rhp4)
		sdk.UpdateCheck(params, wg, task)
		break
	}
}


func (nc *Checker) CheckVersion(publicKey string) (string, error) {
	host,err := explored.GetHostByPublicKey(publicKey)
	if err != nil {
		return "", err
	}
	return host.Settings.Version, nil	
}