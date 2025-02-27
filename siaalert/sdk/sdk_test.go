package sdk_test

import (
	"flag"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/back2basic/siadata/siaalert/config"
	"github.com/back2basic/siadata/siaalert/explored"
	"github.com/back2basic/siadata/siaalert/sdk"
	"github.com/stretchr/testify/assert"
)

func TestPrepareAppwrite(t *testing.T) {
	// Setup a mock environment for configuration and database service
	cfg := config.LoadConfig("../config.yaml")

	// Startup
	log.SetFlags(0)
	log.SetOutput(os.Stdout)
	flag.Parse()

	fmt.Println("Init Appwrite:")
	dbSvc := sdk.PrepareAppwrite(cfg)
	assert.NotNil(t, dbSvc)
}

func TestGetHostByPublicKey(t *testing.T) {
	// Setup a mock database service and the necessary environment
	cfg := config.LoadConfig("../config.yaml")
	fmt.Println("Init Appwrite:")
	dbSvc := sdk.PrepareAppwrite(cfg)
	sdk.GetAppwriteDatabaseService().Client = dbSvc.(*sdk.AppwriteDatabaseService).Client
	cfg.Appwrite.Database = sdk.PrepareDatabase(cfg, sdk.GetAppwriteDatabaseService())
	cfg.Appwrite.ColHosts, cfg.Appwrite.ColStatus, cfg.Appwrite.ColAlert, cfg.Appwrite.ColCheck, cfg.Appwrite.ColRhp2 = sdk.PrepareCollection(sdk.GetAppwriteDatabaseService(), cfg.Appwrite.Database.Id)

	databaseID := cfg.Appwrite.Database.Id
	collectionID := cfg.Appwrite.ColHosts.Id
	publicKey := "1"

	host, err := sdk.GetHostByPublicKey(databaseID, collectionID, publicKey)
	fmt.Println(host.Documents[0])
	assert.NoError(t, err)
	assert.NotNil(t, host)
}

func TestCheckHost(t *testing.T) {
	// Setup a mock environment for configuration and database service
	cfg := config.LoadConfig("../config.yaml")
	fmt.Println("Init Appwrite:")
	dbSvc := sdk.PrepareAppwrite(cfg)
	sdk.GetAppwriteDatabaseService().Client = dbSvc.(*sdk.AppwriteDatabaseService).Client
	cfg.Appwrite.Database = sdk.PrepareDatabase(cfg, sdk.GetAppwriteDatabaseService())
	cfg.Appwrite.ColHosts, cfg.Appwrite.ColStatus, cfg.Appwrite.ColAlert, cfg.Appwrite.ColCheck, cfg.Appwrite.ColRhp2 = sdk.PrepareCollection(sdk.GetAppwriteDatabaseService(), cfg.Appwrite.Database.Id)

	host := explored.Host{
		PublicKey:  "1",
		NetAddress: "[::1]:9982",
		// Other fields...
	}

	doc, err := sdk.CheckHost(host)
	assert.NoError(t, err)
	assert.NotNil(t, doc)
}

func TestCreateHost(t *testing.T) {
	// Setup a mock environment for configuration and database service
	cfg := config.LoadConfig("../config.yaml")
	fmt.Println("Init Appwrite:")
	dbSvc := sdk.PrepareAppwrite(cfg)
	sdk.GetAppwriteDatabaseService().Client = dbSvc.(*sdk.AppwriteDatabaseService).Client
	cfg.Appwrite.Database = sdk.PrepareDatabase(cfg, sdk.GetAppwriteDatabaseService())
	cfg.Appwrite.ColHosts, cfg.Appwrite.ColStatus, cfg.Appwrite.ColAlert, cfg.Appwrite.ColCheck, cfg.Appwrite.ColRhp2 = sdk.PrepareCollection(sdk.GetAppwriteDatabaseService(), cfg.Appwrite.Database.Id)


	knownSince, err := time.Parse(time.RFC3339, "2023-04-02T13:41:37Z")
	if err != nil {
		t.Fatal(err)
	}

	lastScan, err := time.Parse(time.RFC3339, "2025-02-27T10:14:39Z")
	if err != nil {
		t.Fatal(err)
	}

	lastAnnouncement, err := time.Parse(time.RFC3339, "2025-02-24T22:10:21Z")
	if err != nil {
		t.Fatal(err)
	}
	host := explored.Host{
		PublicKey:              "ed25519:1",
		V2:                     false,
		NetAddress:             "storage:9882",
		CountryCode:            "FR",
		KnownSince:             knownSince,
		LastScan:               lastScan,
		LastScanSuccessful:     true,
		LastAnnouncement:       lastAnnouncement,
		TotalScans:             68,
		SuccessfulInteractions: 68,
		FailedInteractions:     0,
	}

	doc, err := sdk.CreateHost(host)
	assert.NoError(t, err)
	assert.NotNil(t, doc)
}

// !! WIP !!

// func TestUpdateHost(t *testing.T) {
// 	// Setup a mock environment for configuration and database service

// 	host := explored.Host{
// 		PublicKey:  "testPublicKey",
// 		NetAddress: "testNetAddress",
// 		// Other fields...
// 	}

// 	hostDoc := sdk.HostDocument{
// 		NetAddress: "testNetAddress",
// 		PublicKey:  "testPublicKey",
// 		// Other fields...
// 	}

// 	err := sdk.UpdateHost(host, "testHostId", hostDoc)
// 	assert.NoError(t, err)
// }

// func TestUpdateNetAddress(t *testing.T) {
// 	// Setup a mock environment for configuration and database service

// 	hostDoc := HostDocument{
// 		PublicKey:  "testPublicKey",
// 		NetAddress: "testNetAddress",
// 		// Other fields...
// 	}

// 	err := sdk.UpdateNetAddress(hostDoc)
// 	assert.NoError(t, err)
// }

// func TestCheckUpdateStatus(t *testing.T) {
// 	// Setup a mock environment for configuration and database service

// 	err := sdk.CheckUpdateStatus("testHostId", "testNetAddress", "testError", true)
// 	assert.NoError(t, err)
// }

// func TestUpdateHostStatus(t *testing.T) {
// 	// Setup a mock environment for configuration and database service

// 	hostDoc := sdk.HostDocument{
// 		Host:      "testHost",
// 		PublicKey: "testPublicKey",
// 		// Other fields...
// 	}

// 	err := UpdateHostStatus("testError", time.Now().Format(time.RFC3339), "", true, hostDoc)
// 	assert.NoError(t, err)
// }
