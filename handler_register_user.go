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
	RefreshToken string	`json:"refresh_token,omitempty"`
	IsChirpyRed bool 	`json:"is_chirpy_red"`
}

type userParams struct {
	Email string `json:"email"`
	Password string `json:"password"`
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
		IsChirpyRed: user.IsChirpyRed,
	}

	respondWithJSON(w, http.StatusCreated, userRes)
}

