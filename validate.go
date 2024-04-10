package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
)

func handlerValidate(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	type resBody struct {
		CleanedBody string `json:"cleaned_body"`
	}
	type parameters struct {
		Body string `json:"body"`
	}
	params := &parameters{}
	err := decoder.Decode(params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode request body")
		return
	}
	if len(strings.TrimSpace(params.Body)) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	res := resBody{
		CleanedBody: cleanChirp(params.Body),
	}
	respondWithJSON(w, http.StatusOK, res)
}

func cleanChirp(msg string) string {
	words := strings.Fields(msg)
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
