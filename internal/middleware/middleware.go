package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/alexnakagama/go-file-storage/internal/auth"
)

type contextKey string

const UserClaimsKey = contextKey("user_claims")

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer") {
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer")

		claims, err := auth.VerifyToken(token)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserClaims(r *http.Request) (*auth.CustomClaims, bool) {
	claims, ok := r.Context().Value(UserClaimsKey).(*auth.CustomClaims)
	return claims, ok
}
