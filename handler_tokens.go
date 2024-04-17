package main

import (
	"net/http"

	"github.com/ammon134/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevokeToken(w http.ResponseWriter, r *http.Request) {
	bearerToken, err := auth.GetBearerToken(r.Header, auth.AuthTypeBearer)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not find JWT token")
		return
	}

	err = cfg.db.RevokeToken(bearerToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, http.StatusText(http.StatusOK))
}

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	bearerToken, err := auth.GetBearerToken(r.Header, auth.AuthTypeBearer)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not find JWT token")
		return
	}

	revoked, err := cfg.db.IsRevoked(bearerToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if revoked {
		respondWithError(w, http.StatusUnauthorized, "JWT revoked")
		return
	}

	accessToken, err := auth.RefreshJWT(cfg.jwtSecret, bearerToken)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "could not create JWT")
		return
	}

	type response struct {
		Token string `json:"token"`
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})
}
