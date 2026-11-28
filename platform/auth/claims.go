package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

const (
	// Issuer — значение iss в токенах, которые выпускает sso.
	Issuer = "sso-zapiski"
	// TokenTypeAccess — значение typ для access-токена.
	TokenTypeAccess = "access"
)

// Claims полезная нагрузка access-токена sso.
type Claims struct {
	UserID    int64  `json:"user_id"`
	TokenType string `json:"typ"`
	jwt.RegisteredClaims
}

// Validator проверяет подпись и обязательные поля access-токена.
type Validator struct {
	secret []byte
	issuer string
}

func NewValidator(secret, issuer string) *Validator {
	if issuer == "" {
		issuer = Issuer
	}

	return &Validator{
		secret: []byte(secret),
		issuer: issuer,
	}
}

// ValidateAccess проверяет access-токен и возвращает claims.
func (v *Validator) ValidateAccess(raw string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(raw, claims,
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}

			return v.secret, nil
		},
		jwt.WithIssuer(v.issuer),
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	if claims.TokenType != TokenTypeAccess {
		return nil, fmt.Errorf("%w: wrong token type", ErrInvalidToken)
	}

	if claims.UserID <= 0 {
		return nil, fmt.Errorf("%w: empty user id", ErrInvalidToken)
	}

	return claims, nil
}
