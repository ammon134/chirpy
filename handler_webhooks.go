package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ammon134/chirpy/internal/auth"
	"github.com/ammon134/chirpy/internal/database"
)

type EventType string

const (
	EventUserUpgraded EventType = "user.upgraded"
)

func (cfg *apiConfig) handlerWebhookUpgradeUser(w http.ResponseWriter, r *http.Request) {
	apikey, err := auth.GetBearerToken(r.Header, auth.AuthTypeAPIKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	if apikey != cfg.polka_apikey {
		respondWithError(w, http.StatusUnauthorized, "invalid apikey")
		return
	}

	decoder := json.NewDecoder(r.Body)
	type parameters struct {
		Event EventType `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		} `json:"data"`
	}

	params := &parameters{}
	err = decoder.Decode(params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not decode request body")
		return
	}

	if params.Event != EventUserUpgraded {
		respondWithJSON(w, http.StatusOK, http.StatusText(http.StatusOK))
		return
	}

	err = cfg.db.UpgradeUser(params.Data.UserID)
	if err != nil {
		if errors.Is(err, database.ErrNotExist) {
			respondWithError(w, http.StatusNotFound, "user does not exist")
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, http.StatusText(http.StatusOK))
}
