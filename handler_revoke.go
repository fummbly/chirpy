package main

import (
	"net/http"

	"github.com/fummbly/chirpy/internal/auth"
)

func (cfg *apiConfig) handleRevoke(w http.ResponseWriter, req *http.Request) {

	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't find token", err)
		return
	}

	_, err = cfg.database.SetRevokedAt(req.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't set revoke", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
