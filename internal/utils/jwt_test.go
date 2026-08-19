package utils

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var testSecret = strings.Repeat("a", 32)

func TestTokenRoundTrip(t *testing.T) {

	t.Setenv("JWT_SECRET", testSecret)

	tokenString, err := GenerateToken(42, "ADMIN")
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	claims, err := ParseToken(tokenString)
	if err != nil {
		t.Fatalf("ParseToken() error = %v", err)
	}

	if claims.UserId != 42 {
		t.Errorf("UserId = %d, want 42", claims.UserId)
	}
	if claims.Role != "ADMIN" {
		t.Errorf("Role = %q, want %q", claims.Role, "ADMIN")
	}
}

func TestParseToken_Expired(t *testing.T) {

	t.Setenv("JWT_SECRET", testSecret)

	claims := Claims{
		UserId: 1,
		Role:   "ADMIN",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}

	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(testSecret))
	if err != nil {
		t.Fatalf("signing setup failed: %v", err)
	}

	if _, err := ParseToken(tokenString); err == nil {
		t.Fatalf("expected error for expired token, got nil")
	}
}

func TestParseToken_RejectsNoneAlgorithm(t *testing.T) {

	t.Setenv("JWT_SECRET", testSecret)

	claims := Claims{
		UserId: 1,
		Role:   "ADMIN",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}

	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodNone, claims).SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		t.Fatalf("signing setup failed: %v", err)
	}

	if _, err := ParseToken(tokenString); err == nil {
		t.Fatalf("expected error for alg=none token, got nil")
	}
}

func TestCheckJWTSecret(t *testing.T) {
	tests := []struct {
		name    string
		secret  string
		wantErr bool
	}{
		{
			name:    "empty",
			secret:  "",
			wantErr: true,
		},
		{
			name:    "too short",
			secret:  strings.Repeat("a", 10),
			wantErr: true,
		},
		{
			name:    "valid length",
			secret:  strings.Repeat("a", 32),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("JWT_SECRET", tt.secret)
			err := CheckJWTSecret()
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckJWTSecret() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseToken_WrongSecret(t *testing.T) {
	t.Setenv("JWT_SECRET", testSecret)
	tokenString, err := GenerateToken(1, "ADMIN")
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	t.Setenv("JWT_SECRET", strings.Repeat("b", 32))
	if _, err := ParseToken(tokenString); err == nil {
		t.Fatal("expected error for token signed with a different secret, got nil")
	}
}
