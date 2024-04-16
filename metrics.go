package main

import (
	"net/http"
)

func (cfg *apiConfig) middlewareHitInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.serverHits++
		next.ServeHTTP(w, r)
	})
}
