package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/back2basic/siaalert/scanner/config"
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
