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

		// udpate release
		// if exists {
		// 	found, err := sdk.GetCheck(checked.Id)
		// 	if err != nil {
		// 		fmt.Println("error updating release", err, checked.NetAddress)
		// 		continue
		// 	}
		// 	var checkDoc sdk.CheckDocument
		// 	found.Decode(&checkDoc)
		// 	params := sdk.Check{
		// 		HostId:             checkDoc.HostId,
		// 		V4Addr:             checkDoc.V4Addr,
		// 		V6Addr:             checkDoc.V6Addr,
		// 		Rhp2Port:           checkDoc.Rhp2Port,
		// 		Rhp2V4Delay:        checkDoc.Rhp2V4Delay,
		// 		Rhp2V6Delay:        checkDoc.Rhp2V6Delay,
		// 		Rhp2V4:             checkDoc.Rhp2V4,
		// 		Rhp2V6:             checkDoc.Rhp2V6,
		// 		Rhp3Port:           checkDoc.Rhp3Port,
		// 		Rhp3V4:             checkDoc.Rhp3V4,
		// 		Rhp3V6:             checkDoc.Rhp3V6,
		// 		Rhp3V4Delay:        checkDoc.Rhp3V4Delay,
		// 		Rhp3V6Delay:        checkDoc.Rhp3V6Delay,
		// 		Rhp4Port:           checkDoc.Rhp4Port,
		// 		Rhp4V4:             checkDoc.Rhp4V4,
		// 		Rhp4V6:             checkDoc.Rhp4V6,
		// 		Rhp4V4Delay:        checkDoc.Rhp4V4Delay,
		// 		Rhp4V6Delay:        checkDoc.Rhp4V6Delay,
		// 		AcceptingContracts: checkDoc.AcceptingContracts,
		// 		Release:            hosts[i-1].Settings.Release,
		// 	}
		// 	sdk.UpdateRelease(checkDoc.Id, params)
		// }

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

// func runPortScan(hostId string, scanned bench.Scan) {
// 	checker := scan.Checker{}
// 	netAddress, rhp2, err := checker.SplitAddressPort(scanned.Settings.Netaddress)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	rhp3 := scanned.Settings.Siamuxport
// 	rhp4 := "9984"

// 	// clasify netaddress
// 	var v4, v6 []net.IP = nil, nil
// 	params := sdk.CheckParams{}
// 	params.HostId = hostId
// 	classify := checker.ClassifyNetAddress(netAddress)
// 	switch classify {
// 	case "Hostname":
// 		params.HasARecord, params.HasAAAARecord, v4, v6 = checker.CheckDNSRecords(netAddress)
// 		if params.HasARecord {
// 			params.Rhp2v4, params.Rhp2v4Delay = checker.CheckPortOpen(v4[0].String(), rhp2)
// 			params.Rhp3v4, params.Rhp3v4Delay = checker.CheckPortOpen(v4[0].String(), rhp3)
// 			params.Rhp4v4, params.Rhp4v4Delay = checker.CheckPortOpen(v4[0].String(), rhp4)
// 		}
// 		if params.HasAAAARecord {
// 			params.Rhp2v6, params.Rhp2v6Delay = checker.CheckPortOpen(v6[0].String(), rhp2)
// 			params.Rhp3v6, params.Rhp3v6Delay = checker.CheckPortOpen(v6[0].String(), rhp3)
// 			params.Rhp4v6, params.Rhp4v6Delay = checker.CheckPortOpen(v6[0].String(), rhp4)
// 		}

// 		if len(params.V4) > 0 {
// 			params.V4 = v4[0].String()
// 		}
// 		if len(params.V6) > 0 {
// 			params.V4 = v6[0].String()
// 		}
// 		sdk.UpdateCheck(params)
// 		break

// 	case "IPv4":
// 		params.Rhp2v4, params.Rhp2v4Delay = checker.CheckPortOpen(netAddress, rhp2)
// 		params.Rhp3v4, params.Rhp3v4Delay = checker.CheckPortOpen(netAddress, rhp3)
// 		params.Rhp4v4, params.Rhp4v4Delay = checker.CheckPortOpen(netAddress, rhp4)
// 		sdk.UpdateCheck(params)
// 		break

// 	case "IPv6":
// 		params.Rhp2v6, params.Rhp2v6Delay = checker.CheckPortOpen(netAddress, rhp2)
// 		params.Rhp3v6, params.Rhp3v6Delay = checker.CheckPortOpen(netAddress, rhp3)
// 		params.Rhp4v6, params.Rhp4v6Delay = checker.CheckPortOpen(netAddress, rhp4)
// 		sdk.UpdateCheck(params)
// 		break
// 	}
// 	// fmt.Println(netAddress)
// 	// fmt.Printf("RHP2v4: %v, RHP2v6: %v, RHP3v4: %v, RHP3v6: %v, RHP4v4: %v, RHP4v6: %v\n", rhp2v4, rhp2v6, rhp3v4, rhp3v6, rhp4v4, rhp4v6)
// 	// fmt.Printf("RHP2v4Delay: %v, RHP2v6Delay: %v, RHP3v4Delay: %v, RHP3v6Delay: %v, RHP4v4Delay: %v, RHP4v6Delay: %v\n", rhp2v4Delay, rhp2v6Delay, rhp3v4Delay, rhp3v6Delay, rhp4v4Delay, rhp4v6Delay)
// }
