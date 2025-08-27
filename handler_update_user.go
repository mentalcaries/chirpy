package main

import (
	"encoding/json"
	"net/http"

	"github.com/mentalcaries/chirpy/internal/auth"
	"github.com/mentalcaries/chirpy/internal/database"
)

func (cfg *apiConfig)handleUpdateUser(w http.ResponseWriter, r *http.Request){
    decoder := json.NewDecoder(r.Body)
    params := userParams{}

    err := decoder.Decode(&params)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request", err)
        return
    }

    token, err := auth.GetBearerToken(r.Header)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "You must be logged in", err)
        return
    }
    hashedPassword, err := auth.HashPassword(params.Password)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request", err)
        return
    }

    userId, err := auth.ValidateJWT(token, cfg.jwtSecret)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "You are not authorized", err)
        return
    }   

    updatedUser, err := cfg.DBQueries.UpdateUser(r.Context(), database.UpdateUserParams{
        ID: userId,
        Email: params.Email,
        HashedPassword: hashedPassword,
    })

    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Could not update user", err)
    }

    respondWithJSON(w, http.StatusOK, User{
        CreatedAt: updatedUser.CreatedAt,
        ID: updatedUser.ID,
        Email: updatedUser.Email,
        UpdatedAt: updatedUser.UpdatedAt,
    })
}