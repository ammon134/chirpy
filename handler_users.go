package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ammon134/chirpy/internal/auth"
	"github.com/ammon134/chirpy/internal/database"
)

type User struct {
	Email    string `json:"email"`
	Password string `json:"-"`
	ID       int    `json:"id"`
}

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
		respondWithError(w, http.StatusInternalServerError, "could not decode request body")
		return
	}

	// use param to create User
	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not hash password")
		return
	}
	// then update createUser to save hash password
	// then return user without password hash
	user, err := cfg.db.CreateUser(params.Email, hash)
	if err != nil {
		if errors.Is(err, database.ErrAlreadyExist) {
			respondWithError(w, http.StatusConflict, "user already exist")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "could not create user in database")
		return
	}
	// return User
	type response struct {
		User
	}
	respondWithJSON(w, http.StatusCreated, response{
		User: User{
			ID:    user.ID,
			Email: user.Email,
		},
	})
}

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.ParseForUserID(cfg.jwtSecret, r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	decoder := json.NewDecoder(r.Body)
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	params := &parameters{}
	err = decoder.Decode(params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not decode request params")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not hash password")
		return
	}

	user, err := cfg.db.UpdateUser(userID, params.Email, hashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not update user")
		return
	}

	type response struct {
		User
	}
	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:    user.ID,
			Email: user.Email,
		},
	})
}
