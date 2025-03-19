package api

import (
	"encoding/json"
	"net/http"

	"github.com/back2basic/siaalert/scanner/config"
	"github.com/back2basic/siaalert/scanner/scan"
	"go.mongodb.org/mongo-driver/bson"
)


func handlePostPortScan(w http.ResponseWriter, r *http.Request) {
	var data RequestData

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cfg := config.GetConfig()
	checker := &scan.Checker{}
	rhp := cfg.DB.FindRhp(bson.M{"publicKey": data.PublicKey})
	if rhp.Err() != nil {
		http.Error(w, rhp.Err().Error(), http.StatusInternalServerError)
		return
	}
	var rhpScan scan.HostScan
	if err := rhp.Decode(&rhpScan); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	checker.PortScan(rhpScan)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
