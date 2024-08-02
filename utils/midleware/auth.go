package middleware

import (
	"HarvestBox/utils"
	"context"
	"fmt"
	"net/http"
	"strings"
)

type contextKey string

const TokenClaimsKey contextKey = "tokenClaims"

func ValidateTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := extractTokenFromHeader(r)
		if tokenString == "" {
			http.Error(w, "Missing authorization token", http.StatusUnauthorized)
			return
		}

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid authorization token", http.StatusUnauthorized)
			return
		}

		email := claims.Audience

		jwtToken, err := utils.RenewToken(email, "admin")
		if err != nil {
			fmt.Println("failed to renew token")
			return
		}
		fmt.Println(jwtToken)
		w.Header().Set("Authorization", "Bearer "+jwtToken)
		w.Header().Set("Content-Type", "application/json")

		// Add the token claims to the request context for downstream handlers to access
		ctx := r.Context()
		ctx = context.WithValue(ctx, TokenClaimsKey, claims)

		// Call the next handler with the modified context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractTokenFromHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	// Check if the Authorization header has the format "Bearer <token>"
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}

	return parts[1]
}
