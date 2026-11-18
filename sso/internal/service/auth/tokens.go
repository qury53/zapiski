package auth_service

import (
	"errors"
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	core_errors "github.com/qury53/GC_learn/sso/internal/core/errors"
)

const (
	Issuer           = "zapiski"
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

// Claims данные, которые хранятся в jwt токене, нужны в основном для фронта
type Claims struct {
	UserID int64
	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	UserID int64
	jwt.RegisteredClaims
}

// GetRefreshTokenTTL возвращает срок жизни refresh токена
//
// Returns: time.Duration

func (s *AuthService) GetRefreshTokenTTL() time.Duration {
	return s.RefreshTokenTTL
}

// GenerateAccessToken генерирует access токен
//
// Returns: string, error

func (s *AuthService) GenerateAccessToken(userID int64) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    Issuer,
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(s.JWTSecret)
	if err != nil {
		return "", fmt.Errorf(
			"failed to create access token: %w",
			err,
		)
	}

	return tokenStr, nil
}

// GenerateRefreshToken генерирует refresh токен
//
// Returns: string, error

func (s *AuthService) GenerateRefreshToken(userID int64, jti string) (string, error) {
	claims := RefreshTokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.RefreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    Issuer,
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(s.JWTSecret)
	if err != nil {
		return "", fmt.Errorf(
			"failed to generate refresh token: %w",
			err,
		)
	}

	return tokenStr, nil
}

// ValidateToken валидирует токен и возвращает Claims
//
// Returns: *Claims, error

func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return s.JWTSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// ValidateRefreshToken валидирует refresh токен.
// Для истёкшего токена возвращает его claims и core_errors.ErrTokenExpired —
// по jti вызывающий код отличает истечение от повторного использования.
//
// Returns: claims, признак валидности, ошибку

func (s *AuthService) ValidateRefreshToken(tokenStr string) (RefreshTokenClaims, bool, error) {
	claims := &RefreshTokenClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.JWTSecret, nil
	})

	if err != nil {
		// Истёкший токен возвращаем вместе с claims: по jti вызывающий код
		// находит сессию и отличает истечение от повторного использования.
		if errors.Is(err, jwt.ErrTokenExpired) {
			return *claims, false, core_errors.ErrTokenExpired
		}

		return RefreshTokenClaims{}, false, fmt.Errorf("failed to parse jwt: %w", err)
	}

	if claims, ok := token.Claims.(*RefreshTokenClaims); ok && token.Valid {
		return *claims, true, nil
	}

	return RefreshTokenClaims{}, false, nil
}

// ParseRefreshToken парсит RefreshTokenClaims из refresh токена
func (s *AuthService) ParseRefreshToken(
	tokenString string,
) (RefreshTokenClaims, error) {
	claims := RefreshTokenClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return s.JWTSecret, nil
	})

	if err != nil {
		return RefreshTokenClaims{}, err
	}

	if claims, ok := token.Claims.(RefreshTokenClaims); ok && token.Valid {
		return claims, nil
	}

	return RefreshTokenClaims{}, fmt.Errorf("invalid token: %w %v", core_errors.ErrUnauthorized, err)
}
