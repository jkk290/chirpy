package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jkk290/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) createChirp(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
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

	newChirp := Chirp{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Body:      cleanedBody,
		UserId:    params.UserId,
	}

	_, err := cfg.dbQueries.CreateChirp(req.Context(), database.CreateChirpParams{
		ID:        newChirp.ID,
		CreatedAt: newChirp.CreatedAt,
		UpdatedAt: newChirp.UpdatedAt,
		Body:      newChirp.Body,
		UserID:    newChirp.UserId,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating chirp", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, newChirp)
}

func (cfg *apiConfig) getAllChirps(w http.ResponseWriter, req *http.Request) {
	chirps, err := cfg.dbQueries.GetChirps(req.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error getting all chirps", err)
		return
	}

	final := []Chirp{}
	for _, c := range chirps {
		chirp := Chirp{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body:      c.Body,
			UserId:    c.UserID,
		}
		final = append(final, chirp)
	}

	respondWithJSON(w, http.StatusOK, final)

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
