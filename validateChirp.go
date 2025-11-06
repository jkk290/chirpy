package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func validateChirp(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		Error       string `json:"error"`
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const charLimit = 140
	if len(params.Body) > charLimit {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	cleanedBody := profaneCheck(params.Body)

	respondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: cleanedBody,
	})

}

func profaneCheck(body string) string {
	// profaneWords := []string{
	// 	"kerfuffle",
	// 	"sharbert",
	// 	"fornax",
	// }

	// using map is O(1) instead of O(n) with slice
	profaneWords := map[string]bool{
		"kerfuffle": true,
		"sharbert":  true,
		"fornax":    true,
	}

	words := strings.Split(body, " ")
	final := []string{}
	for _, word := range words {
		converted := strings.ToLower(word)
		if profaneWords[converted] {
			final = append(final, "****")
		} else {
			final = append(final, word)
		}
	}
	return strings.Join(final, " ")
}
