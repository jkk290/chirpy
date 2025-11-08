package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/jkk290/chirpy/internal/database"
)

func (cfg *apiConfig) upgradeUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	params := parameters{}
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "error parsing json", err)
		return
	}
	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	userId, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error parsing user id", err)
	}

	if err := cfg.dbQueries.UpdateChirpyRed(req.Context(), database.UpdateChirpyRedParams{
		ID:          userId,
		IsChirpyRed: true,
	}); err != nil {
		respondWithError(w, http.StatusNotFound, "user not found", err)
	}
	respondWithJSON(w, http.StatusNoContent, nil)
}
