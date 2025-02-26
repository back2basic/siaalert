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
		checker := scan.Checker{}
		RunScan(sdk.HostCache, checker)
	})
}

func cronEvery15Minutes(c *cron.Cron) {
	c.AddFunc("0 */15 * * * *", func() {
		hosts, err := explored.GetAllHosts()
		if err != nil {
			fmt.Println(err)
			return
		}
		if len(hosts) == 0 {
			fmt.Println("No hosts found")
			return
		}
		fmt.Println("New Hosts available:", len(hosts)-len(sdk.HostCache))
		CheckNewExporedHosts(hosts)
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
	// Start cron
	c.Start()
}
