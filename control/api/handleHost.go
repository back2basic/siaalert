package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/back2basic/siaalert/control/config"
	"github.com/back2basic/siaalert/control/explored"
	"github.com/back2basic/siaalert/control/scan"
	"github.com/back2basic/siaalert/shared/logger"
	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson"
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

	cfg := config.GetConfig()
	checker := &scan.Checker{}
	log := logger.GetLogger(cfg.Logging.Path)
	defer log.Sync()
	host, err := explored.GetHostByPublicKey(publicKey)
	if err != nil {
		log.Error("GetHostByPublicKey", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// log.Info("Scanning host", zap.String("host", host.PublicKey.String()))
	scanned, err := scan.RunRhpScan(host, checker)
	if err != nil {
		log.Error("RunRhpScan", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	hostsJson, err := json.Marshal(scanned)
	if err != nil {
		// log.Error("Marshal", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(hostsJson)
	// w.Write([]byte("TODO"))
}
