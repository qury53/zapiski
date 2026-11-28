package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func signToken(t *testing.T, secret string, claims Claims) string {
	t.Helper()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	raw, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}

	return raw
}

func TestValidateAccess(t *testing.T) {
	const secret = "test-secret"

	validator := NewValidator(secret, Issuer)

	registered := jwt.RegisteredClaims{
		Issuer:    Issuer,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	tests := []struct {
		name    string
		raw     string
		wantErr bool
	}{
		{
			name: "valid access token",
			raw: signToken(t, secret, Claims{
				UserID:           42,
				TokenType:        TokenTypeAccess,
				RegisteredClaims: registered,
			}),
		},
		{
			name: "wrong secret",
			raw: signToken(t, "other-secret", Claims{
				UserID:           42,
				TokenType:        TokenTypeAccess,
				RegisteredClaims: registered,
			}),
			wantErr: true,
		},
		{
			name: "refresh token instead of access",
			raw: signToken(t, secret, Claims{
				UserID:           42,
				TokenType:        "refresh",
				RegisteredClaims: registered,
			}),
			wantErr: true,
		},
		{
			name: "expired token",
			raw: signToken(t, secret, Claims{
				UserID:    42,
				TokenType: TokenTypeAccess,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    Issuer,
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Minute)),
				},
			}),
			wantErr: true,
		},
		{
			name: "empty user id",
			raw: signToken(t, secret, Claims{
				TokenType:        TokenTypeAccess,
				RegisteredClaims: registered,
			}),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := validator.ValidateAccess(tt.raw)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if claims.UserID != 42 {
				t.Fatalf("user id = %d, want 42", claims.UserID)
			}
		})
	}
}
