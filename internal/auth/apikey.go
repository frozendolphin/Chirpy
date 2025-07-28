package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {

	unclean_apikey := headers.Get("Authorization")
	if unclean_apikey == "" {
		return "", errors.New("apikey doesn't exist")
	}

	values := strings.Fields(unclean_apikey)

	if values[0] != "ApiKey" {
		return "", errors.New("no ApiKey format")
	}

	return values[1], nil
	
}