package sdk_test

import (
	"flag"
	"fmt"
	"log"
	"os"
	"testing"

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

	// Startup
	log.SetFlags(0)
	log.SetOutput(os.Stdout)
	flag.Parse()

	fmt.Println("Init Appwrite:")
	dbSvc := sdk.PrepareAppwrite(cfg)
	sdk.GetAppwriteDatabaseService().Client = dbSvc.(*sdk.AppwriteDatabaseService).Client


	databaseID := "67bf83fe0014e0b1b4ef"
	collectionID := "67bf8471002b2ae85dbd"
	publicKey := "1"

	host, err := sdk.GetHostByPublicKey(databaseID, collectionID, publicKey)
	fmt.Println(host.Documents[0])
	assert.NoError(t, err)
	assert.NotNil(t, host)
}

func TestCheckHost(t *testing.T) {
	// Setup a mock environment for configuration and database service

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

	host := explored.Host{
		PublicKey:  "testPublicKey",
		NetAddress: "testNetAddress",
		// Other fields...
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
