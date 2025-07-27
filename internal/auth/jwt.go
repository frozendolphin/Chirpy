package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	
	signingKey := []byte(tokenSecret)

	claim := jwt.RegisteredClaims{
		Issuer: "chirpy",
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject: userID.String(),
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	token_str, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}
	
	return token_str, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {

	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claimsStruct, func(t *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.UUID{}, err
	}

	userid, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.UUID{}, err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}
	if issuer != "chirpy" {
		return uuid.Nil, errors.New("invalid issuer")
	}
	
	id, err := uuid.Parse(userid)
	if err != nil {
		return uuid.UUID{}, err
	}

	return id, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	
	unclean_token := headers.Get("Authorization")
	if unclean_token == "" {
		return "", errors.New("token doesn't exist")
	}

	values := strings.Fields(unclean_token)

	if values[0] != "Bearer" {
		return "", errors.New("no Bearer format")
	}

	return values[1], nil
}