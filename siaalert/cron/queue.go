package cron

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/back2basic/siadata/siaalert/bench"
	"github.com/back2basic/siadata/siaalert/scan"
	"github.com/back2basic/siadata/siaalert/sdk"
)

func NewWorker(id int, jobQueue chan Job, wg *sync.WaitGroup) Worker {
	return Worker{
		ID:        id,
		JobQueue:  jobQueue,
		Waitgroup: wg,
	}
}

func (w Worker) Start(checker scan.Checker) {
	go func() {
		defer w.Waitgroup.Done()

		for job := range w.JobQueue {
			// do work
			switch job.Type {
			case "scan":
				scanned, err := bench.ScanHosts(job.Address, job.HostKey)
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
					scan := bench.Scan{Settings: bench.Settings{Acceptingcontracts: false, Netaddress: job.Address, Siamuxport: port}, PriceTable: bench.PriceTable{}}
					checker.PortScan(job.Name, scan)
				} else {
					sdk.CheckUpdateStatus(job.Name, job.Address, "", true)
					checker.PortScan(job.Name, scanned)
				}
			}
		}
	}()
}
