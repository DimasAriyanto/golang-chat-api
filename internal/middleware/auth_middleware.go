package middleware

import (
	"context"
	"net/http"
	"strings"
	"fmt"
	"time"
	"errors"

	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("secret_key")

type ContextKey string

const UserIDKey ContextKey = "user_id"

func GenerateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}


func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		fmt.Printf("Authorization Header: %s\n", authHeader)

		if authHeader == "" {
			http.Error(w, "Unauthorized - No token provided", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		fmt.Printf("Extracted Token String: %s\n", tokenString)

		if tokenString == authHeader {
			http.Error(w, "Unauthorized - Invalid token format", http.StatusUnauthorized)
			return
		}

		claims := &jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		fmt.Printf("Parsed Token Claims: %+v\n", claims)

		if err != nil || !token.Valid {
			fmt.Println("Error parsing token:", err)
			http.Error(w, "Unauthorized - Invalid token", http.StatusUnauthorized)
			return
		}

		userID, ok := (*claims)["user_id"].(float64)
		if !ok {
			http.Error(w, "Unauthorized - Invalid token data", http.StatusUnauthorized)
			return
		}

		fmt.Printf("Extracted user_id: %d\n", int(userID))

		ctx := context.WithValue(r.Context(), UserIDKey, int(userID))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ParseToken(tokenString string) (*jwt.MapClaims, error) {
	claims := &jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}