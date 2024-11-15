package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
)

func removeProfainity(str string) string {

	profainity := []string{"kerfuffle", "sharbert", "fornax"}

	words := strings.Split(str, " ")
	for i, word := range words {
		if slices.Contains(profainity, strings.ToLower(word)) {
			words[i] = "****"
		}
	}

	newStr := strings.Join(words, " ")

	return newStr

}

func handleValidate(w http.ResponseWriter, req *http.Request) {

	type parameters struct {
		Body string `json:"body"`
	}
	type returnVal struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const chripLength = 140

	if len(params.Body) > chripLength {
		respondWithError(w, http.StatusBadRequest, "Chirp too long", nil)
		return
	}

	respondWithJson(w, http.StatusOK, returnVal{
		CleanedBody: removeProfainity(params.Body),
	})

}
