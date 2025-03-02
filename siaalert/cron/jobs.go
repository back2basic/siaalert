package cron

import (
	"fmt"
	"sync"
	"time"

	"github.com/back2basic/siadata/siaalert/explored"
	"github.com/back2basic/siadata/siaalert/scan"
	"github.com/back2basic/siadata/siaalert/sdk"
	"github.com/back2basic/siadata/siaalert/strict"
)

func CheckNewExploredHosts(host explored.Host) {
	sdk.Mutex.RLock()
	checked, exists := sdk.HostCache[host.PublicKey]
	sdk.Mutex.RUnlock()

	// Check host if not already cached
	if !exists {
		var err error
		checked, err = sdk.CheckHost(host)
		if err != nil {
			return
		}

		sdk.Mutex.Lock()
		sdk.HostCache[host.PublicKey] = checked
		sdk.Mutex.Unlock()
	} else {
		if checked.NetAddress != host.NetAddress {
			checked.NetAddress = host.NetAddress
			sdk.Mutex.Lock()
			sdk.HostCache[host.PublicKey] = checked
			sdk.Mutex.Unlock()
			sdk.UpdateNetAddress(checked)
		}
	}
}

func RunScan(hosts map[string]strict.HostDocument, checker *scan.Checker) {
	needScanning := []strict.HostDocument{}
	skipped := 0
	failed := 0
	malicious := 0
	hostd := 0
	siad := 0
	fmt.Println("Scanning", len(hosts), "hosts")
	// Filter hosts
	for _, host := range hosts {
		// fmt.Println("Scanning host", host.NetAddress)
		lastAnnounced, err := time.Parse(time.RFC3339, host.LastAnnouncement)
		if err != nil {
			failed++
			continue
		}

		// Bad Host detection
		if scan.DetectBadHost(host.NetAddress) {
			malicious++
			continue
		}

		version, err := checker.CheckVersion(host.PublicKey)
		if err != nil {
			failed++
			continue
		}

		// if hostd only check 1 year since last announcement
		if version == "1.6.0" {
			if time.Since(lastAnnounced).Hours() > (24 * 183) {
				skipped++
				continue
			} else {
				hostd++
			}
		} else {
			if host.TotalScans != 0 && host.FailedInteractions != 0 {
				failed++
				continue
			}
			siad++
		}

		// append to needscanning
		needScanning = append(needScanning, host)
	}

	if len(needScanning) == 0 {
		Running5Minutes = false
		return
	}
	// Workers max 400 min 2
	numWorkers := max(min(len(needScanning)/20, 100), 1)

	fmt.Println("Starting", numWorkers, "workers for scanning", len(needScanning), "hosts")
	fmt.Printf("Skipped %d hosts\n", skipped)
	fmt.Printf("Failed %d hosts\n", failed)
	fmt.Printf("Malicious %d hosts\n", malicious)
	fmt.Printf("Hostd %d hosts\n", hostd)
	fmt.Printf("Siad %d hosts\n", siad)
	fmt.Printf("Scanning %d hosts\n", len(hosts)-skipped-failed-malicious)

	// Queue
	jobQueue := make(chan Job, len(needScanning))
	var wg sync.WaitGroup

	for i := 1; i <= numWorkers; i++ {
		go worker(i, jobQueue, &wg, *checker)

	}
	// Add jobs to the queue
	// var jobId int
	for j, host := range needScanning {
		// schedule scan
		// jobId++
		jobscan := Job{
			ID:      j,
			Type:    "scan",
			Name:    host.Id,
			Address: host.NetAddress,
			HostKey: host.PublicKey,
			V2:      host.V2,
		}
		wg.Add(1)
		jobQueue <- jobscan
	}

	close(jobQueue)
	// fmt.Println("Closed jobQueue")
	wg.Wait()
	// fmt.Println("wg done waiting")
}

func RunBench(hosts []explored.Host, checker scan.Checker) {
	// Queue
	jobQueue := make(chan Job, 2000)
	var wg sync.WaitGroup

	numWorkers := 5
	for i := 1; i <= numWorkers; i++ {
		go worker(i, jobQueue, &wg, checker)
	}

	// Add jobs to the queue
	for i := 1; i <= len(hosts); i++ {
		jobbench := Job{
			ID:      i,
			Type:    "bench",
			Name:    hosts[i-1].NetAddress,
			Address: hosts[i-1].NetAddress,
			HostKey: hosts[i-1].PublicKey,
		}
		wg.Add(1)
		jobQueue <- jobbench
	}

	close(jobQueue)
	wg.Wait()

}

func RunRhp(hosts map[string]explored.Host) {
	for _, host := range hosts {
		if time.Since(host.LastAnnouncement).Hours() > 24*365 {
			sdk.UpdateRhp(host)
		}
	}
}
