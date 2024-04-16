package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ammon134/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	params := &parameters{}
	err := decoder.Decode(params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not decode request params")
		return
	}

	// compare password and user
	user, err := cfg.db.GetUserByEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "could not get user")
		return
	}
	err = auth.CheckPasswordHash(user.HashedPassword, params.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid password")
		return
	}
	accessToken, err := auth.CreateJWT(cfg.jwtSecret, user.ID, time.Hour, auth.TokenTypeAccess)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "could not create JWT")
		return
	}
	refreshToken, err := auth.CreateJWT(cfg.jwtSecret, user.ID, time.Hour*1440, auth.TokenTypeRefresh)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "could not create JWT")
		return
	}

	// type response struct
	type response struct {
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
		User
	}
	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:    user.ID,
			Email: user.Email,
		},
		Token:        accessToken,
		RefreshToken: refreshToken,
	})
}
