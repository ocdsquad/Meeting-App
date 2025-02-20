package helper

import (
	"E-Meeting/configs"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	UserID  int    `json:"user_id"`
	IsAdmin bool   `json:"is_admin"`
	Email   string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID int, isAdmin bool, language string, email string) (string, error) {
	claims := jwt.MapClaims{
		"email":    email,
		"user_id":  userID,
		"is_admin": isAdmin,
		"language": language,
		"exp":      time.Now().Add(configs.Token.TokenExpiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(configs.Token.JWTSecret)
}
