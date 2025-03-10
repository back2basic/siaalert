package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/back2basic/siaalert/scanner/api"
	"github.com/back2basic/siaalert/scanner/config"
	"github.com/back2basic/siaalert/scanner/db"
	"github.com/back2basic/siaalert/scanner/explored"
	"github.com/back2basic/siaalert/scanner/logger"
	"github.com/back2basic/siaalert/scanner/mail"
	"github.com/back2basic/siaalert/scanner/scan"
	"github.com/pingcap/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

func checkRhpResult(mongoDB *db.MongoDB, netAddress string, online bool, result scan.HostScan, log *zap.Logger) {
	prev := mongoDB.FindRhp(bson.M{"publicKey": result.PublicKey})
	if prev.Err() != nil {
		err := mongoDB.UpdateRhp(result.PublicKey, online, result.ToBSON(), log)
		if err != nil {
			log.Error("Failed to create rhp", zap.Error(err))
		}
		return
		// log.Warn("Failed to find document", zap.Error(prev.Err()))
		// if result.Error != "" {
		// 	result.OfflineSince = time.Now()
		// } else {
		// 	result.OnlineSince = time.Now()
		// }
		// return err
	}
	var prevScan scan.HostScan
	if err := prev.Decode(&prevScan); err != nil {
		log.Warn("Failed to decode rhp", zap.Error(err))
		// return err
	}

	result.OnlineSince = prevScan.OnlineSince
	result.OfflineSince = prevScan.OfflineSince

	if prevScan.Success != online {
		if online {
			log.Info("Host " + result.PublicKey + " is online")
			result.OnlineSince = time.Now()
			result.OfflineSince = time.Time{}
			// Send Mail
			mail.PrepareAlertEmails(netAddress, "Online", result.PublicKey, log, mongoDB)
		} else {
			log.Warn("Host " + result.PublicKey + " is offline")
			result.OnlineSince = time.Time{}
			result.OfflineSince = time.Now()
			// Send Mail
			mail.PrepareAlertEmails(netAddress, "Offline", result.PublicKey, log, mongoDB)
		}
	}
	// result.NetAddress = netAddress
	err := mongoDB.UpdateRhp(result.PublicKey, online, result.ToBSON(), log)
	if err != nil {
		log.Error("Failed to update rhp", zap.Error(err))
	}
}

func scanHost(host explored.Host, wg *sync.WaitGroup, log *zap.Logger, mongodDB *db.MongoDB, checker *scan.Checker) {
	defer wg.Done() // Signal the WaitGroup when this goroutine is done
	mongodDB.UpdateHost(
		host.PublicKey,
		host.ToBSON(),
	)
	// log.Info("Scanning host", zap.String("host", host.PublicKey.String()))
	scanned, err := scan.RunRhpScan(host, log, mongodDB, checker)
	if err != nil {
		scanned.Error = err.Error()
		checkRhpResult(mongodDB, host.NetAddress, false, scanned, log)
		return
	}

	checkRhpResult(mongodDB, host.NetAddress, true, scanned, log)
	checker.PortScan(host.PublicKey, scanned, mongodDB)
}

// Fetch the list of hosts (replace with your actual fetching logic)
func fetchHosts(cache bool) (map[string]explored.Host, error) {
	return explored.GetAllHosts(cache)
}

func filterHosts(hosts map[string]explored.Host, log *zap.Logger) []explored.Host {
	checker := &scan.Checker{}
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

		version, err := checker.CheckVersion(host.PublicKey.String())
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

func main() {
	// Handle signals for graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		log.Warn("Shutting down...", zap.String("module", "api"), zap.String("module", "scanner"))
		os.Exit(0)
	}()

	// Initialize the logger
	log := logger.GetLogger()
	defer logger.Sync()

	// Load the configuration
	cfg := config.LoadConfig("./data/config-scanner.yaml")
	checker := &scan.Checker{}
	defer cfg.Close(log)

	// Start API server
	go api.StartServer(log)
	log.Info("API is starting...", zap.String("module", "scanner"))

	workTime := 10
	// Start scanning
	runs := 0
	for {
		start := time.Now()
		hosts := make(map[string]explored.Host)
		err := error(nil)
		// Step 1: Fetch hosts
		if runs != 10 {
			runs++
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
		filterred := filterHosts(hosts, log)
		maxWorkers := max(min(len(filterred)/workTime, 100), 1)

		log.Info("Starting scan", zap.Int("workers", maxWorkers), zap.Int("run", runs))

		// Step 2: Scan hosts concurrently
		var wg sync.WaitGroup
		sem := make(chan struct{}, maxWorkers * 10)

		for _, host := range filterred {
			wg.Add(1)
			sem <- struct{}{}
			go func(h explored.Host) {
				defer func() { <-sem }()
				scanHost(h, &wg, log, cfg.DB, checker)
			}(host)
		}
		close(sem) // Close the semaphore
		wg.Wait()  // Wait for all scans to finish
		log.Info("All hosts scanned.")

		// Step 2.9: adjust workers totaltime
		elapsed := time.Since(start)
		if elapsed > 6*time.Minute {
			workTime--
			if workTime < 1 {
				workTime = 1
			}
		}
		if elapsed < 4*time.Minute {
			workTime++
		}
		// Step 3: Wait for the next 5-minute interval
		if elapsed < 5*time.Minute {
			log.Info("Waiting for the next 5-minute interval...", zap.Duration("remaining", 5*time.Minute-elapsed))
			time.Sleep(5*time.Minute - elapsed)
		} else {
			log.Warn("Next run should have been started already...", zap.Duration("late", 5*time.Minute-elapsed))
		}
	}
}
