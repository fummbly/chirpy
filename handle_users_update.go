package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/fummbly/chirpy/internal/auth"
	"github.com/fummbly/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUserUpdate(w http.ResponseWriter, req *http.Request) {

	type Parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	requestToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Failed to get token from request", err)
		return
	}

	userID, err := auth.ValidateJWT(requestToken, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not validate JWT", err)
		return
	}

	params := Parameters{}

	decoder := json.NewDecoder(req.Body)

	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode request body", err)
		return
	}

	if params.Email == "" || params.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Email and password are required to update user", errors.New("Email and Password are required to update user"))
		return
	}

	hashed_password, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create hash password", err)
		return
	}

	user, err := cfg.database.UpdateUserEmailPassword(req.Context(), database.UpdateUserEmailPasswordParams{
		ID:             userID,
		HashedPassword: hashed_password,
		Email:          params.Email,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not update user", err)
		return
	}

	respondWithJson(w, http.StatusOK, User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})
}
