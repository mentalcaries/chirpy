package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/mentalcaries/chirpy/internal/auth"
)

func (cfg *apiConfig) verifyRefreshToken(w http.ResponseWriter, r *http.Request) {

	refreshToken, err := extractRefreshToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid request headers", err)
		return
	}

	tokenData, err := cfg.DBQueries.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Authorized required", err)
		return
	}

	if tokenData.ExpiresAt.Before(time.Now()) || tokenData.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	updatedJWT, err := auth.MakeJWT(tokenData.UserID, cfg.jwtSecret, time.Duration(time.Second*60))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not renew token", err)
	}

	respondWithJSON(w, http.StatusOK, struct {
		Token string `json:"token"`
	}{Token: updatedJWT})

}

func (cfg *apiConfig) revokeRefreshToken(w http.ResponseWriter, r *http.Request) {
    refreshToken, err := extractRefreshToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid request headers", err)
		return
	}

	tokenData, err := cfg.DBQueries.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Authorized required", err)
		return
	}

    err = cfg.DBQueries.RevokeRefreshToken(r.Context(), tokenData.Token)
    if err != nil{
        respondWithError(w, http.StatusBadRequest, "could not revoke token", err)
    }
    respondWithJSON(w, http.StatusNoContent, nil)
}

func extractRefreshToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return "", fmt.Errorf("invalid request header")
	}
	token := strings.Replace(authHeader, "Bearer ", "", 1)

    fmt.Println(token)
	return token, nil
}
