package main

import (
	"encoding/json"
	"net/http"

	"github.com/frozendolphin/Chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) upgradeUserChirpyRed(w http.ResponseWriter, r *http.Request) {
	
	type event struct {
		Event string `json:"event"`
		Data struct {
			User_id string `json:"user_id"`
		} `json:"data"`
	}

	param := event{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&param)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode response", err)
	}

	apikey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "no apikey in header", err)
		return
	}

	if cfg.polka_key != apikey {
		respondWithError(w, http.StatusUnauthorized, "apikey didn't match", err)
		return
	}

	if param.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	u_id, err := uuid.Parse(param.Data.User_id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't convert id into uuid", err)
		return
	}

	err = cfg.db.UpgradeUserChirpyRed(r.Context(), u_id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "couldn't find the user", err)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}