package main

import (
	"net/http"

	"github.com/fummbly/chirpy/internal/auth"
)

func (cfg *apiConfig) handleRevoke(w http.ResponseWriter, req *http.Request) {

	type Response struct {
	}

	requestToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get request token", err)
		return
	}

	token, err := cfg.database.GetRefreshToken(req.Context(), requestToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token could not be found", err)
		return
	}

	err = cfg.database.SetRevokedAt(req.Context(), token.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to set revoked at time", err)
		return
	}

	respondWithJson(w, http.StatusNoContent, Response{})

}
