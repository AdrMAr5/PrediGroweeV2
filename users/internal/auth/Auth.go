package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"time"
)

//func CreateAndSendAuthCookie(rw http.ResponseWriter, userID int) (string, error) {
//	secret := []byte(os.Getenv("JWT_SECRET"))
//	token, err := CreateJWT(secret, userID)
//	if err != nil {
//		http.Error(rw, "Unable to create token", http.StatusInternalServerError)
//		return "", err
//	}
//	http.SetCookie(rw, &http.Cookie{
//		Name:  "Authorization",
//		Value: token,
//	})
//	return token, nil
//}

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
		"exp": time.Now().Add(15 * time.Minute).Unix(), // Short-lived access token
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func GenerateRefreshToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(), // Longer-lived refresh token
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func ExtractRefreshTokenFromRequest(r *http.Request) (string, error) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
