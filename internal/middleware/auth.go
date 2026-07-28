package middleware

import (
	"context"
	"net/http"
	"slices"
	"strings"

	"github.com/Nuttachai-K/cafe-inventory-management/internal/model"
	"github.com/Nuttachai-K/cafe-inventory-management/internal/utils"
)

type contextKey int

const claimsKey contextKey = iota

func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		const prefix = "Bearer "
		if !strings.HasPrefix(authHeader, prefix) {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		claims, err := utils.ParseToken(strings.TrimPrefix(authHeader, prefix))
		if err != nil {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), claimsKey, claims)
		next(w, r.WithContext(ctx))
	}
}

func GetClaims(ctx context.Context) (*utils.Claims, bool) {
	claims, ok := ctx.Value(claimsKey).(*utils.Claims)
	return claims, ok
}

func RequireRole(roles ...model.UserRole) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			claims, ok := GetClaims(r.Context())
			if !ok {
				http.Error(w, "invalid or expired token", http.StatusUnauthorized)
				return
			}

			if slices.Contains(roles, model.UserRole(claims.Role)) {
				next(w, r)
				return
			}
			http.Error(w, "permission denied", http.StatusForbidden)
		}
	}
}
