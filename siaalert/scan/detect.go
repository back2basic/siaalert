package scan

import (
	"fmt"
	"net"
)

var HostPortMap = make(map[string]map[string]bool)

// Helper function to get keys (ports) from a map
func getKeys(portMap map[string]bool) []string {
	keys := make([]string, 0, len(portMap))
	for key := range portMap {
		keys = append(keys, key)
	}
	return keys
}


func DetectBadHost(netAddress string) bool {
	// Split host:port into IP and port
	host, port, err := net.SplitHostPort(netAddress)
	if err != nil {
			fmt.Printf("Invalid host:port format: %s\n", netAddress)
			return true
	}

	// Initialize map for the IP if not already present
	if _, exists := HostPortMap[host]; !exists {
			HostPortMap[host] = make(map[string]bool)
	}

	// Check if the port is already tracked
	if !HostPortMap[host][port] {
			HostPortMap[host][port] = true
	}

	// Debug log: show current state
	// fmt.Printf("Host: %s, Current Ports: %v\n", host, getKeys(HostPortMap[host]))

	// Check threshold
	if len(HostPortMap[host]) >= 10 {
			// fmt.Printf("Malicious IP detected: %s, Ports: %v\n", host, getKeys(HostPortMap[host]))
			return true
	}

	return false
}
