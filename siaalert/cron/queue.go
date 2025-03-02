package cron

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/back2basic/siadata/siaalert/scan"
	"github.com/back2basic/siadata/siaalert/sdk"
	"github.com/back2basic/siadata/siaalert/strict"

	"go.sia.tech/core/rhp/v2"
	"go.sia.tech/core/types"
)

// func NewWorker(id int, jobQueue chan Job, wg *sync.WaitGroup) Worker {
// 	return Worker{
// 		ID:        id,
// 		JobQueue:  jobQueue,
// 		Waitgroup: wg,
// 	}
// }

func stringToPublicKey(key string) (types.PublicKey, error) {
	var publicKey types.PublicKey
	err := publicKey.UnmarshalText([]byte(key))
	return publicKey, err
}

func worker(id int, tasks <-chan Job, wg *sync.WaitGroup, checker scan.Checker) {
	defer func() {
		Running5Minutes = false
		// print wg state
		fmt.Println("Scan complete")
		wg.Done()
	}()

	// Create channel for upadting sdk queue
	const numWorkers = 2
	sdkQueue := make(chan strict.TaskCheckDoc)
	var sdkWg sync.WaitGroup
	// Start worker goroutines
	for i := 1; i <= numWorkers; i++ {
		go sdk.SdkWorker(i, sdkQueue, &sdkWg)
	}

	for job := range tasks {
		// fmt.Printf("WIP %d: %s\n", w.ID, job.Type)
		publicKey, err := stringToPublicKey(job.HostKey)
		if err != nil {
			// handle the error
			fmt.Println("failed getting public key", err)
			continue
		}
		// do work
		switch job.Type {
		case "scan":
			var scanned scan.HostScan
			if job.V2 {
				fmt.Println("!!! V2 host Found", job.Address, job.Name)
				scanned, err = checker.ScanV2Host(scan.UnscannedHost{
					NetAddress: job.Address,
					PublicKey:  publicKey,
				})
				if err != nil {
					sdk.CheckUpdateStatus(job.Name, job.Address, err.Error(), false)
					_, port, err := checker.SplitAddressPort(job.Address)
					if err != nil {
						fmt.Println(err)
						continue
					}
					// convert string to int add 1 to port and change back to string
					mux, err := strconv.Atoi(port)
					if err != nil {
						fmt.Println(err)
						continue
					}
					mux++
					port = strconv.Itoa(mux)
					scan := scan.HostScan{Settings: rhp.HostSettings{AcceptingContracts: false, NetAddress: job.Address, SiaMuxPort: port}}
					checker.PortScan(job.Name, scan, &sdkWg, sdkQueue)
					continue
				}
			} else {
				scanned, err = checker.ScanV1Host(scan.UnscannedHost{
					NetAddress: job.Address,
					PublicKey:  publicKey,
				})
				if err != nil {
					sdk.CheckUpdateStatus(job.Name, job.Address, err.Error(), false)
					_, port, err := checker.SplitAddressPort(job.Address)
					if err != nil {
						fmt.Println(err)
						continue
					}
					// convert string to int add 1 to port and change back to string
					mux, err := strconv.Atoi(port)
					if err != nil {
						fmt.Println(err)
						continue
					}
					mux++
					port = strconv.Itoa(mux)
					scan := scan.HostScan{Settings: rhp.HostSettings{AcceptingContracts: false, NetAddress: job.Address, SiaMuxPort: port}}
					// scan := scan.HostScan{Settings: bench.Settings{Acceptingcontracts: false, Netaddress: job.Address, Siamuxport: port}, PriceTable: bench.PriceTable{}}
					checker.PortScan(job.Name, scan, &sdkWg, sdkQueue)
					continue
				}
			}
			sdk.CheckUpdateStatus(job.Name, job.Address, "", true)
			checker.PortScan(job.Name, scanned, &sdkWg, sdkQueue)
		default:
			fmt.Printf("WIP %d: Unknown job: %s\n", id, job.Address)
		}
	}

	close(sdkQueue)
	Running5Minutes = false
	sdkWg.Wait()
}
