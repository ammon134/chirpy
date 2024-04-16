package auth

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return []byte{}, err
	}
	return hash, nil
}

func CheckPasswordHash(hashedPassword []byte, password string) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
}

func CreateJWT(jwtSecret string, userID, expiresInSeconds int) (string, error) {
	currentTime := time.Now().UTC()
	if expiresInSeconds == 0 || expiresInSeconds > 86400 {
		expiresInSeconds = 86400 // 24 hours
	}
	claims := jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(currentTime),
		ExpiresAt: jwt.NewNumericDate(currentTime.Add(time.Duration(expiresInSeconds) * time.Second)),
		Issuer:    "chirpy",
		Subject:   strconv.Itoa(userID),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return ss, nil
}

func ValidateJWT(jwtSecret string, r *http.Request) (*jwt.Token, error) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return nil, errors.New("no authorization token found")
	}
	bearerStr, tokenStr, found := strings.Cut(auth, " ")
	if !found || bearerStr != "Bearer" {
		return nil, errors.New("malformed authorization header")
	}

	token, err := jwt.ParseWithClaims(
		tokenStr,
		&jwt.RegisteredClaims{},
		func(t *jwt.Token) (interface{}, error) { return []byte(jwtSecret), nil },
	)
	if err != nil {
		return nil, errors.New("token is invalid or has expired")
	}
	return token, nil
}