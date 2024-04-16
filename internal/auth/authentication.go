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

type TokenType string

const (
	TokenTypeAccess  TokenType = "chirpy-access"
	TokenTypeRefresh TokenType = "chirpy-refresh"
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

func CreateJWT(jwtSecret string, userID int, duration time.Duration, tt TokenType) (string, error) {
	currentTime := time.Now().UTC()
	claims := jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(currentTime),
		ExpiresAt: jwt.NewNumericDate(currentTime.Add(duration)),
		Issuer:    string(tt),
		Subject:   strconv.Itoa(userID),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return ss, nil
}

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

func ValidateJWT(jwtSecret string, header http.Header) (*jwt.Token, error) {
	bearerToken, err := GetBearerToken(header)
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(
		bearerToken,
		&jwt.RegisteredClaims{},
		func(t *jwt.Token) (interface{}, error) { return []byte(jwtSecret), nil },
	)
	if err != nil {
		return nil, errors.New("token is invalid or has expired")
	}
	return token, nil
}

func GetBearerToken(header http.Header) (string, error) {
	auth := header.Get("Authorization")
	if auth == "" {
		return "", errors.New("no authorization token found")
	}
	bearerStr, tokenStr, found := strings.Cut(auth, " ")
	if !found || bearerStr != "Bearer" {
		return "", errors.New("malformed authorization header")
	}
	return tokenStr, nil
}
