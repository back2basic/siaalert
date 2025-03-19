package explored

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/back2basic/siaalert/control/config"
	"github.com/back2basic/siaalert/shared/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.sia.tech/coreutils/rhp/v4/siamux"
	"go.uber.org/zap"
)

var (
	ExploredCache = make(map[string]Host)
	Mutex         sync.RWMutex
)

func (h *Host) updateCacheIfDifferent(found Host) {
	if h != &found {
		// update cache
		Mutex.Lock()
		ExploredCache[h.PublicKey.String()] = *h
		Mutex.Unlock()
	}
}

// V2SiamuxAddr returns the `Address` of the first TCP siamux `NetAddress` it
// finds in the host's list of net addresses.  The protocol for this address is
// ProtocolTCPSiaMux.
func (h Host) V2SiamuxAddr() (string, bool) {
	for _, netAddr := range h.V2NetAddresses {
		if netAddr.Protocol == siamux.Protocol {
			return netAddr.Address, true
		}
	}
	return "", false
}

func (h *Host) ToBSON() bson.M {
	return bson.M{
		"publicKey":      h.PublicKey.String(),
		"v2":             h.V2,
		"netAddress":     h.NetAddress,
		"v2NetAddresses": h.V2NetAddresses,

		"knownSince":             h.KnownSince,
		"lastScan":               h.LastScan,
		"nextScan":               h.NextScan,
		"lastScanSuccessful":     h.LastScanSuccessful,
		"lastAnnouncement":       h.LastAnnouncement,
		"totalScans":             h.TotalScans,
		"successfulInteractions": h.SuccessfulInteractions,
		"failedInteractions":     h.FailedInteractions,

		"settings":   h.Settings,
		"priceTable": h.PriceTable,

		"rhpV4Settings": h.RHPV4Settings,
	}
}

func GetConsensus() (Consensus, error) {
	cfg := config.GetConfig()
	url := cfg.External.ExploredUrl + "api/consensus/tip"
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

func GetAllHosts(cached bool) (map[string]Host, error) {
	cfg := config.GetConfig()
	log := logger.GetLogger(cfg.Logging.Path)
	defer logger.Sync()
	if cached {
		if len(ExploredCache) != 0 {
			return ExploredCache, nil
		}
	}
	// clear cache
	// ExploredCache = make(map[string]Host)
	var hosts []Host
	// try grab new hosts
	for i := range 200 {
		host, err := GetHosts(i * 500)
		if err != nil {
			log.Error("Error getting hosts", zap.Error(err))
			return ExploredCache, err
		}
		hosts = append(hosts, host...)
		if len(host) < 500 {
			log.Warn("Finished getting hosts", zap.Int("hosts", len(hosts)))
			break
		}
	}
	// update cache
	for _, host := range hosts {
		Mutex.RLock()
		found, exists := ExploredCache[host.PublicKey.String()]
		Mutex.RUnlock()
		if exists {
			host.updateCacheIfDifferent(found)
			continue
		} else {
			Mutex.Lock()
			ExploredCache[host.PublicKey.String()] = host
			Mutex.Unlock()
		}
	}
	log.Info("Updated cache", zap.Int("hosts", len(hosts)))
	return ExploredCache, nil
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
	// try grabbing first from cache else grab from api
	host, exists := ExploredCache[publicKey]
	if exists {
		// fmt.Println("found in cache")
		return host, nil
	}
	// failback to api
	cfg := config.GetConfig()
	url := cfg.External.ExploredUrl + "api/hosts"
	reqBody := strings.NewReader(`{"PublicKeys": ["` + publicKey + `"]}`)
	// Create the request

	// fmt.Println(url)
	// fmt.Println(reqBody)

	req, err := http.NewRequest("POST", url, reqBody)
	if err != nil {
		return Host{}, fmt.Errorf("failed to create request: %v", err)
	}

	// fmt.Println(req)

	// Send the request
	client := &http.Client{
		Timeout: 60 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return Host{}, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// fmt.Println(resp.StatusCode)
	// fmt.Println(resp.Status)
	// Read the response body
	if resp.StatusCode != 200 {
		return Host{}, fmt.Errorf("response failed: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Host{}, fmt.Errorf("failed to read response body: %v", err)
	}
	// fmt.Println(string(body))
	// Parse the JSON response
	var response []Host
	err = json.Unmarshal(body, &response)
	if err != nil {
		return Host{}, fmt.Errorf("failed to parse JSON: %v", err)
	}
	if len(response) == 0 {
		return Host{}, fmt.Errorf("no host found")
	}
	// // add to cache
	Mutex.Lock()
	ExploredCache[publicKey] = response[0]
	Mutex.Unlock()
	return response[0], nil
}

