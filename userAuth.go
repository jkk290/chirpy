package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jkk290/chirpy/internal/auth"
)

func (cfg *apiConfig) loginUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decoding parameters", err)
		return
	}

	dbUser, err := cfg.dbQueries.GetUserByEmail(req.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "user not found", err)
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, dbUser.HashedPassword)
	if err != nil || !match {
		respondWithError(w, http.StatusUnauthorized, "error verifying password", err)
		return
	}

	expiresDefault := 3600
	expiresIn := time.Duration(max(expiresDefault, params.ExpiresInSeconds))

	tokenString, err := auth.MakeJWT(dbUser.ID, cfg.tokenSecret, expiresIn*time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error generating token", err)
	}

	newUser := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
		Token:     tokenString,
	}
	respondWithJSON(w, http.StatusOK, newUser)
}
