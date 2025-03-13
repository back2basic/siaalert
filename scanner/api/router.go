package api

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"*"},
		AllowedOrigins: []string{"https://siaalert.euregiohosting.nl"}, // Use this to allow specific origin hosts
		// AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	// r.Use(middleware.Throttle(10))
	// r.Use(httprate.LimitByRealIP(100, 1*time.Minute))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("SiaAlert API"))
	})

	r.Route("/auth/otp", func(r chi.Router) {
		r.Post("/", handlePostOtp)
		r.Put("/", handlePutOtp)
	})
	r.Mount("/v1", v1Router())

	return r
}

func v1Router() *chi.Mux {
	r := chi.NewRouter()

	r.Route("/consensus", func(r chi.Router) { r.Get("/", handleGetConsensus) })
	r.Route("/host", func(r chi.Router) { r.Get("/", handleGetHost) })
	r.Route("/host/scan", func(r chi.Router) { r.Get("/", handleGetHostScan) })
	r.Route("/rhp", func(r chi.Router) { r.Get("/", handleGetRhp) })
	r.Route("/scan", func(r chi.Router) { r.Get("/", handleGetScan) })
	return r
}
