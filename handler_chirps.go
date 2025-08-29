package main

import (
	"net/http"
	"sort"
	"time"

	"github.com/google/uuid"
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
	Body string `json:"body"`
	// UserId uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	authorParams := r.URL.Query().Get("author_id")
	userId, _ := uuid.Parse(authorParams)
	sortParams := r.URL.Query().Get("sort")

	var err error
	allChirps := []database.Chirp{}
	if userId != uuid.Nil {
		allChirps, err = cfg.DBQueries.GetChirpsByAuthor(r.Context(), userId)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Could not get posts", err)
			return
		}
	} else {

		allChirps, err = cfg.DBQueries.GetAllChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Could not get posts", err)
			return
		}
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

	if sortParams == "desc"{
		sort.Slice(chirps, func(i, j int) bool { return chirps[i].CreatedAt.After(chirps[j].CreatedAt) } )
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
