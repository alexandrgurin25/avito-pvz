package auth

import (
	"avito-pvz/internal/constants"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func generateJWT(id int, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"role":    role,
		"iat":     time.Now().Add(constants.ExpirationTime).Unix(),
	})

	secretKeyString := os.Getenv("AUTH_SECRET_KEY")

	tokenString, err := token.SignedString([]byte(secretKeyString))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
