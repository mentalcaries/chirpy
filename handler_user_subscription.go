package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type EventWebhook struct{
    Event string `json:"event"`
    Data struct{
        UserId uuid.UUID `json:"user_id"`
    } `json:"data"`
}

func (cfg *apiConfig)handlerUpgradeUserSubscription(w http.ResponseWriter, r *http.Request){
    decoder := json.NewDecoder(r.Body)
    
    reqBody := EventWebhook{}
    decoder.Decode(&reqBody)

    fmt.Println(reqBody)

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

        respondWithJSON(w, http.StatusNoContent, nil)
    }
}