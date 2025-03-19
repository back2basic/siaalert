package api

import (
	"net/http"

	"go.uber.org/zap"
)

func StartServer(log *zap.Logger) {
	log.Info("API is starting...", zap.String("module", "api"), zap.Int("port", 8080))
	http.ListenAndServe(":8080", NewRouter())
}
