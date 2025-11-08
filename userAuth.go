package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jkk290/chirpy/internal/auth"
	"github.com/jkk290/chirpy/internal/database"
)

func (cfg *apiConfig) loginUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
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

	expiresIn := time.Duration(3600)

	tokenString, err := auth.MakeJWT(dbUser.ID, cfg.tokenSecret, expiresIn*time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error generating token", err)
	}

	hexString, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error generating refresh token", err)
	}

	dbRefreshToken, err := cfg.dbQueries.CreateRefreshToken(req.Context(), database.CreateRefreshTokenParams{
		Token:     hexString,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    dbUser.ID,
		ExpiresAt: time.Now().Add(1440 * time.Hour),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error storing refresh token to db", err)
	}

	newUser := User{
		ID:           dbUser.ID,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
		Email:        dbUser.Email,
		IsChirpyRed:  dbUser.IsChirpyRed,
		Token:        tokenString,
		RefreshToken: dbRefreshToken.Token,
	}
	respondWithJSON(w, http.StatusOK, newUser)
}

func (cfg *apiConfig) refreshToken(w http.ResponseWriter, req *http.Request) {
	type returnVal struct {
		Token string `json:"token"`
	}

	tokenString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "token invalid", err)
		return
	}
	refreshToken, err := cfg.dbQueries.GetRefreshToken(req.Context(), tokenString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "token invalid", err)
		return
	}
	if refreshToken.RevokedAt.Valid || time.Now().After(refreshToken.ExpiresAt) {
		respondWithError(w, http.StatusUnauthorized, "refresh token revoked", nil)
		return
	}

	user, err := cfg.dbQueries.GetUserFromRefreshToken(req.Context(), tokenString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error getting user info", err)
		return
	}

	newToken, err := auth.MakeJWT(user.ID, cfg.tokenSecret, 1*time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error generating new token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, returnVal{
		Token: newToken,
	})
}

func (cfg *apiConfig) revokeToken(w http.ResponseWriter, req *http.Request) {
	tokenString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "token invalid", err)
		return
	}

	if err := cfg.dbQueries.RevokeToken(req.Context(), database.RevokeTokenParams{
		Token:     tokenString,
		UpdatedAt: time.Now(),
		RevokedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}); err != nil {
		respondWithError(w, http.StatusInternalServerError, "error revoking token", err)
	}
	w.WriteHeader(http.StatusNoContent)
}
