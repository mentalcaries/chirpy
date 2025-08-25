package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)


func MakeJWT(userId uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userId.String(),
	})

	signedToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", fmt.Errorf("could not generate token")
	}

	return signedToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {

	claimsStruct := jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenString, &claimsStruct, func(t *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid token: %v", err)
	}

	userIDString, err := token.Claims.GetSubject() 
	if err != nil {
		return uuid.Nil, err
	}

	parsedUUID, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.Nil, err
	}
	
	return parsedUUID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return "", fmt.Errorf("invalid request header")
	}

	token := strings.Replace(authHeader, "Bearer ", "", 1)

    return token, nil
}
