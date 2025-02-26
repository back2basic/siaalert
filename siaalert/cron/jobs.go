package cron

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/back2basic/siadata/siaalert/explored"
	"github.com/back2basic/siadata/siaalert/scan"
	"github.com/back2basic/siadata/siaalert/sdk"
)

func CheckNewExporedHosts(hosts []explored.Host) {
	for i := 1; i <= len(hosts); i++ {
		sdk.Mutex.Lock()
		checked, exists := sdk.HostCache[hosts[i-1].PublicKey]
		sdk.Mutex.Unlock()

		// Check host if not already cached
		if !exists {
			var err error
			checked, err = sdk.CheckHost(hosts[i-1])
			if err != nil {
				continue
			}

			sdk.Mutex.Lock()
			sdk.HostCache[hosts[i-1].PublicKey] = checked
			sdk.Mutex.Unlock()
		} else {
			if checked.NetAddress != hosts[i-1].NetAddress {
				checked.NetAddress = hosts[i-1].NetAddress
				sdk.Mutex.Lock()
				sdk.HostCache[hosts[i-1].PublicKey] = checked
				sdk.Mutex.Unlock()
				go sdk.UpdateNetAddress(checked)
			}
		}
	}
}

func RunScan(hosts map[string]sdk.HostDocument, checker scan.Checker) {
	needScanning := []sdk.HostDocument{}
	skipped := 0
	failed := 0
	for _, host := range hosts {
		lastAnnounced, err := time.Parse(time.RFC3339, host.LastAnnouncement)
		if err != nil {
			failed++
			continue
		}

		if !host.Online && host.Error != "" && host.OfflineSince != "" && time.Since(lastAnnounced).Hours() > (24*365*2) {
			skipped++
			continue
		}
		// append to needscanning
		needScanning = append(needScanning, host)
	}

	if len(needScanning) == 0 {
		return
	}
	// Queue
	jobQueue := make(chan Job, len(needScanning))
	var wg sync.WaitGroup

	// Workers
	numWorkers := 20
	if os.Getenv("NETWORK") == "main" {
		if len(needScanning)/5 > 100 {
			numWorkers = 100
		} else {
			numWorkers = len(needScanning) / 5
		}
	}
	for i := 1; i <= numWorkers; i++ {
		worker := NewWorker(i, jobQueue, &wg)
		worker.Start(checker)
	}
	// Add jobs to the queue
	var jobId int
	for _, host := range needScanning {
		// schedule scan
		jobId++
		jobscan := Job{
			ID:      jobId,
			Type:    "scan",
			Name:    host.Id,
			Address: host.NetAddress,
			HostKey: host.PublicKey,
		}
		wg.Add(1)
		jobQueue <- jobscan
	}

	fmt.Printf("Skipped %d hosts\n", skipped)
	fmt.Printf("Failed %d hosts\n", failed)
	fmt.Printf("Scanning %d hosts\n", len(hosts)-skipped-failed)

	close(jobQueue)
	wg.Wait()
}

func RunBench(hosts []explored.Host, checker scan.Checker) {
	// Queue
	jobQueue := make(chan Job, 2000)
	var wg sync.WaitGroup

	numWorkers := 5
	for i := 1; i <= numWorkers; i++ {
		worker := NewWorker(i, jobQueue, &wg)
		worker.Start(checker)
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

func RunRhp(hosts []explored.Host) {
	for _, host := range hosts {
		if time.Since(host.LastAnnouncement).Hours() > 24*365*2 {
			sdk.UpdateRhp(host)
		}
	}
}
