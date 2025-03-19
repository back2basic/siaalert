package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/back2basic/siaalert/control/config"
	"github.com/back2basic/siaalert/shared/logger"

	"go.mongodb.org/mongo-driver/bson"
)

func handleGetRhp(w http.ResponseWriter, r *http.Request) {
	netaddress := r.URL.Query().Get("search")
	online, err := strconv.ParseBool(r.URL.Query().Get("online"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// fmt.Println(online, netaddress)
	cfg := config.GetConfig()
	log := logger.GetLogger(cfg.Logging.Path)
	defer log.Sync()
	if netaddress != "" {
		var filter bson.M
		if online {
			filter = bson.M{
				"netAddress": bson.M{"$regex": netaddress, "$options": "i"}, // Case-insensitive regex
				"success":    online,
			}
		} else {
			filter = bson.M{
				"netAddress": bson.M{"$regex": netaddress, "$options": "i"}, // Case-insensitive regex
			}
		}

		rhp, err := cfg.DB.FindRhps(filter)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rhpsJson, err := json.Marshal(rhp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(rhpsJson)
		return
	}
	rhp, err := cfg.DB.FindRhps(bson.M{"success": online})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// fmt.Println(len(rhp))
	rhpsJson, err := json.Marshal(rhp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(rhpsJson)
}
