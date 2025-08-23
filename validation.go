package main

import (
	"errors"
	"slices"
	"strings"
)


func validateChirp(chirp string)(string, error) {
	// decoder := json.NewDecoder(body)

	// params := parameters{}

	// err := decoder.Decode(&params)
	// if err != nil {
	// 	respondWithError(w, 500, "Invalid paramters", err)
	// 	return
	// }

	const maxChirpLength = 140
	if len(chirp) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}

    filteredRes := filterProfanity(chirp)
	return filteredRes, nil
	// respondWithJSON(w, http.StatusOK, returnVal{Filtered: filteredRes})

}

func filterProfanity(s string) string {
    profanities := []string{ "kerfuffle", "sharbert", "fornax"}
    censored := "****"

    if len(s) < 1 {
        return ""
    }
    cleanedWords := []string{}
    words := strings.Split(s, " ")

    for _, word := range words {
        if slices.Contains(profanities, strings.ToLower(word)){
            cleanedWords = append(cleanedWords, censored)
            continue
        }
        cleanedWords = append(cleanedWords, word)
    }
    return strings.Join(cleanedWords, " ")
}