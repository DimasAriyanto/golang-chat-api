package middleware

import (
	"context"
	"net/http"
	"strings"
	"fmt"
	"time"
	"errors"

	"github.com/golang-jwt/jwt/v4"
	"github.com/DimasAriyanto/golang-chat-api/config"
)

type ContextKey string

const UserIDKey ContextKey = "user_id"

func GenerateToken(userID int, secretKey string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func AuthMiddleware(cfg config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			fmt.Printf("Authorization Header: %s\n", authHeader)

			if authHeader == "" {
				http.Error(w, "Unauthorized - No token provided", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				http.Error(w, "Unauthorized - Invalid token format", http.StatusUnauthorized)
				return
			}

			claims := &jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(cfg.JWTSecret), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Unauthorized - Invalid token", http.StatusUnauthorized)
				return
			}

			userID, ok := (*claims)["user_id"].(float64)
			if !ok {
				http.Error(w, "Unauthorized - Invalid token data", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, int(userID))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func ParseToken(tokenString, secretKey string) (*jwt.MapClaims, error) {
	claims := &jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
