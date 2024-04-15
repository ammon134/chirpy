package main

import (
	"strconv"
	"time"

	"github.com/ammon134/chirpy/internal/database"
	"github.com/golang-jwt/jwt/v5"
)

// TODO: this func looks really clumsy, there should be a better way
func createJWTToken(cfg *apiConfig, user database.User, expiresInSeconds int) (string, error) {
	currentTime := time.Now().UTC()
	if expiresInSeconds == 0 || expiresInSeconds > 86400 {
		expiresInSeconds = 86400 // 24 hours
	}
	claims := jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(currentTime),
		ExpiresAt: jwt.NewNumericDate(currentTime.Add(time.Duration(expiresInSeconds) * time.Second)),
		Issuer:    "chirpy",
		Subject:   strconv.Itoa(user.ID),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(cfg.jwtSecret))
	if err != nil {
		return "", err
	}
	return ss, nil
}
