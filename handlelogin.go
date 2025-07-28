package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/frozendolphin/Chirpy/internal/auth"
	"github.com/frozendolphin/Chirpy/internal/database"
)

func (cfg *apiConfig) loginUser(w http.ResponseWriter, r *http.Request) {

	type reqBody struct {
		Password string `json:"password"`
		Email string `json:"email"`
	}

	params := reqBody {}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	refresh_token, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create refresh token", err)
		return
	}

	rt_params := database.CreateRefreshTokenParams {
		Token: refresh_token,
		UserID: user.ID,
		ExpiresAt: time.Now().Add(24 * 60 * time.Hour),
	}

	rt, err := cfg.db.CreateRefreshToken(r.Context(), rt_params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't add refresh token in database", err)
		return
	}

	tokenexpiresIn := 3600 * time.Second
	token, err := auth.MakeJWT(user.ID, cfg.secret, tokenexpiresIn)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create a jwt", err)
		return
	}

	res := usersInfo {
		Id: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		Token: token,
		RefreshToken: rt.Token,
		IsChiryRed: user.IsChirpyRed,
	}

	respondWithJSON(w, http.StatusOK, res)
}