package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(header http.Header) (string, error) {
	authHeader := header.Get("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "ApiKey") {
		return "", fmt.Errorf("invalid api key")
	}

	apiKey := strings.Replace(authHeader, "ApiKey ", "", 1)

	return apiKey, nil
}
