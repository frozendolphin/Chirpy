package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
)

func MakeRefreshToken() (string, error) {

	key := make([]byte, 32)

	rand.Read(key)

	encodedStr := hex.EncodeToString(key)

	return encodedStr, nil
}

func GetBearerRefreshToken(headers http.Header) (string, error) {
	
	unclean_token := headers.Get("Authorization")
	if unclean_token == "" {
		return "", errors.New("refresh token doesn't exist")
	}

	values := strings.Fields(unclean_token)

	if values[0] != "Bearer" {
		return "", errors.New("no Bearer format")
	}

	return values[1], nil
}