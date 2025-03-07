package api

import (
	"encoding/json"
	"net/http"

	"github.com/back2basic/siaalert/scanner/explored"
)

func handleGetConsensus(w http.ResponseWriter, r *http.Request) {
	consensus, err := explored.GetConsensus()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	consensusJson, err := json.Marshal(consensus)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(consensusJson)
}
