package main

import (
	"net/http"

	"github.com/ammon134/chirpy/internal/database"
)

type apiConfig struct {
	db         *database.DB
	jwtSecret  string
	serverHits int
}

func (cfg *apiConfig) middlewareHitInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.serverHits++
		next.ServeHTTP(w, r)
	})
}
