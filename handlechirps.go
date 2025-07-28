package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/frozendolphin/Chirpy/internal/auth"
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
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find token", err)
		return
	}

	user_id, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "jwt validation failed", err)
		return
	}

	cleaned, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	cparams := database.CreateChirpParams{
		Body: cleaned,
		UserID: user_id,
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

	
	var all_chirps []database.Chirp
	var err error

	s := r.URL.Query().Get("author_id")
	sortchirps := r.URL.Query().Get("sort")

	if s != "" {

		a_id, err := uuid.Parse(s)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "id cannot be parsed to uuid format", err)
			return
		}

		all_chirps, err = cfg.db.GetChirpsFromAuthor(r.Context(), a_id)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "chirp cannot be found", err)
			return
		}

	} else {

		all_chirps, err = cfg.db.GetAllChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "db request to get all chirps Failed", err)
			return
		}

	}

	if sortchirps == "desc" {
		sort.Slice(all_chirps, func(i, j int) bool {return all_chirps[i].CreatedAt.After(all_chirps[j].CreatedAt) })
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

func (cfg *apiConfig) deleteAChirp(w http.ResponseWriter, r *http.Request) {

	c_id := r.PathValue("chirpID")
	if c_id == "" {
		respondWithError(w, http.StatusBadRequest, "no such wildcard in path", nil)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find access token", err)
		return
	}

	user_id, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "jwt validation failed", err)
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

	if chirp.UserID != user_id {
		respondWithError(w, http.StatusForbidden, "given chirp is not yours", err)
		return
	}

	err = cfg.db.DeleteAChirp(r.Context(), u_id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "couldn't find id in the server", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}