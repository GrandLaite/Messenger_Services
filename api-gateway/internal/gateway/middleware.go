package gateway

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

type ctxKey int

var userContextKey ctxKey = 0

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Malformed token", http.StatusUnauthorized)
			return
		}
		tokenStr := parts[1]

		secret := os.Getenv("GATEWAY_JWT_SECRET")
		if secret == "" {
			secret = "default_secret"
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid claims", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserClaims(r *http.Request) jwt.MapClaims {
	val := r.Context().Value(userContextKey)
	if claims, ok := val.(jwt.MapClaims); ok {
		return claims
	}
	return nil
}
