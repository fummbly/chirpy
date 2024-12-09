package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {

	auth := headers.Get("Authorization")
	if auth == "" {
		return "", ErrNoAuthHeaderIncluded
	}

	splitAuth := strings.Split(auth, " ")

	if len(splitAuth) != 2 {
		return "", errors.New("Incorrect format for api key")
	}

	if splitAuth[0] != "ApiKey" {
		return "", errors.New("Request missing apikey")
	}

	apiKey := splitAuth[1]

	return apiKey, nil

}
