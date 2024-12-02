package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fummbly/chirpy/internal/auth"
)

func (cfg *apiConfig) handleRefresh(w http.ResponseWriter, req *http.Request) {

	type Response struct {
		Token string `json:"token"`
	}

	requestToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token was not found", err)
		return
	}

	token, err := cfg.database.GetRefreshToken(req.Context(), requestToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token not found", err)
		return
	}

	if token.ExpiresAt.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "Token has expired", err)
		return
	}

	if token.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Token was revoked", err)
		return
	}

	fmt.Println(token.RevokedAt)

	user, err := cfg.database.GetUserFromRefreshToken(req.Context(), token.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get user by refresh token", err)
		return
	}

	newJWT, err := auth.MakeJWT(user.ID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create new jwt", err)
		return
	}

	respondWithJson(w, http.StatusOK, Response{
		Token: newJWT,
	})

}
