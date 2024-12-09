package main

import (
	"encoding/json"
	"net/http"

	"github.com/fummbly/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handleUserUpgrade(w http.ResponseWriter, req *http.Request) {

	type Parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	params := Parameters{}

	reqAPIKey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Failed to get Api key from request", err)
		return
	}

	if reqAPIKey != cfg.polka_key {
		respondWithError(w, http.StatusUnauthorized, "Request api key doesn't match key stored", nil)
		return
	}

	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to decode request body", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not parse user id", err)
		return
	}

	user, err := cfg.database.GetUserByID(req.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Could not find user", err)
		return
	}

	_, err = cfg.database.UpgradeUser(req.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not upgrade user", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
