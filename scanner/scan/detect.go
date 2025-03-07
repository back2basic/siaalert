package scan

import (
	"net"

	"go.uber.org/zap"
)

var HostPortMap = make(map[string]map[string]bool)

// Helper function to get keys (ports) from a map (used for debug)
// func getKeys(portMap map[string]bool) []string {
// 	keys := make([]string, 0, len(portMap))
// 	for key := range portMap {
// 		keys = append(keys, key)
// 	}
// 	return keys
// }

func DetectBadHost(netAddress string, log *zap.Logger) bool {
	// Split host:port into IP and port
	host, port, err := net.SplitHostPort(netAddress)
	if err != nil {
		log.Error("Bad Host", zap.String("host", netAddress), zap.Error(err))
		return true
	}

	// host should never be
	if host == "[::1]" || host == "127.0.0.1" || host == "localhost" {
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
	if len(HostPortMap[host]) >= 5 {
		// fmt.Printf("Malicious IP detected: %s, Ports: %v\n", host, getKeys(HostPortMap[host]))
		return true
	}

	return false
}
