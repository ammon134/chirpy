package main

import "net/http"

type apiConfig struct {
	serverHits int
}

func (cfg *apiConfig) middlewareHitInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.serverHits++
		next.ServeHTTP(w, r)
	})
}
