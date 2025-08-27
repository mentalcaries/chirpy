package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mentalcaries/chirpy/internal/auth"
	"github.com/mentalcaries/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

type chirpParams struct {
	Body   string    `json:"body"`
	// UserId uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	params := chirpParams{}

    token, err:= auth.GetBearerToken(req.Header)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "You must be logged in", err)
        return
    }

    userId, err := auth.ValidateJWT(token, cfg.jwtSecret)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "Invalid user", err)
        return
    }


	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request", err)
		return
	}

	chirp, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	newChirp, err := cfg.DBQueries.CreateChirp(req.Context(), database.CreateChirpParams{
		Body:   chirp,
		UserID: userId,
	})

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not create chirp", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        newChirp.ID,
		CreatedAt: newChirp.CreatedAt,
		UpdatedAt: newChirp.UpdatedAt,
		Body:      newChirp.Body,
		UserId:    newChirp.UserID,
	})

}

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	allChirps, err := cfg.DBQueries.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not get posts", err)
	}

	chirps := []Chirp{}

	for _, chirp := range allChirps {
		formattedChirp := Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserId:    chirp.UserID,
		}

		chirps = append(chirps, formattedChirp)
	}

	respondWithJSON(w, http.StatusOK, chirps)

}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(id)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
        return
	}

	chirp, err := cfg.DBQueries.GetChirpById(r.Context(), chirpID)
    if err != nil {
        respondWithError(w, http.StatusNotFound, "Invalid ID", err)
    }

	chirpRes := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
	}
	respondWithJSON(w, http.StatusOK, chirpRes)
}
