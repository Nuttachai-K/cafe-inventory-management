package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nuttachai-K/cafe-inventory-management/internal/model"
	"github.com/Nuttachai-K/cafe-inventory-management/internal/utils"
)

func okHandler(called *bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		*called = true
		w.WriteHeader(http.StatusOK)
	}
}

func TestRequireRole(t *testing.T) {
	tests := []struct {
		name       string
		claims     *utils.Claims // nil = no claims in context
		wantStatus int
		wantCalled bool
	}{
		{
			name:       "admin allowed",
			claims:     &utils.Claims{Role: string(model.RoleAdmin)},
			wantStatus: http.StatusOK,
			wantCalled: true,
		},
		{
			name:       "staff denied",
			claims:     &utils.Claims{Role: string(model.RoleStaff)},
			wantStatus: http.StatusForbidden,
			wantCalled: false,
		},
		{
			name:       "no claims",
			claims:     nil,
			wantStatus: http.StatusUnauthorized,
			wantCalled: false,
		},
	}

	for _, tt := range tests {
		called := false
		handler := RequireRole(model.RoleAdmin)(okHandler(&called))

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		if tt.claims != nil {
			req = req.WithContext(context.WithValue(req.Context(), claimsKey, tt.claims))
		}
		rec := httptest.NewRecorder()

		handler(rec, req)

		if rec.Code != tt.wantStatus {
			t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
		}

		if called != tt.wantCalled {
			t.Errorf("next called = %v, want %v", called, tt.wantCalled)
		}
	}
}

func TestAuthenticate(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-at-least-32-bytes-long!!")
	validToken, _ := utils.GenerateToken(1, string(model.RoleAdmin))

	tests := []struct {
		name       string
		authHeader string
		wantStatus int
	}{
		{
			name:       "missing header",
			authHeader: "",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "no bearer prefix",
			authHeader: validToken,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "malformed token",
			authHeader: "Bearer not-a-jwt",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "valid token",
			authHeader: "Bearer " + validToken,
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			handler := Authenticate(okHandler(&called))

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			rec := httptest.NewRecorder()

			handler(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}
