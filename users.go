package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handleAddUser(w http.ResponseWriter, req *http.Request) {

	type Paramaters struct {
		Email string `json:"email"`
	}

	type Response struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	decoder := json.NewDecoder(req.Body)
	params := Paramaters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to decode request", err)
		return
	}

	user, err := cfg.database.CreateUser(req.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create user", err)
		return
	}

	respondWithJson(w, http.StatusCreated, Response{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})

}
