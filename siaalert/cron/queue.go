package cron

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/back2basic/siadata/siaalert/scan"
	"github.com/back2basic/siadata/siaalert/sdk"
	"go.sia.tech/core/rhp/v2"
	"go.sia.tech/core/types"
)

func NewWorker(id int, jobQueue chan Job, wg *sync.WaitGroup) Worker {
	return Worker{
		ID:        id,
		JobQueue:  jobQueue,
		Waitgroup: wg,
	}
}

func stringToPublicKey(key string) (types.PublicKey, error) {
	var publicKey types.PublicKey
	err := publicKey.UnmarshalText([]byte(key))
	return publicKey, err
}

func (w Worker) Start(checker scan.Checker) {
	go func() {
		defer w.Waitgroup.Done()

		for job := range w.JobQueue {
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
					scanned, err = checker.ScanV2Host(scan.UnscannedHost{
						NetAddress: job.Address,
						PublicKey:  publicKey,
					})
					if err != nil {
						sdk.CheckUpdateStatus(job.Name, job.Address, err.Error(), false)
						_, port, err := checker.SplitAddressPort(job.Address)
						if err != nil {
							fmt.Println(err)
						}
						// convert string to int add 1 to port and change back to string
						mux, err := strconv.Atoi(port)
						if err != nil {
							fmt.Println(err)
						}
						mux++
						port = strconv.Itoa(mux)
						scan := scan.HostScan{Settings: rhp.HostSettings{AcceptingContracts: false, NetAddress: job.Address, SiaMuxPort: port}}
						checker.PortScan(job.Name, scan)
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
						}
						// convert string to int add 1 to port and change back to string
						mux, err := strconv.Atoi(port)
						if err != nil {
							fmt.Println(err)
						}
						mux++
						port = strconv.Itoa(mux)
						scan := scan.HostScan{Settings: rhp.HostSettings{AcceptingContracts: false, NetAddress: job.Address, SiaMuxPort: port}}
						// scan := scan.HostScan{Settings: bench.Settings{Acceptingcontracts: false, Netaddress: job.Address, Siamuxport: port}, PriceTable: bench.PriceTable{}}
						checker.PortScan(job.Name, scan)
						continue
					}
				}
				sdk.CheckUpdateStatus(job.Name, job.Address, "", true)
				checker.PortScan(job.Name, scanned)
				// scanned, err := bench.ScanHosts(job.Address, job.HostKey)
				// if err != nil {
				// 	sdk.CheckUpdateStatus(job.Name, job.Address, err.Error(), false)
				// 	_, port, err := checker.SplitAddressPort(job.Address)
				// 	if err != nil {
				// 		fmt.Println(err)
				// 	}
				// 	// convert string to int add 1 to port and change back to string
				// 	mux, err := strconv.Atoi(port)
				// 	if err != nil {
				// 		fmt.Println(err)
				// 	}
				// 	mux++
				// 	port = strconv.Itoa(mux)
				// 	scan := bench.Scan{Settings: bench.Settings{Acceptingcontracts: false, Netaddress: job.Address, Siamuxport: port}, PriceTable: bench.PriceTable{}}
				// 	checker.PortScan(job.Name, scan)
				// } else {
				// 	sdk.CheckUpdateStatus(job.Name, job.Address, "", true)
				// 	checker.PortScan(job.Name, scanned)
				// }
			}
		}
	}()
}
