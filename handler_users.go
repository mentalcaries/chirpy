package main

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"time"

	"github.com/google/uuid"
	"github.com/mentalcaries/chirpy/internal/auth"
	"github.com/mentalcaries/chirpy/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token	  string 	`json:"token,omitempty"`
}

type userParams struct {
	Email string `json:"email"`
	Password string `json:"password"`
	ExpiresIn int `json:"expires_in_seconds"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := userParams{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not decode request", err)
		return
	}

    _, err = mail.ParseAddress(params.Email)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request", err)
		return
    }

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Password is not valid", err)
		return
	}

	user, err := cfg.DBQueries.CreateUser(r.Context(), database.CreateUserParams{
		Email: params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not create user", err)
		return
	}

	userRes := User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	}

	respondWithJSON(w, http.StatusCreated, userRes)
}

func (cfg *apiConfig)handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := userParams{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request", err)
		return
	}

	user, err := cfg.DBQueries.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Invalid credentials", err)
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	if params.ExpiresIn == 0 || params.ExpiresIn > int(time.Second * 3600){
		params.ExpiresIn = int(time.Second * 3600)
	}

	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Duration(params.ExpiresIn))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not generate token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		Token: token,
	})

}