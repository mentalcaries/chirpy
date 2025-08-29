package main

import (
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	"github.com/mentalcaries/chirpy/internal/auth"
	"github.com/mentalcaries/chirpy/internal/database"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request){
	id := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(id)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "You must be logged in", err)
		return	
	}

	userId, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "You are not logged in", err)
		return
	}

	_ , err = cfg.DBQueries.DeleteChirpById(r.Context(), database.DeleteChirpByIdParams{
		ID: chirpID,
		UserID: userId,
	})
	if err != nil{
		if err == sql.ErrNoRows{
			respondWithError(w, http.StatusForbidden, "Post cannot be deleted", err)
			return
		}
		respondWithError(w, http.StatusBadRequest, "Could not delete post", err)
		return
	}
	respondWithJSON(w, http.StatusNoContent, nil)
}