package api

import (
	"encoding/json"
	"net/http"

	"github.com/back2basic/siaalert/scanner/config"
	"go.mongodb.org/mongo-driver/bson"
)

func handleGetScan(w http.ResponseWriter, r *http.Request) {
	publicKey := r.URL.Query().Get("publicKey")

	cfg := config.GetConfig()
	scan, err := cfg.DB.FindScan(bson.M{"publicKey": publicKey})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// fmt.Println(scan)
	scanJson, err := json.Marshal(scan)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(scanJson) == 0 {
		http.Error(w, "scan not found", http.StatusNotFound)
		return
	}
	// fmt.Println(scanJson)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(scanJson)
}
