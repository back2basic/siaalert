package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/back2basic/siaalert/scanner/config"
	"github.com/back2basic/siaalert/scanner/explored"
	"github.com/back2basic/siaalert/scanner/mail"
	"github.com/back2basic/siaalert/scanner/scan"
	"github.com/back2basic/siaalert/shared/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

// func handleGetHost(w http.ResponseWriter, r *http.Request) {
// 	online, err := strconv.ParseBool(r.URL.Query().Get("online"))
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	cfg := config.GetConfig()
// 	hosts, err := cfg.DB.FindHosts(bson.M{"lastScanSuccessful": online})
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	hostsJson, err := json.Marshal(hosts)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	w.Header().Add("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(hostsJson)
// }

type RequestData struct {
	PublicKey string `json:"publicKey"`
}

func handlePostRhpScan(w http.ResponseWriter, r *http.Request) {
	var data RequestData

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	host, err := explored.GetHostByPublicKey(data.PublicKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	scanned, err := scan.RunRhpScan(host, &scan.Checker{})
	if err != nil {
		// fmt.Println("scan error", err)
		scanned.Error = err.Error()
		checkRhpResult(host.NetAddress, false, scanned)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	checkRhpResult(host.NetAddress, true, scanned)

	hostsJson, err := json.Marshal(scanned)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(hostsJson)
}

func checkRhpResult(netAddress string, online bool, result scan.HostScan) {
	cfg := config.GetConfig()

	log := logger.GetLogger(cfg.Logging.Path)
	defer log.Sync()

	prev := cfg.DB.FindRhp(bson.M{"publicKey": result.PublicKey})
	if prev.Err() != nil {
		err := cfg.DB.UpdateRhp(result.PublicKey, online, result.ToBSON(), log)
		if err != nil {
			log.Error("Failed to create rhp", zap.Error(err))
		}
		return
	}

	var prevScan scan.HostScan
	if err := prev.Decode(&prevScan); err != nil {
		log.Warn("Failed to decode rhp", zap.Error(err))
		// return err
	}

	result.OnlineSince = prevScan.OnlineSince
	result.OfflineSince = prevScan.OfflineSince

	if prevScan.Success != online {
		if online {
			log.Info("Host " + result.PublicKey + " is online")
			result.OnlineSince = time.Now()
			result.OfflineSince = time.Time{}
			// Send Mail
			mail.PrepareAlertEmails(netAddress, "Online", result.PublicKey, log)
		} else {
			log.Warn("Host " + result.PublicKey + " is offline")
			result.OnlineSince = time.Time{}
			result.OfflineSince = time.Now()
			// Send Mail
			mail.PrepareAlertEmails(netAddress, "Offline", result.PublicKey, log)
		}
	}
	// result.NetAddress = netAddress
	err := cfg.DB.UpdateRhp(result.PublicKey, online, result.ToBSON(), log)
	if err != nil {
		log.Error("Failed to update rhp", zap.Error(err))
	}
	// log.Info("Finished checking RHP", zap.String("publicKey", result.PublicKey))
}
