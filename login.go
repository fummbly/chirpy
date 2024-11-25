package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/fummbly/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, req *http.Request) {

	type Parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type Response struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	decoder := json.NewDecoder(req.Body)
	params := Parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to decode request", err)
		return
	}

	user, err := cfg.database.GetUser(req.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	respondWithJson(w, http.StatusOK, Response{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})

}
