package bench_test

import (
	"net/http"
	"testing"

	"github.com/back2basic/siadata/siaalert/bench"
	"github.com/back2basic/siadata/siaalert/config"
	"github.com/stretchr/testify/assert"
)

// Mock HTTP client for testing
type MockClient struct {
	Response *http.Response
	Error    error
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return m.Response, m.Error
}

func TestScanHosts(t *testing.T) {
	// mockResponse := &http.Response{
	//     StatusCode: 200,
	//     Body:       io.NopCloser(strings.NewReader(`{"address": "123"}`)), // Example response body
	// }
	// mockClient := &MockClient{Response: mockResponse}

	// Set mock config
	config.LoadConfig("../config.yaml")

	// Replace HTTP client with mock
	// bench.Client = mockClient

	address := "123"
	hostkey := "123"
	response, err := bench.ScanHosts(address, hostkey)

	assert.NoError(t, err, "ScanHosts should not return an error")
	assert.Equal(t, "123", response.Settings.Netaddress, "Response field mismatch")
}

// func TestBenchHosts(t *testing.T) {
//     mockResponse := &http.Response{
//         StatusCode: 200,
//         Body:       io.NopCloser(strings.NewReader(`{"someField": "someValue"}`)), // Example response body
//     }
//     mockClient := &MockClient{Response: mockResponse}

//     // Set mock config
//     mockConfig := config.Config{
//         External: config.ExternalConfig{
//             BenchUrl: "http://mockurl.com/",
//         },
//     }
//     config.SetConfig(mockConfig)

//     // Replace HTTP client with mock
//     bench.Client = mockClient

//     address := "mockAddress"
//     hostkey := "mockHostKey"
//     sectors := 10
//     response, err := bench.BenchHosts(address, hostkey, sectors)

//     assert.NoError(t, err, "BenchHosts should not return an error")
//     assert.Equal(t, "someValue", response.SomeField, "Response field mismatch")
// }

func TestLoadBenchPeers(t *testing.T) {
	// mockResponse := &http.Response{
	//     StatusCode: 200,
	//     Body:       io.NopCloser(strings.NewReader(`[{"someField": "someValue"}]`)), // Example response body
	// }
	// mockClient := &MockClient{Response: mockResponse}

	// // Set mock config
	// mockConfig := config.Config{
	//     External: config.ExternalConfig{
	//         BenchUrl: "http://mockurl.com/",
	//     },
	// }
	// config.SetConfig(mockConfig)
	config.LoadConfig("../config.yaml")

	// // Replace HTTP client with mock
	// bench.Client = mockClient

	response, err := bench.LoadBenchPeers()

	assert.NoError(t, err, "LoadBenchPeers should not return an error")
	assert.Len(t, response, 1, "Response length mismatch")
	assert.Equal(t, "someValue", response[0].Address, "Response field mismatch")
}

func TestLoadBenchConsensus(t *testing.T) {
	// mockResponse := &http.Response{
	//     StatusCode: 200,
	//     Body:       io.NopCloser(strings.NewReader(`{"someField": "someValue"}`)), // Example response body
	// }
	// mockClient := &MockClient{Response: mockResponse}

	// // Set mock config
	// mockConfig := config.Config{
	//     External: config.ExternalConfig{
	//         BenchUrl: "http://mockurl.com/",
	//     },
	// }
	// config.SetConfig(mockConfig)
	config.LoadConfig("../config.yaml")

	// // Replace HTTP client with mock
	// bench.Client = mockClient

	response, err := bench.LoadBenchConsensus()

	assert.NoError(t, err, "LoadBenchConsensus should not return an error")
	assert.Equal(t, "someValue", response.ChainIndex.Height, "Response field mismatch")
}
