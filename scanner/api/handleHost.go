package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/back2basic/siaalert/scanner/config"
	"github.com/back2basic/siaalert/scanner/explored"
	"github.com/back2basic/siaalert/scanner/logger"
	"github.com/back2basic/siaalert/scanner/scan"

	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

func handleGetHost(w http.ResponseWriter, r *http.Request) {
	online, err := strconv.ParseBool(r.URL.Query().Get("online"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cfg := config.GetConfig()
	hosts, err := cfg.DB.FindHosts(bson.M{"lastScanSuccessful": online})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	hostsJson, err := json.Marshal(hosts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(hostsJson)
}

func handleGetHostScan(w http.ResponseWriter, r *http.Request) {
	publicKey := r.URL.Query().Get("publicKey")

	// cfg := config.GetConfig()
	checker := &scan.Checker{}
	log := logger.GetLogger()
	defer log.Sync()
	host, err := explored.GetHostByPublicKey(publicKey)
	if err != nil {
		log.Error("GetHostByPublicKey", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// log.Info("Scanning host", zap.String("host", host.PublicKey.String()))
	scanned, err := scan.RunRhpScan(host, log, checker)
	if err != nil {
		log.Error("RunRhpScan", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// log.Info("Finished scanning host", zap.String("host", host.PublicKey.String()))
	// hosts, err := cfg.DB.FindHosts(bson.M{"publicKey": publicKey})
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	hostsJson, err := json.Marshal(scanned)
	if err != nil {
		// log.Error("Marshal", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(hostsJson)
}
