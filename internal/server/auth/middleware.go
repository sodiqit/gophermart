package auth

import (
	"context"
	"net/http"
	"strings"
)

const (
	AuthHeaderKey = "Authorization"
	BearerPrefix  = "Bearer "
)

func JWTAuth(tokenService TokenService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get(AuthHeaderKey)
			if authHeader == "" {
				http.Error(w, "Authorization header is required", http.StatusUnauthorized)
				return
			}

			if !strings.HasPrefix(authHeader, BearerPrefix) {
				http.Error(w, "Authorization header must start with 'Bearer'", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, BearerPrefix)
			if tokenString == "" {
				http.Error(w, "Token is not provided", http.StatusUnauthorized)
				return
			}

			claims, err := tokenService.Validate(tokenString)

			if err != nil {
				http.Error(w, "", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), CLAIMS_CONTEXT_KEY, claims.TokenUser)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
