package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/frozendolphin/Chirpy/internal/auth"
	"github.com/frozendolphin/Chirpy/internal/database"
	"github.com/google/uuid"
)

type usersInfo struct {
		Id uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email string `json:"email"`
		Token string `json:"token"`
		RefreshToken string `json:"refresh_token"`
		IsChiryRed bool `json:"is_chirpy_red"`
	}

func (cfg *apiConfig) createUsers(w http.ResponseWriter, r *http.Request) {
	
	type mail struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	params := mail{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	hashedp, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash the password", err)
		return
	}

	createUser_params := database.CreateUserParams {
		Email: params.Email,
		HashedPassword: hashedp,
	}

	user, err := cfg.db.CreateUser(r.Context(), createUser_params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	res := usersInfo {
		Id: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		IsChiryRed: user.IsChirpyRed,
	}

	respondWithJSON(w, http.StatusCreated, res)
}

func (cfg *apiConfig) resetHits(w http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "403 Forbidden: Access Denied", nil)
		return
	}

	err := cfg.db.DeleteAllUsers(req.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Request to delete to database Failed", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}


func (cfg *apiConfig) changeEmailPass(w http.ResponseWriter, r *http.Request) {

	type requirementreq struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find access token", err)
		return
	}

	requirement := requirementreq{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&requirement)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode the request data", err)
		return
	}

	if requirement.Email == "" || requirement.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Email or Password not found", err)
		return
	}

	user_id, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "jwt validation failed", err)
		return
	}

	hashpass, err := auth.HashPassword(requirement.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't hash the password", err)
		return
	}

	updateuserparams := database.UpdateUserParams {
		Email: requirement.Email,
		HashedPassword: hashpass,
		ID: user_id,
	}

	user, err := cfg.db.UpdateUser(r.Context(), updateuserparams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't update the database", err)
		return
	}

	res := usersInfo {
		Id: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		IsChiryRed: user.IsChirpyRed,
	}	

	respondWithJSON(w, http.StatusOK, res)
}