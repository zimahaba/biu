package main

import (
	"crypto/rand"
	"encoding/base64"
	"log/slog"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

const tokenExpiration = 60 * time.Minute

var JwtKey = []byte(os.Getenv("JWT_KEY"))

func GenerateToken(username string) (string, error) {
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExpiration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString(JwtKey)
}

func GenerateRandomToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		slog.Error(err.Error())
		return "", err
	}

	return base64.URLEncoding.EncodeToString(token), nil
}
