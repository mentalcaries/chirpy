package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/mentalcaries/chirpy/internal/auth"
)

type EventWebhook struct{
    Event string `json:"event"`
    Data struct{
        UserId uuid.UUID `json:"user_id"`
    } `json:"data"`
}

func (cfg *apiConfig)handlerUpgradeUserSubscription(w http.ResponseWriter, r *http.Request){
    apiKey, err := auth.GetAPIKey(r.Header)

    if err != nil || apiKey != cfg.apiKey {
        respondWithError(w, http.StatusUnauthorized, "Invalid provider key", err)
        return
    }

    
    decoder := json.NewDecoder(r.Body)
    reqBody := EventWebhook{}
    decoder.Decode(&reqBody)

    if reqBody.Event != "user.upgraded" {
        respondWithJSON(w, http.StatusNoContent, nil)
    }


    if reqBody.Event == "user.upgraded" {
        _, err := cfg.DBQueries.UpgradeUserSubscription(r.Context(), reqBody.Data.UserId)

        if err != nil {
            if err == sql.ErrNoRows{
                respondWithError(w, http.StatusNotFound, "Invalid user", err)
                return
            }
            respondWithError(w, http.StatusBadRequest, "Could not upgrade user", err)
            return
        }

        w.WriteHeader(http.StatusNoContent)
    }
}