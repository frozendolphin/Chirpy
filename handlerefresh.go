package main

import (
	"net/http"
	"time"

	"github.com/frozendolphin/Chirpy/internal/auth"
)

func (cfg *apiConfig) newRefresh(w http.ResponseWriter, r *http.Request) {
	
	type respstruct struct {
		Token string `json:"token"`
	}

	rtoken, err := auth.GetBearerRefreshToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find refresh token", err)
		return
	}

	rt_info, err := cfg.db.GetRefreshToken(r.Context(), rtoken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find refresh token", err)
		return
	}

	if rt_info.ExpiresAt.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "Refresh token expired", nil)
		return
	}

	if rt_info.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Refresh token has been revoked", nil)
		return
	}

	expiresIn := 1 * time.Hour
	new_token, err := auth.MakeJWT(rt_info.UserID, cfg.secret, expiresIn)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create new jwt", err)
		return
	}

	res := respstruct {
		Token: new_token,
	}

	respondWithJSON(w, http.StatusOK, res)
}

func (cfg *apiConfig) revokeRefresh(w http.ResponseWriter, r *http.Request) {

	rtoken, err := auth.GetBearerRefreshToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find refresh token", err)
		return
	}

	err = cfg.db.RevokeRefreshToken(r.Context(), rtoken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find refresh token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}