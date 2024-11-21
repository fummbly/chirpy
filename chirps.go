package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/fummbly/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handleAddChirp(w http.ResponseWriter, req *http.Request) {

	type Parameters struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}

	type Response struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserId    uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(req.Body)
	params := Parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to decode request", err)
		return
	}

	statusCode, validatedBody := validateChrip(params.Body)
	if statusCode > 200 {
		respondWithError(w, statusCode, "Not a valid chirp", nil)
		return
	}

	chirp, err := cfg.database.CreateChirp(req.Context(), database.CreateChirpParams{
		Body:   validatedBody,
		UserID: params.UserId,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create chirp", err)
		return
	}

	respondWithJson(w, http.StatusCreated, Response{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
	})
}

func (cfg *apiConfig) handleGetChirps(w http.ResponseWriter, req *http.Request) {

	type Response struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserId    uuid.UUID `json:"user_id"`
	}

	chirps, err := cfg.database.GetChirps(req.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get chirps", err)
		return
	}

	res := []Response{}

	for _, chirp := range chirps {
		res = append(res, Response{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserId:    chirp.UserID,
		})
	}

	respondWithJson(w, http.StatusOK, res)

}

func (cfg *apiConfig) handleGetChirp(w http.ResponseWriter, req *http.Request) {

	type Response struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	stringID := req.PathValue("chirp_id")

	id, err := uuid.Parse(stringID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to convert to uuid", err)
		return
	}

	chirp, err := cfg.database.GetChirp(req.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not find id", err)
	}

	respondWithJson(w, http.StatusOK, Response{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})

}
