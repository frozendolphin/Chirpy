package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/frozendolphin/Chirpy/internal/database"
	"github.com/google/uuid"
)

type chirpInfo struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time	`json:"updated_at"`
	Body      string	`json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) createChirps(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	cleaned, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	cparams := database.CreateChirpParams{
		Body: cleaned,
		UserID: params.UserId,
	}
	
	chirp, err := cfg.db.CreateChirp(r.Context(), cparams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "db request to create chirp Failed", err)
		return
	}

	res := chirpInfo {
		Id: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserId: chirp.UserID,
	}

	respondWithJSON(w, http.StatusCreated, res)
}

func validateChirp(body string) (string, error) {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("chirp is too long")
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleaned := getCleanedBody(body, badWords)
	return cleaned, nil
}

func getCleanedBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}

func (cfg *apiConfig) getAllChirps(w http.ResponseWriter, r *http.Request) {

	all_chirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "db request to get all chirps Failed", err)
		return
	}

	var chirp_list []chirpInfo
	var single_chirp chirpInfo

	for _, chirp := range all_chirps {
		single_chirp = chirpInfo{
			Id: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			UserId: chirp.UserID,
		}
		chirp_list = append(chirp_list, single_chirp)
	}
	
	respondWithJSON(w, http.StatusOK, chirp_list)
}

func (cfg *apiConfig) getAChirp(w http.ResponseWriter, r *http.Request) {

	c_id := r.PathValue("chirpID")
	if c_id == "" {
		respondWithError(w, http.StatusBadRequest, "no such wildcard in path", nil)
		return
	}

	u_id, err := uuid.Parse(c_id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't convert id into uuid", err)
		return
	}

	chirp, err := cfg.db.GetAChirp(r.Context(), u_id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "couldn't find id in the server", err)
		return
	}

	res := chirpInfo {
		Id: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserId: chirp.UserID,
	}

	respondWithJSON(w, http.StatusOK, res)
} 