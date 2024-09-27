package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"strings"
	"time"
)

func ValidateJWT(tokenString string) (token *jwt.Token, err error) {
	secret := os.Getenv("JWT_SECRET")
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
}

func GenerateAccessToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(15 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
func ExtractAccessTokenFromRequest(r *http.Request) (string, error) {
	// First, check the Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		splitToken := strings.Split(authHeader, "Bearer ")
		if len(splitToken) == 2 {
			return splitToken[1], nil
		}
	}

	// If not in the header, check the cookie
	cookie, err := r.Cookie("access_token")
	if err == nil && cookie.Value != "" {
		return cookie.Value, nil
	}

	return "", fmt.Errorf("no valid access token found")
}
func ExtractSessionIDFromRequest(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return "", err
	}
	if cookie.Value == "" {
		return "", fmt.Errorf("empty session id")
	}
	return cookie.Value, nil
}
func GenerateSessionID(length int) (string, error) {
	bytes := make([]byte, length)
	// Fill the byte slice with random bytes
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	// Convert the byte slice to a hex string
	return hex.EncodeToString(bytes), nil
}
