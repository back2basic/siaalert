package cron

import (
	"fmt"

	"github.com/back2basic/siadata/siaalert/explored"
	"github.com/back2basic/siadata/siaalert/scan"
	"github.com/back2basic/siadata/siaalert/sdk"
	"github.com/robfig/cron"
)

func cronEveryMinute(c *cron.Cron) {
	c.AddFunc("0 * * * * *", func() {
		state, err := explored.GetConsensus()
		if err != nil {
			fmt.Println(err)
			return
		}
		sdk.UpdateStatus(state)
	})
}

func cronEvery5Minutes(c *cron.Cron) {
	c.AddFunc("0 */5 * * * *", func() {
		if running15Minutes {
			fmt.Println("Skipping scan, cache is being updated")
			return
		}
		checker := scan.Checker{}
		sdk.Mutex.RLock()
		cache := sdk.HostCache
		sdk.Mutex.RUnlock()
		RunScan(cache, checker)
	})
}

var running15Minutes bool

func cronEvery15Minutes(c *cron.Cron) {
	c.AddFunc("10 */15 * * * *", func() {
		if running15Minutes {
			return
		}
		running15Minutes = true
		defer func() {
			running15Minutes = false
		}()
		hosts, err := explored.GetAllHosts()
		sdk.Mutex.RLock()
		cache := sdk.HostCache
		sdk.Mutex.RUnlock()
		if err != nil {
			fmt.Println(err)
			return
		}
		if len(hosts) == 0 {
			fmt.Println("No hosts found")
			return
		}
		fmt.Println("New Hosts available:", len(hosts)-len(cache))
		CheckNewExporedHosts(hosts)
	})
}

var running2Hour bool

func cronEvery2Hour(c *cron.Cron) {
	c.AddFunc("0 0 */2 * * *", func() {
		if running2Hour {
			return
		}
		running2Hour = true
		defer func() {
			running2Hour = false
		}()
		hosts, err := explored.GetAllHosts()
		if err != nil {
			fmt.Println(err)
			return
		}
		RunRhp(hosts)
	})
}

func StartCron() {
	c := cron.New()
	// Every minute
	cronEveryMinute(c)
	// Every 5 minutes
	cronEvery5Minutes(c)
	// Every hour
	cronEvery15Minutes(c)
	// Every 8 hours
	cronEvery2Hour(c)
	// Start cron
	c.Start()
}
