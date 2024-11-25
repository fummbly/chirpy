package main

import (
  "encoding/json"
  "net/http"
  "time"
  "strings"
  "slices"

  "github.com/google/uuid"
  "github.com/fummbly/chirpy/internal/database"

)

type Chirp struct {
  ID uuid.UUID `json:"id"`
  CreatedAt time.Time `json:"created_at"`
  UpdatedAt time.Time `json:"updated_at"`
  UserID uuid.UUID `json:"user_id"`
  Body string `json:"body"`
}

func (cfg *apiConfig) handleAddChirp(w http.ResponseWriter, req *http.Request) {

	type Parameters struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
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

	respondWithJson(w, http.StatusCreated, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}


func validateChrip(str string) (int, string) {
  if len(str) > 140 {
    return 400, ""

  }

	profainity := []string{"kerfuffle", "sharbert", "fornax"}

	words := strings.Split(str, " ")
	for i, word := range words {
		if slices.Contains(profainity, strings.ToLower(word)) {
			words[i] = "****"
		}
	}

	newStr := strings.Join(words, " ")

	return 200, newStr

}


