package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/ammon134/chirpy/internal/auth"
	"github.com/ammon134/chirpy/internal/database"
)

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	type parameters struct {
		Body string `json:"body"`
	}
	params := &parameters{}
	err := decoder.Decode(params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode request body")
		return
	}

	userID, err := auth.ParseForUserID(cfg.jwtSecret, r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Validate chirp length
	if len(strings.TrimSpace(params.Body)) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	chirp, err := cfg.db.CreateChirp(cleanChirp(params.Body), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, chirp)
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	authorIDStr := r.URL.Query().Get("author_id")
	authorID, err := strconv.Atoi(authorIDStr)
	if err != nil {
		authorID = -1
	}

	chirps, err := cfg.db.GetChirpsByAuthor(authorID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sortParam := r.URL.Query().Get("sort")
	if sortParam == "desc" {
		sort.Slice(chirps, func(i, j int) bool { return chirps[i].ID > chirps[j].ID })
	} else {
		sort.Slice(chirps, func(i, j int) bool { return chirps[i].ID < chirps[j].ID })
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	chirp, err := cfg.db.GetChirp(idInt)
	if err != nil {
		if errors.Is(err, database.ErrNotExist) {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	respondWithJSON(w, http.StatusOK, chirp)
}

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.ParseForUserID(cfg.jwtSecret, r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
	}

	chirpIDStr := r.PathValue("id")
	chirpID, err := strconv.Atoi(chirpIDStr)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	chirp, err := cfg.db.GetChirp(chirpID)
	if err != nil {
		if !errors.Is(err, database.ErrNotExist) {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	if chirp.AuthorID != userID {
		respondWithError(w, http.StatusForbidden, "user does not have permission")
		return
	}

	err = cfg.db.DeleteChirp(chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, http.StatusText(http.StatusOK))
}

func cleanChirp(msg string) string {
	words := strings.Fields(msg)
	// PERF: refactor badWords from list to map
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	cleanedWords := []string{}
	for _, word := range words {
		if slices.Contains(badWords, strings.ToLower(word)) {
			cleanedWords = append(cleanedWords, "****")
		} else {
			cleanedWords = append(cleanedWords, word)
		}
	}
	return strings.Join(cleanedWords, " ")
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type errBody struct {
		Error string
	}
	respondWithJSON(w, code, errBody{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marchalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}
