package scan_test

import (
	"testing"

	"github.com/back2basic/siaalert/scanner/logger"
	"github.com/back2basic/siaalert/scanner/scan"
)

func TestDetectBadHost(t *testing.T) {
	// Reset the global map before starting tests
	scan.HostPortMap = make(map[string]map[string]bool)
	log := logger.GetLogger()
	defer logger.Sync()
	tests := []struct {
		netAddress     string
		expectedResult bool
	}{
		{"192.168.1.1:9982", false}, // Normal case
		{"192.168.1.1:9983", false}, // Different port, under threshold
		{"192.168.1.1:9984", false},
		{"192.168.1.1:9985", false},
		{"192.168.1.1:9986", false},
		{"192.168.1.1:9987", false},
		{"192.168.1.1:9988", false},
		{"192.168.1.1:9989", false},
		{"192.168.1.1:9990", false},
		{"192.168.1.1:9982", false},
		{"192.168.1.1:9482", true},
		{"192.168.1.1:9991", true},  // Exceeds threshold at this point
		{"192.168.1.1:9992", true},  // Exceeds threshold at this point
		{"192.168.1.2:9982", false}, // Different IP, always false
		{"192.168.1.3:9983", false}, // Another different IP
	}

	for _, test := range tests {
		result := scan.DetectBadHost(test.netAddress, log)
		if result != test.expectedResult {
			t.Errorf("DetectBadHost(%s): expected %v, got %v", test.netAddress, test.expectedResult, result)
		}
	}
}
