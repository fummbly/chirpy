package main

import (
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

	user, err := cfg.database.GetUserFromRefreshToken(req.Context(), requestToken)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't get user from token", err)
		return
	}

	accessToken, err := auth.MakeJWT(
		user.ID,
		cfg.secret,
		time.Hour,
	)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate and create access token", err)
		return
	}

	respondWithJson(w, http.StatusOK, Response{
		Token: accessToken,
	})

}
