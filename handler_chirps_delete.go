package main

import (
	"net/http"

	"github.com/fummbly/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handleChirpDelete(w http.ResponseWriter, req *http.Request) {

	stringID := req.PathValue("chirp_id")

	id, err := uuid.Parse(stringID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not parse id", err)
		return
	}

	chirp, err := cfg.database.GetChirp(req.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Could not get chirp", err)
		return
	}

	requestToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not get auth token", err)
		return
	}

	userID, err := auth.ValidateJWT(requestToken, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not validate JWT", err)
		return
	}

	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "User does not own this chirp", err)
		return
	}

	err = cfg.database.DeleteChirpById(req.Context(), chirp.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not delete chirp because not authorized", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
