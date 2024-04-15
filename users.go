package main

import (
	"encoding/json"
	"net/http"

	"github.com/ammon134/chirpy/internal/database"
	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	// decoder to read param
	decoder := json.NewDecoder(r.Body)
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	params := &parameters{}
	err := decoder.Decode(params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode request body")
		return
	}

	// use param to create User
	hash, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// then update createUser to save hash password
	// then return user without password hash
	user, err := cfg.db.CreateUser(params.Email, hash)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// return User
	respondWithJSON(w, http.StatusCreated, database.User{
		ID:    user.ID,
		Email: user.Email,
	})
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}
	params := &parameters{}
	err := decoder.Decode(params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// compare password and user
	user, err := cfg.db.GetUser(params.Email)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	err = bcrypt.CompareHashAndPassword(user.Hash, []byte(params.Password))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	ss, err := createJWTToken(cfg, user, params.ExpiresInSeconds)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	type response struct {
		Email string
		Token string
		ID    int
	}
	respondWithJSON(w, http.StatusOK, response{
		ID:    user.ID,
		Email: user.Email,
		Token: ss,
	})
}
