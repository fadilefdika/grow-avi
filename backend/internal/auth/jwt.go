package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	NPK        string `json:"npk"`
	UserName   string `json:"user_name"`
	Department string `json:"department"`
	Role       string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(npk, userName, department, role string) (string, error) {
	ttlStr := os.Getenv("JWT_ACCESS_TTL")
	if ttlStr == "" {
		ttlStr = "24h"
	}
	ttl, err := time.ParseDuration(ttlStr)
	if err != nil {
		ttl = 24 * time.Hour
	}

	expirationTime := time.Now().Add(ttl)

	claims := &Claims{
		NPK:        npk,
		UserName:   userName,
		Department: department,
		Role:       role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := []byte(os.Getenv("JWT_SECRET"))

	return token.SignedString(secret)
}

func GenerateRefreshToken() (string, string) {
	// Generate random string
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, 64)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	plainToken := string(b)

	// Hash it
	hasher := sha256.New()
	hasher.Write([]byte(plainToken))
	hashedToken := hex.EncodeToString(hasher.Sum(nil))

	return plainToken, hashedToken
}
