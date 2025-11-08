package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jkk290/chirpy/internal/auth"
	"github.com/jkk290/chirpy/internal/database"
)

type User struct {
	ID             uuid.UUID `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"hashed_password"`
	IsChirpyRed    bool      `json:"is_chirpy_red"`
	Token          string    `json:"token"`
	RefreshToken   string    `json:"refresh_token"`
}

func (cfg *apiConfig) createUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	hashedpw, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error hashing password", err)
		return
	}

	dbUser, err := cfg.dbQueries.CreateUser(req.Context(), database.CreateUserParams{
		ID:             uuid.New(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Email:          params.Email,
		HashedPassword: hashedpw,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't create new user", err)
		return
	}

	newUser := User{
		ID:          dbUser.ID,
		CreatedAt:   dbUser.CreatedAt,
		UpdatedAt:   dbUser.UpdatedAt,
		Email:       dbUser.Email,
		IsChirpyRed: dbUser.IsChirpyRed,
	}

	respondWithJSON(w, http.StatusCreated, newUser)
}

func (cfg *apiConfig) updateUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	tokenString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "token invalid", err)
		return
	}

	userId, err := auth.ValidateJWT(tokenString, cfg.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "token invalid", err)
		return
	}

	params := parameters{}
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "error parsing json", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error hashing password", err)
		return
	}

	dbUpdatedUser, err := cfg.dbQueries.UpdateUser(req.Context(), database.UpdateUserParams{
		ID:             userId,
		Email:          params.Email,
		HashedPassword: hashedPassword,
		UpdatedAt:      time.Now(),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error updating user", err)
		return
	}

	updatedUser := User{
		ID:          dbUpdatedUser.ID,
		Email:       dbUpdatedUser.Email,
		CreatedAt:   dbUpdatedUser.CreatedAt,
		UpdatedAt:   dbUpdatedUser.UpdatedAt,
		IsChirpyRed: dbUpdatedUser.IsChirpyRed,
	}

	respondWithJSON(w, http.StatusOK, updatedUser)
}
