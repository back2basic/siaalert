package cron_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/back2basic/siadata/siaalert/config"
	"github.com/back2basic/siadata/siaalert/sdk"
	"github.com/robfig/cron"
	"github.com/stretchr/testify/assert"
)

func TestCronEveryMinute(t *testing.T) {
	cfg := config.LoadConfig("../config.yaml")
	fmt.Println("Init Appwrite:")
	dbSvc := sdk.PrepareAppwrite(cfg)
	sdk.GetAppwriteDatabaseService().Client = dbSvc.(*sdk.AppwriteDatabaseService).Client
	cfg.Appwrite.Database = sdk.PrepareDatabase(cfg, sdk.GetAppwriteDatabaseService())
	cfg.Appwrite.ColHosts, cfg.Appwrite.ColStatus, cfg.Appwrite.ColAlert, cfg.Appwrite.ColCheck, cfg.Appwrite.ColRhp2, cfg.Appwrite.ColRhp3 = sdk.PrepareCollection(sdk.GetAppwriteDatabaseService(), cfg.Appwrite.Database.Id)

	c := cron.New()
	// cron.ronEveryMinute(c)

	// Check if the cron job was added successfully
	assert.Len(t, c.Entries(), 1)

	// Check if the cron job is scheduled to run every minute
	entry := c.Entries()[0]
	schedule := entry.Schedule
	assert.NotNil(t, schedule)

	// Verify the scheduling rules
	next := schedule.Next(time.Now())
	assert.NotNil(t, next)
	assert.WithinDuration(t, time.Now().Add(1*time.Minute), next, 1*time.Minute)
}

func TestCronEvery5Minutes(t *testing.T) {
	cfg := config.LoadConfig("../config.yaml")
	fmt.Println("Init Appwrite:")
	dbSvc := sdk.PrepareAppwrite(cfg)
	sdk.GetAppwriteDatabaseService().Client = dbSvc.(*sdk.AppwriteDatabaseService).Client
	cfg.Appwrite.Database = sdk.PrepareDatabase(cfg, sdk.GetAppwriteDatabaseService())
	cfg.Appwrite.ColHosts, cfg.Appwrite.ColStatus, cfg.Appwrite.ColAlert, cfg.Appwrite.ColCheck, cfg.Appwrite.ColRhp2, cfg.Appwrite.ColRhp3 = sdk.PrepareCollection(sdk.GetAppwriteDatabaseService(), cfg.Appwrite.Database.Id)

	c := cron.New()
	// cronEvery5Minutes(c)

	// Check if the cron job was added successfully
	assert.Len(t, c.Entries(), 1)

	// Check if the cron job is scheduled to run every 5 minutes
	entry := c.Entries()[0]
	schedule := entry.Schedule
	next := schedule.Next(time.Now())
	assert.NotNil(t, next)
	assert.WithinDuration(t, time.Now().Add(5*time.Minute), next, 5*time.Minute)
}

func TestCronEvery15Minutes(t *testing.T) {
	cfg := config.LoadConfig("../config.yaml")
	fmt.Println("Init Appwrite:")
	dbSvc := sdk.PrepareAppwrite(cfg)
	sdk.GetAppwriteDatabaseService().Client = dbSvc.(*sdk.AppwriteDatabaseService).Client
	cfg.Appwrite.Database = sdk.PrepareDatabase(cfg, sdk.GetAppwriteDatabaseService())
	cfg.Appwrite.ColHosts, cfg.Appwrite.ColStatus, cfg.Appwrite.ColAlert, cfg.Appwrite.ColCheck, cfg.Appwrite.ColRhp2, cfg.Appwrite.ColRhp3 = sdk.PrepareCollection(sdk.GetAppwriteDatabaseService(), cfg.Appwrite.Database.Id)

	c := cron.New()
	// cronEvery15Minutes(c)

	// Check if the cron job was added successfully
	assert.Len(t, c.Entries(), 1)

	// Check if the cron job is scheduled to run every 15 minutes
	entry := c.Entries()[0]
	schedule := entry.Schedule
	next := schedule.Next(time.Now())
	assert.NotNil(t, next)
	assert.WithinDuration(t, time.Now().Add(15*time.Minute), next, 15*time.Minute)
}

func TestCronEvery8Hour(t *testing.T) {
	cfg := config.LoadConfig("../config.yaml")
	fmt.Println("Init Appwrite:")
	dbSvc := sdk.PrepareAppwrite(cfg)
	sdk.GetAppwriteDatabaseService().Client = dbSvc.(*sdk.AppwriteDatabaseService).Client
	cfg.Appwrite.Database = sdk.PrepareDatabase(cfg, sdk.GetAppwriteDatabaseService())
	cfg.Appwrite.ColHosts, cfg.Appwrite.ColStatus, cfg.Appwrite.ColAlert, cfg.Appwrite.ColCheck, cfg.Appwrite.ColRhp2, cfg.Appwrite.ColRhp3 = sdk.PrepareCollection(sdk.GetAppwriteDatabaseService(), cfg.Appwrite.Database.Id)

	c := cron.New()
	// cronEvery2Hour(c)

	// Check if the cron job was added successfully
	assert.Len(t, c.Entries(), 1)

	// Check if the cron job is scheduled to run every 8 hours
	entry := c.Entries()[0]
	schedule := entry.Schedule
	next := schedule.Next(time.Now())
	assert.NotNil(t, next)
	assert.WithinDuration(t, time.Now().Add(8*time.Hour), next, 8*time.Hour)
}

func TestStartCron(t *testing.T) {
	// Create a new cron instance
	c := cron.New()

	// Start the cron
	// StartCron()

	// Check if the cron jobs were added successfully
	assert.Len(t, c.Entries(), 4)

	// Check if the cron jobs are scheduled to run at the correct intervals
	entry1 := c.Entries()[0]
	entry2 := c.Entries()[1]
	entry3 := c.Entries()[2]
	entry4 := c.Entries()[3]

	// Verify the scheduling rules
	assert.NotNil(t, entry1.Schedule)
	assert.NotNil(t, entry2.Schedule)
	assert.NotNil(t, entry3.Schedule)
	assert.NotNil(t, entry4.Schedule)

	// You can also verify the next run time
	next1 := entry1.Schedule.Next(time.Now())
	next2 := entry2.Schedule.Next(time.Now())
	next3 := entry3.Schedule.Next(time.Now())
	next4 := entry4.Schedule.Next(time.Now())

	assert.NotNil(t, next1)
	assert.NotNil(t, next2)
	assert.NotNil(t, next3)
	assert.NotNil(t, next4)
}

func TestRunScan(t *testing.T) {
	// Create a mock host cache
	// hostCache := make(map[string]sdk.HostDocument)

	// Create a mock checker
	// checker := scan.Checker{}

	// Run the scan
	// RunScan(hostCache, &checker)

	// Check if the scan was run successfully
	// NOTE: This test is incomplete, as it doesn't check the actual behavior of the RunScan function.
	// You may need to add additional assertions or use a mocking library to test the behavior of the RunScan function.
}

func TestRunRhp(t *testing.T) {
	// Create a mock list of hosts
	// hosts := map[string]explored.Host{}

	// Run the RHP
	// RunRhp(hosts)

	// Check if the RHP was run successfully
	// NOTE: This test is incomplete, as it doesn't check the actual behavior of the RunRhp function.
	// You may need to add additional assertions or use a mocking library to test the behavior of the RunRhp function.
}
