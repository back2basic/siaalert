package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/back2basic/siaalert/control/api"
	"github.com/back2basic/siaalert/control/config"
	"github.com/back2basic/siaalert/control/explored"
	"github.com/back2basic/siaalert/control/scan"
	"github.com/back2basic/siaalert/control/scheduler"
	"github.com/back2basic/siaalert/shared/logger"
	"go.mongodb.org/mongo-driver/bson"

	"go.uber.org/zap"
)

var runs = 0
var maxRhpWorkers = 60
var maxPortWorkers = 10
var scan5minRHPrunning bool
var scan5minPORTrunning bool

func processDuration(elapsed time.Duration, log *zap.Logger, workType string) {
	if elapsed > 6*time.Minute {
		log.Warn("Process took too long", zap.Duration("duration", elapsed))
		switch workType {
		case "rhp":
			if maxRhpWorkers < 60 {
				maxRhpWorkers++
			}
		case "portscan":
			if maxPortWorkers < 20 {
				maxPortWorkers++
			}
		}
	}
	if elapsed < 4*time.Minute {
		log.Warn("Process took too short", zap.Duration("duration", elapsed))
		switch workType {
		case "rhp":
			if maxRhpWorkers > 1 {
				maxRhpWorkers--
			}
		case "portscan":
			if maxPortWorkers > 1 {
				maxPortWorkers--
			}
		}
	}
}

func schedlPort(log *zap.Logger) {
	start := time.Now()
	cfg := config.GetConfig()

	hosts := make(map[string]explored.Host)
	err := error(nil)
	// Step 1: Fetch hosts
	hosts, err = fetchHosts(true)
	if err != nil {
		log.Error("Failed to fetch hosts:", zap.Error(err))
	}
	// Step 1.5: Filter hosts
	filterred := filterHosts(hosts, log)

	var rhpScans []scan.HostScan
	for _, host := range filterred {
		rhp := cfg.DB.FindRhp(bson.M{"publicKey": host.PublicKey.String()})
		if rhp.Err() != nil {
			log.Error("Failed to find rhp for portscan:", zap.Error(err))
			continue
		}
		var rhpScan scan.HostScan
		if err := rhp.Decode(&rhpScan); err != nil {
			log.Error("Failed to decode rhp for portscan:", zap.Error(err))
			continue
		}
		if (rhpScan.OnlineSince != time.Time{}) {
			// if (rhpScan.OnlineSince != time.Time{}) || (time.Since(rhpScan.OfflineSince) < 2*time.Hour) {
			rhpScans = append(rhpScans, rhpScan)
		}
	}

	if len(rhpScans) == 0 {
		log.Info("No hosts to portscan")
		return
	}
	log.Info("Starting portscan", zap.Int("hosts", len(rhpScans)), zap.Int("workers", maxPortWorkers))
	var wg sync.WaitGroup
    taskChannel := make(chan scan.HostScan, len(rhpScans))

    // Step 2: Start workers.
    for i := 0; i < maxRhpWorkers; i++ {
        wg.Add(1)
        go func(workerID int) {
            defer wg.Done()
            for host := range taskChannel {
                err := scanPort(host)
                if err != nil {
                    // fmt.Printf("Worker %d: Error portscanning host %s: %v\n", workerID, host.PublicKey, err)
                } else {
                    // fmt.Printf("Worker %d: Successfully portscanned host %s\n", workerID, host.PublicKey)
                }
            }
        }(i)
    }

    // Step 3: Add tasks to the channel.
    for _, host := range rhpScans {
        taskChannel <- host
    }
    close(taskChannel) // Close the channel to signal no more tasks.
		log.Info("Channel PortScan closed")

    // Step 4: Wait for all workers to finish.
    wg.Wait()
	// Step 2: Scan hosts concurrently
	// var wg sync.WaitGroup
	// sem := make(chan struct{}, maxRhpWorkers)

	// for _, scanHost := range rhpScans {
	// 	wg.Add(1)
	// 	sem <- struct{}{}
	// 	go func(s scan.HostScan) {
	// 		defer func() { <-sem }()
	// 		// scanHost(h, &wg, log, cfg.DB, checker)
	// 		err := scanPort(s, &wg)
	// 		if err != nil {
	// 			// log.Error("Failed to schedule portscan:", zap.Error(err))
	// 		}
	// 	}(scanHost)
	// }
	// Step 2: Schedule scans
	// for _, host := range filterred {

	// }
	
	//  if runtime is lower then 1 minutes sleep till 1 minute expires
	if time.Since(start) < time.Minute {
		time.Sleep(time.Minute - time.Since(start))
	}

	log.Info("Processed portScans in", zap.Duration("duration", time.Since(start)))
	processDuration(time.Since(start), log, "portscan")

}

func scanPort(host scan.HostScan) error {
	// defer wg.Done()
	cfg := config.GetConfig()
	url := cfg.External.ScannerUrl + "v1/scan/port/"
	reqBody := strings.NewReader(`{"PublicKey": "` + host.PublicKey + `"}`)
	// Create the request
	req, err := http.NewRequest("POST", url, reqBody)
	if err != nil {
		return err
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
		// fmt.Println("Error making request:", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("non-200 status code: %d", resp.StatusCode)
	}
	// // Read the response body
	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	// fmt.Printf("Error reading response: %s\n", err)
	// 	return err
	// }
	// // fmt.Println(string(body))
	// // Parse the JSON response
	// var response scan.Check
	// err = json.Unmarshal(body, &response)
	// if err != nil {
	// 	// fmt.Printf("Error parsing JSON: %s\n", err)
	// 	return err
	// }
	// // return response, nil
	return nil
}

func scanRhp(host explored.Host) error {
	// defer wg.Done()

	cfg := config.GetConfig()
	url := cfg.External.ScannerUrl + "v1/scan/rhp/"
	reqBody := strings.NewReader(`{"PublicKey": "` + host.PublicKey.String() + `"}`)
	// Create the request
	req, err := http.NewRequest("POST", url, reqBody)
	if err != nil {
		return err
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
		// fmt.Println("Error making request:", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("Failed to scan rhp: %s", resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// fmt.Printf("Error reading response: %s\n", err)
		return err
	}
	// fmt.Println(string(body))
	// Parse the JSON response
	var response scan.HostScan
	err = json.Unmarshal(body, &response)
	if err != nil {
		// fmt.Printf("Error parsing JSON: %s\n", err)
		return err
	}
	// return response, nil
	return nil
}

func fetchHosts(cache bool) (map[string]explored.Host, error) {
	return explored.GetAllHosts(cache)
}

func CheckVersion(publicKey string) (string, error) {
	host, err := explored.GetHostByPublicKey(publicKey)
	if err != nil {
		return "", err
	}
	return host.Settings.Version, nil
}

func filterHosts(hosts map[string]explored.Host, log *zap.Logger) []explored.Host {
	// checker := &scan.Checker{}
	needScanning := []explored.Host{}
	skipped := 0
	failed := 0
	malicious := 0
	hostd := 0
	siad := 0

	for _, host := range hosts {
		// Get netAddress
		var netAddresss string

		if host.V2 {
			addr, v2 := host.V2SiamuxAddr()
			if !v2 {
				log.Warn("Failed to get v2 siamux address", zap.String("host", host.PublicKey.String()))
				failed++
				continue
			}
			netAddresss = addr
		} else {
			netAddresss = host.NetAddress
		}

		// Bad Host detection
		if scan.DetectBadHost(netAddresss, log) {
			// log.Warn("Bad host detected", zap.Bool("v2", host.V2), zap.String("host", netAddresss), zap.String("publicKey", host.PublicKey.String()))
			malicious++
			continue
		}

		version, err := CheckVersion(host.PublicKey.String())
		if err != nil {
			// log.Warn("Failed to get version", zap.String("host", host.PublicKey.String()), zap.Error(err))
			failed++
			continue
		}

		// if hostd only check 1 year since last announcement
		// once hardfork is active in june we can skip v1
		if version != "" && version != "1.5.9" || host.V2 {
			if time.Since(host.LastAnnouncement).Hours() > (24 * 365) {
				skipped++
				continue
			} else {
				hostd++
			}
		} else {
			if host.TotalScans != 0 && host.FailedInteractions != 0 && host.SuccessfulInteractions == 0 {
				// log.Warn("Failed interactions", zap.Uint64("total", host.TotalScans), zap.Uint64("failed", host.FailedInteractions), zap.Uint64("successful", host.SuccessfulInteractions))
				failed++
				continue
			}
			siad++
		}

		if (host.NextScan != time.Time{}) && host.NextScan.After(time.Now()) {
			// log.Warn("Next scan is in the future", zap.Time("nextScan", host.NextScan), zap.String("host", host.PublicKey.String()))
			skipped++
			continue
		}
		// append to needscanning
		needScanning = append(needScanning, host)
	}

	log.Info("Filterred",
		zap.Int("scanning", len(needScanning)),
		zap.Int("skipped", skipped), zap.Int("failed", failed),
		zap.Int("malicious", malicious), zap.Int("hostd", hostd),
		zap.Int("siad", siad))

	return needScanning
}

func schedlRHP(log *zap.Logger) {
	defer func() {
		runs++
	}()
	start := time.Now()
	hosts := make(map[string]explored.Host)
	err := error(nil)
	// Step 1: Fetch hosts
	if runs < 10 {
		hosts, err = fetchHosts(true)
		if err != nil {
			log.Error("Failed to fetch hosts:", zap.Error(err))
		}
	} else {
		runs = 0
		hosts, err = fetchHosts(false)
		if err != nil {
			log.Error("Failed to fetch hosts:", zap.Error(err))
		}
	}
	// Step 1.5: Filter hosts
	filtered := filterHosts(hosts, log)
	// maxWorkers := max(min(len(filterred)/50, 5), 1)
	// maxWorkers := 2
	log.Info("Starting scan", zap.Int("workers", maxRhpWorkers), zap.Int("run", runs))

	// // Step 2: Scan hosts concurrently
	var wg sync.WaitGroup
    taskChannel := make(chan explored.Host, len(filtered))

    // Step 2: Start workers.
    for i := 0; i < maxRhpWorkers; i++ {
        wg.Add(1)
        go func(workerID int) {
            defer wg.Done()
            for host := range taskChannel {
                err := scanRhp(host)
                if err != nil {
                    // fmt.Printf("Worker %d: Error scanning host %s: %v\n", workerID, host.PublicKey, err)
                } else {
                    // fmt.Printf("Worker %d: Successfully scanned host %s\n", workerID, host.PublicKey)
                }
            }
        }(i)
    }

    // Step 3: Add tasks to the channel.
    for _, host := range filtered {
        taskChannel <- host
    }
    close(taskChannel) // Close the channel to signal no more tasks.
	log.Info("Channel RHP closed")
    // Step 4: Wait for all workers to finish.
    wg.Wait()

	// var wg sync.WaitGroup
	// sem := make(chan struct{}, maxRhpWorkers)

	// for _, host := range filterred {
	// 	wg.Add(1)
	// 	sem <- struct{}{}
	// 	go func(h explored.Host) {
	// 		defer func() { <-sem }()
	// 		// scanHost(h, &wg, log, cfg.DB, checker)
	// 		err := scanRhp(h, &wg)
	// 		if err != nil {
	// 			// log.Error("control RHP", zap.Error(err))
	// 		}
	// 	}(host)
	// }
	// close(sem) // Close the semaphore
	// wg.Wait()  // Wait for all scans to finish

	//  if runtime is lower then 1 minutes sleep till 1 minute expires
	if time.Since(start) < time.Minute {
		time.Sleep(time.Minute - time.Since(start))
	}

	log.Info("Processed rhpScans in", zap.Duration("duration", time.Since(start)))
	processDuration(time.Since(start), log, "rhp")
}

func main() {
	// Handle signals for graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		// log.Warn("Shutting down...", zap.String("module", "api"), zap.String("module", "scanner"))
		os.Exit(0)
	}()

	// Load the configuration
	cfg := config.LoadConfig("./data/config-control.yaml")

	// Initialize the logger
	log := logger.GetLogger(cfg.Logging.Path)
	defer logger.Sync()

	// checker := &scan.Checker{}
	defer cfg.Close(log)

	// Initialize the scheduler
	sched := scheduler.NewScheduler()

	// Add a cron job to run every 5 minutes for RHP scanning
	_, err := sched.AddJob("30 */5 * * * *", func() {
		if !scan5minRHPrunning {
			scan5minRHPrunning = true
			schedlRHP(log)
			scan5minRHPrunning = false
		}
	})
	if err != nil {
		log.Error("Failed to add job:", zap.Error(err))
	}

	// Add a cron job to run every 5 minutes for PortScan
	_, err = sched.AddJob("* */5 * * * *", func() {
		if !scan5minPORTrunning && runs != 0 {
			scan5minPORTrunning = true
			schedlPort(log)
			scan5minPORTrunning = false
		}
	})
	if err != nil {
		log.Error("Failed to add job:", zap.Error(err))
	}

	// Start the scheduler
	sched.Start()

	// Start API server
	log.Info("API is starting...", zap.String("module", "scanner"))
	api.StartServer(log)

}
