package main

import (
	"slices"
	"strings"
)

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

