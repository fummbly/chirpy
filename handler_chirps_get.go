package main

import (
	"net/http"
  "github.com/google/uuid"
)

func (cfg *apiConfig) handleGetChirp(w http.ResponseWriter, req *http.Request) {

	stringID := req.PathValue("chirp_id")

	id, err := uuid.Parse(stringID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to convert to uuid", err)
		return
	}

	chirp, err := cfg.database.GetChirp(req.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not find id", err)
    return
	}

	respondWithJson(w, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})

}


func (cfg *apiConfig) handleGetChirps(w http.ResponseWriter, req *http.Request) {

	chirps, err := cfg.database.GetChirps(req.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get chirps", err)
		return
	}


	res := []Chirp{}

	for _, chirp := range chirps {
		res = append(res, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	respondWithJson(w, http.StatusOK, res)

}



