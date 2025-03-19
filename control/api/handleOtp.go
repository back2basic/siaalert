package api

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
	"net/http"
	"time"

	"github.com/back2basic/siaalert/control/config"
	"github.com/back2basic/siaalert/control/mail"
	"github.com/back2basic/siaalert/shared/logger"
	"github.com/back2basic/siaalert/shared/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

func getRandomString(length int) (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}
	return string(result), nil
}

func handlePostOtp(w http.ResponseWriter, r *http.Request) {
	cfg := config.GetConfig()
	publicKey := r.URL.Query().Get("publicKey")
	// secret := r.URL.Query().Get("secret")
	// expire := r.URL.Query().Get("expire")
	email := r.URL.Query().Get("email")

	log := logger.GetLogger(cfg.Logging.Path)
	defer log.Sync()

	log.Info("handlePostOtp", zap.String("publicKey", publicKey), zap.String("email", email))
	sec, err := getRandomString(32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	exp := time.Now().Format(time.RFC3339)

	log.Info("Security", zap.String("sec", sec), zap.String("exp", exp))

	err = cfg.DB.UpdateOtp(publicKey, email, exp, sec)
	if err != nil {
		log.Error("UpdateOtp", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = mail.SendOtp(publicKey, email, sec, cfg.Network.Name)
	if err != nil {
		log.Error("SendOtp", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		jsonResp, err := json.Marshal(bson.M{"message": "Could not send email"})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(jsonResp)
		return
	}
	jsonResp, err := json.Marshal(bson.M{"message": "Email sent"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResp)
}

func handlePutOtp(w http.ResponseWriter, r *http.Request) {
	cfg := config.GetConfig()
	log := logger.GetLogger(cfg.Logging.Path)
	defer log.Sync()

	publicKey := r.URL.Query().Get("publicKey")
	secret := r.URL.Query().Get("secret")
	// expire := r.URL.Query().Get("expire")
	email := r.URL.Query().Get("email")

	doc := cfg.DB.FindOtp(bson.M{"publicKey": publicKey, "email": email, "secret": secret})
	if doc.Err() != nil {
		http.Error(w, doc.Err().Error(), http.StatusInternalServerError)
		cfg.DB.DeleteOtp(publicKey, email)
		return
	}

	var otp bson.M
	if err := doc.Decode(&otp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		cfg.DB.DeleteOtp(publicKey, email)
		return
	}

	parsedExpire, err := time.Parse(time.RFC3339, otp["expire"].(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		cfg.DB.DeleteOtp(publicKey, email)
		return
	}

	if time.Since(parsedExpire).Minutes() > 30 {
		http.Error(w, "otp expired", http.StatusUnauthorized)
		cfg.DB.DeleteOtp(publicKey, email)
		return
	}
	cfg.DB.DeleteOtp(publicKey, email)

	host := cfg.DB.FindRhp(bson.M{"publicKey": publicKey})
	if host.Err() != nil {
		http.Error(w, host.Err().Error(), http.StatusInternalServerError)
		return
	}
	var foundHost types.HostScan
	if err := host.Decode(&foundHost); err != nil {
		log.Warn("Failed to decode document", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		// return err
	}

	found, err := cfg.DB.FindAlerts(bson.M{"publicKey": publicKey, "email": email})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(found) == 0 {
		err := cfg.DB.UpdateAlert(publicKey, bson.M{"email": email, "type": "email"})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonResp, err := json.Marshal(bson.M{"message": "enabled", "publicKey": publicKey, "email": email, "address": foundHost.NetAddress})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write(jsonResp)
	} else {
		// fmt.Println("found", len(found))
		for _, f := range found {
			err := cfg.DB.DeleteAlert(f["publicKey"].(string), f["email"].(string))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		jsonResp, err := json.Marshal(bson.M{"message": "disabled", "publicKey": publicKey, "email": email, "address": foundHost.NetAddress})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write(jsonResp)
	}
}
