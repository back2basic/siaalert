package explored

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/back2basic/siadata/siaalert/config"
)

func GetConsensus() (Consensus, error) {
	cfg := config.GetConfig()
	url := cfg.External.ExploredUrl + "api/consensus/state"
	// Create the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	// // Add Basic Auth
	// auth := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	// req.Header.Add("Authorization", "Basic "+auth)

	// Send the request
	client := &http.Client{
		Timeout: 30 * time.Second,
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

func GetAllHosts() ([]Host, error) {
	var hosts []Host
	for i := 0; i < 1000; i++ {
		host, err := GetHosts(i * 500)
		if err != nil {
			return []Host{}, err
		}
		hosts = append(hosts, host...)
		if len(host) < 500 {
			break
		}
	}
	return hosts, nil
}

func GetHosts(offset int) ([]Host, error) {
	cfg := config.GetConfig()
	url := cfg.External.ExploredUrl + "api/hosts?limit=500&offset=" + fmt.Sprint(offset)
	reqBody := strings.NewReader(`{}`)
	// Create the request
	req, err := http.NewRequest("POST", url, reqBody)
	if err != nil {
		return []Host{}, fmt.Errorf("failed to create request: %v", err)
	}

	// Send the request
	client := &http.Client{
		Timeout: 60 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return []Host{}, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []Host{}, fmt.Errorf("failed to read response body: %v", err)
	}
	// fmt.Println(string(body))
	// Parse the JSON response
	var response []Host
	err = json.Unmarshal(body, &response)
	if err != nil {
		return []Host{}, fmt.Errorf("failed to parse JSON: %v", err)
	}
	return response, nil
}

func GetHostByPublicKey(publicKey string) (Host, error) {
	cfg := config.GetConfig()
	url := cfg.External.ExploredUrl + "api/hosts/"
	reqBody := strings.NewReader(`{"PublicKeys": "[` + publicKey + `]"}`)
	// Create the request
	req, err := http.NewRequest("GET", url, reqBody)
	if err != nil {
		return Host{}, fmt.Errorf("failed to create request: %v", err)
	}

	// Send the request
	client := &http.Client{
		Timeout: 60 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return Host{}, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Host{}, fmt.Errorf("failed to read response body: %v", err)
	}
	// fmt.Println(string(body))
	// Parse the JSON response
	var response Host
	err = json.Unmarshal(body, &response)
	if err != nil {
		return Host{}, fmt.Errorf("failed to parse JSON: %v", err)
	}
	return response, nil
}