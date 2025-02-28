package bench

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/back2basic/siadata/siaalert/config"
)

func ScanHosts(address, hostkey string) (Scan, error) {
	cfg := config.GetConfig()
	url := cfg.External.BenchUrl + "scan"
	// Create the request
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return Scan{}, fmt.Errorf("%v", err)
	}

	// Add Body
	req.Body = io.NopCloser(strings.NewReader(fmt.Sprintf(`{"address": "%s", "hostkey": "%s"}`, address, hostkey)))

	// Send the request
	client := &http.Client{
		Timeout: 15 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return Scan{}, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Scan{}, fmt.Errorf("%v", err)
	}

	if resp.StatusCode != 200 {
		return Scan{}, fmt.Errorf("%s", string(body))
	}

	// Parse the JSON response
	var response Scan
	err = json.Unmarshal(body, &response)
	if err != nil {
		return Scan{}, fmt.Errorf("%v", err)
	}

	return response, nil
}

func BenchHosts(address, hostkey string, sectors int) (Bench, error) {
	cfg := config.GetConfig()
	url := cfg.External.BenchUrl + "benchmark"
	// Create the request
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return Bench{}, fmt.Errorf("failed to create request: %v", err)
	}

	// Add Body
	req.Body = io.NopCloser(strings.NewReader(fmt.Sprintf(`{"address": "%s", "hostkey": "%s", "sectors": %d}`, address, hostkey, sectors)))

	// Send the request
	client := &http.Client{
		// Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return Bench{}, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Bench{}, fmt.Errorf("failed to read response body: %v", err)
	}

	fmt.Println(string(body))

	// Parse the JSON response
	var response Bench
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println(response)
		return Bench{}, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	// fmt.Println(string(body))
	return response, nil

}

func LoadBenchPeers() ([]Peers, error) {
	cfg := config.GetConfig()
	url := cfg.External.BenchUrl + "syncer/peers"
	// Create the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	// Send the request
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %s\n", err)
		return nil, err
	}
	// fmt.Println(string(body))
	// Parse the JSON response
	var response []Peers
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("Error parsing JSON: %s\n", err)
		return nil, err
	}
	return response, nil
}

func LoadBenchConsensus() (Consensus, error) {
	cfg := config.GetConfig()
	url := cfg.External.BenchUrl + "state/consensus"
	// Create the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	// Send the request
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return Consensus{}, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %s\n", err)
		return Consensus{}, err
	}
	// fmt.Println(string(body))
	// Parse the JSON response
	var response Consensus
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("Error parsing JSON: %s\n", err)
		return Consensus{}, err
	}
	return response, nil
}
