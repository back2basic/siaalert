package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/back2basic/siadata/siaalert/config"
	"github.com/back2basic/siadata/siaalert/cron"
	"github.com/back2basic/siadata/siaalert/explored"
	"github.com/back2basic/siadata/siaalert/sdk"
	"github.com/back2basic/siadata/siaalert/strict"
)

func SaveHostCacheToFile(filename string, cache map[string]strict.HostDocument) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(cache)
	if err != nil {
		return err
	}

	return nil
}

func LoadHostCacheFromFile(filename string) (map[string]strict.HostDocument, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	cache := make(map[string]strict.HostDocument)
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cache)
	if err != nil {
		return nil, err
	}

	return cache, nil
}

func init() {
	// Load the cache from the file
	cache, err := LoadHostCacheFromFile("cache/host_cache.json")
	if err != nil {
		fmt.Println("Error loading cache:", err)
	} else {
		sdk.Mutex.Lock()
		sdk.HostCache = cache
		sdk.Mutex.Unlock()
		fmt.Println("Cache loaded:", len(cache), "hosts")
	}
}

func main() {
	// Handle signals for graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		fmt.Println("Saving cache...")
		err := SaveHostCacheToFile("cache/host_cache.json", sdk.HostCache)
		if err != nil {
			fmt.Println("Error saving cache:", err)
		}
		os.Exit(0)
	}()

	cfg := config.LoadConfig("./data/config.yaml")

	// Startup
	log.SetFlags(0)
	log.SetOutput(os.Stdout)
	flag.Parse()

	fmt.Println("Init Appwrite:")
	dbSvc := sdk.PrepareAppwrite(cfg)
	sdk.GetAppwriteDatabaseService().Client = dbSvc.(*sdk.AppwriteDatabaseService).Client

	cfg.Appwrite.Database = sdk.PrepareDatabase(cfg, sdk.GetAppwriteDatabaseService())
	cfg.Appwrite.ColHosts, cfg.Appwrite.ColStatus, cfg.Appwrite.ColAlert, cfg.Appwrite.ColCheck, cfg.Appwrite.ColRhp2, cfg.Appwrite.ColRhp3 = sdk.PrepareCollection(sdk.GetAppwriteDatabaseService(), cfg.Appwrite.Database.Id)

	// cache explored hosts
	explored.GetAllHosts()

	// start cron
	cron.StartCron()

	fmt.Println("Startup Complete")

	defer func() {
		fmt.Println("Saving cache...")
		err := SaveHostCacheToFile("cache/host_cache.json", sdk.HostCache)
		if err != nil {
			fmt.Println("Error saving cache:", err)
		}
	}()
	select {}
}
