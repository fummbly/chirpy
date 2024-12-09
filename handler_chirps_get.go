package main

import (
	"net/http"

	"github.com/fummbly/chirpy/internal/database"
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
		respondWithError(w, http.StatusNotFound, "Could not find id", err)
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

	sort := req.URL.Query().Get("sort")

	var dbChirps []database.Chirp
	var err error

	if sort == "" || sort == "asc" {
		dbChirps, err = cfg.database.GetChirpsAsc(req.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps from database", err)
		}
	} else {
		dbChirps, err = cfg.database.GetChirpsDesc(req.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "couldn't get chirps from database", err)
			return
		}
	}

	authorID := uuid.Nil
	authorIDString := req.URL.Query().Get("author_id")
	if authorIDString != "" {
		authorID, err = uuid.Parse(authorIDString)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author id", err)
			return
		}
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		if authorID != uuid.Nil && dbChirp.UserID != authorID {
			continue
		}

		chirps = append(chirps, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			UserID:    dbChirp.UserID,
			Body:      dbChirp.Body,
		})
	}

	respondWithJson(w, http.StatusOK, chirps)
}
