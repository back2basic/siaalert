package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/back2basic/siaalert/scanner/api"
	"github.com/back2basic/siaalert/scanner/config"
	"github.com/back2basic/siaalert/shared/logger"

	"go.uber.org/zap"
)

// func checkRhpResult(mongoDB *db.MongoDB, netAddress string, online bool, result scan.HostScan, log *zap.Logger) {
// 	prev := mongoDB.FindRhp(bson.M{"publicKey": result.PublicKey})
// 	if prev.Err() != nil {
// 		err := mongoDB.UpdateRhp(result.PublicKey, online, result.ToBSON(), log)
// 		if err != nil {
// 			log.Error("Failed to create rhp", zap.Error(err))
// 		}
// 		return
// 	}

// 	var prevScan scan.HostScan
// 	if err := prev.Decode(&prevScan); err != nil {
// 		log.Warn("Failed to decode rhp", zap.Error(err))
// 		// return err
// 	}

// 	result.OnlineSince = prevScan.OnlineSince
// 	result.OfflineSince = prevScan.OfflineSince

// 	if prevScan.Success != online {
// 		if online {
// 			log.Info("Host " + result.PublicKey + " is online")
// 			result.OnlineSince = time.Now()
// 			result.OfflineSince = time.Time{}
// 			// Send Mail
// 			// mail.PrepareAlertEmails(netAddress, "Online", result.PublicKey, log, mongoDB)
// 		} else {
// 			log.Warn("Host " + result.PublicKey + " is offline")
// 			result.OnlineSince = time.Time{}
// 			result.OfflineSince = time.Now()
// 			// Send Mail
// 			// mail.PrepareAlertEmails(netAddress, "Offline", result.PublicKey, log, mongoDB)
// 		}
// 	}
// 	// result.NetAddress = netAddress
// 	err := mongoDB.UpdateRhp(result.PublicKey, online, result.ToBSON(), log)
// 	if err != nil {
// 		log.Error("Failed to update rhp", zap.Error(err))
// 	}
// 	log.Info("Finished checking RHP", zap.String("publicKey", result.PublicKey))
// }

// func scanHost(host explored.Host, wg *sync.WaitGroup, log *zap.Logger, mongodDB *db.MongoDB, checker *scan.Checker) {
// 	defer wg.Done()

// 	// mongodDB.UpdateHost(
// 	// 	host.PublicKey,
// 	// 	host.ToBSON(),
// 	// )

// 	// log.Info("Scanning host", zap.String("host", host.PublicKey.String()))
// 	scanned, err := scan.RunRhpScan(host, log, checker)
// 	if err != nil {
// 		// log.Error("RunRhpScan", zap.Error(err))
// 		scanned.Error = err.Error()
// 		checkRhpResult(mongodDB, host.NetAddress, false, scanned, log)
// 	} else {
// 		// log.Info("RunRhpScan", zap.String("host", host.PublicKey.String()))
// 		checkRhpResult(mongodDB, host.NetAddress, true, scanned, log)
// 	}
// 	// if (scanned.OnlineSince != time.Time{}) || (time.Since(scanned.OfflineSince) < 2*time.Hour) {
// 	// 	checker.PortScan(host.PublicKey, scanned, mongodDB)
// 	// }
// 	// log.Info("Finished scanning host", zap.String("host", host.PublicKey.String()))
// }

// // Fetch the list of hosts (replace with your actual fetching logic)
// func fetchHosts(cache bool) (map[string]explored.Host, error) {
// 	return explored.GetAllHosts(cache)
// }

// func filterHosts(hosts map[string]explored.Host, log *zap.Logger) []explored.Host {
// 	checker := &scan.Checker{}
// 	needScanning := []explored.Host{}
// 	skipped := 0
// 	failed := 0
// 	malicious := 0
// 	hostd := 0
// 	siad := 0

// 	for _, host := range hosts {
// 		// Get netAddress
// 		var netAddresss string

// 		if host.V2 {
// 			addr, v2 := host.V2SiamuxAddr()
// 			if !v2 {
// 				log.Warn("Failed to get v2 siamux address", zap.String("host", host.PublicKey.String()))
// 				failed++
// 				continue
// 			}
// 			netAddresss = addr
// 		} else {
// 			netAddresss = host.NetAddress
// 		}

// 		// Bad Host detection
// 		if scan.DetectBadHost(netAddresss, log) {
// 			// log.Warn("Bad host detected", zap.Bool("v2", host.V2), zap.String("host", netAddresss), zap.String("publicKey", host.PublicKey.String()))
// 			malicious++
// 			continue
// 		}

// 		version, err := checker.CheckVersion(host.PublicKey.String())
// 		if err != nil {
// 			// log.Warn("Failed to get version", zap.String("host", host.PublicKey.String()), zap.Error(err))
// 			failed++
// 			continue
// 		}

// 		// if hostd only check 1 year since last announcement
// 		// once hardfork is active in june we can skip v1
// 		if version != "" && version != "1.5.9" || host.V2 {
// 			if time.Since(host.LastAnnouncement).Hours() > (24 * 365) {
// 				skipped++
// 				continue
// 			} else {
// 				hostd++
// 			}
// 		} else {
// 			if host.TotalScans != 0 && host.FailedInteractions != 0 && host.SuccessfulInteractions == 0 {
// 				// log.Warn("Failed interactions", zap.Uint64("total", host.TotalScans), zap.Uint64("failed", host.FailedInteractions), zap.Uint64("successful", host.SuccessfulInteractions))
// 				failed++
// 				continue
// 			}
// 			siad++
// 		}
// 		// append to needscanning
// 		needScanning = append(needScanning, host)
// 	}

// 	log.Info("Filterred",
// 		zap.Int("scanning", len(needScanning)),
// 		zap.Int("skipped", skipped), zap.Int("failed", failed),
// 		zap.Int("malicious", malicious), zap.Int("hostd", hostd),
// 		zap.Int("siad", siad))

// 	return needScanning
// }

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
	cfg := config.LoadConfig("./data/config-scanner.yaml")
	
	// Initialize the logger
	log := logger.GetLogger(cfg.Logging.Path)
	defer logger.Sync()
	
	// checker := &scan.Checker{}
	defer cfg.Close(log)

	// Start API server
	log.Info("API is starting...", zap.String("module", "scanner"))
	api.StartServer(log)

	// workTime := 25
	// // Start scanning
	// runs := 0
	// for {
	// 	start := time.Now()
	// 	hosts := make(map[string]explored.Host)
	// 	err := error(nil)
	// 	// Step 1: Fetch hosts
	// 	if runs != 10 {
	// 		runs++
	// 		hosts, err = fetchHosts(true)
	// 		if err != nil {
	// 			log.Error("Failed to fetch hosts:", zap.Error(err))
	// 		}
	// 	} else {
	// 		runs = 0
	// 		hosts, err = fetchHosts(false)
	// 		if err != nil {
	// 			log.Error("Failed to fetch hosts:", zap.Error(err))
	// 		}
	// 	}
	// 	// Step 1.5: Filter hosts
	// 	filterred := filterHosts(hosts, log)
	// 	maxWorkers := max(min(len(filterred)/workTime, 150), 1)

	// 	log.Info("Starting scan", zap.Int("workers", maxWorkers), zap.Int("run", runs))

	// 	// Step 2: Scan hosts concurrently
	// 	var wg sync.WaitGroup
	// 	sem := make(chan struct{}, workTime*4)

	// 	for _, host := range filterred {
	// 		wg.Add(1)
	// 		sem <- struct{}{}
	// 		go func(h explored.Host) {
	// 			defer func() { <-sem }()
	// 			scanHost(h, &wg, log, cfg.DB, checker)
	// 		}(host)
	// 	}
	// 	close(sem) // Close the semaphore
	// 	wg.Wait()  // Wait for all scans to finish
	// 	log.Info("All hosts scanned.")

	// 	// Step 2.9: adjust workers totaltime
	// 	elapsed := time.Since(start)
	// 	if elapsed > 6*time.Minute {
	// 		workTime--
	// 		if workTime < 1 {
	// 			workTime = 1
	// 		}
	// 	}
	// 	if elapsed < 4*time.Minute {
	// 		workTime++
	// 	}
	// 	// Step 3: Wait for the next 15-minute interval
	// 	if elapsed < 15*time.Minute {
	// 		log.Info("Waiting for the next 5-minute interval...", zap.Duration("remaining", 15*time.Minute-elapsed))
	// 		time.Sleep(15*time.Minute - elapsed)
	// 	} else {
	// 		log.Warn("Next run should have been started already...", zap.Duration("late", 15*time.Minute-elapsed))
	// 	}
	// }
}
