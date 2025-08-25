package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const passwordCost = 10

func HashPassword(password string) (string, error) {

	if len(password) < 4 {
		return "", fmt.Errorf("password is too short")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), passwordCost)
	if err != nil {
		return "", fmt.Errorf("could not hash password")
	}

	return string(hashedPassword), nil
}

func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	if err != nil {
		return fmt.Errorf("incorrect password. try again")
	}
	return nil
}
