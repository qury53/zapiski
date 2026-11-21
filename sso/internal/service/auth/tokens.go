package auth_service

import (
	"errors"
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	core_errors "github.com/qury53/GC_learn/sso/internal/core/errors"
)

const (
	Issuer           = "sso-zapiski"
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

// Claims данные, которые хранятся в jwt токене, нужны в основном для фронта
type Claims struct {
	UserID    int64  `json:"user_id"`
	TokenType string `json:"typ"`
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
		UserID:    userID,
		TokenType: TokenTypeAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    Issuer,
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(s.jwtSecret)
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
	claims := Claims{
		UserID:    userID,
		TokenType: TokenTypeRefresh,
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
	tokenStr, err := token.SignedString(s.jwtSecret)
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

		return s.jwtSecret, nil
	},
		jwt.WithIssuer(Issuer), // Проверка iss
		jwt.WithValidMethods([]string{"HS256"}),
	)

	if err != nil {
		return nil, fmt.Errorf("%w: %w", core_errors.ErrUnauthorized, err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("%w: invalid token", core_errors.ErrUnauthorized)
	}

	if claims.TokenType != TokenTypeAccess {
		return nil, fmt.Errorf("%w: invalid token type", core_errors.ErrUnauthorized)
	}

	return claims, nil
}

// ValidateRefreshToken валидирует refresh токен.
// Для истёкшего токена возвращает его claims и core_errors.ErrTokenExpired —
// по jti вызывающий код отличает истечение от повторного использования.
//
// Returns: claims, признак валидности, ошибку

func (s *AuthService) ValidateRefreshToken(tokenStr string) (Claims, bool, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	},
		jwt.WithIssuer(Issuer),
		jwt.WithValidMethods([]string{"HS256"}),
	)

	if err != nil {
		// Истёкший токен возвращаем вместе с claims: по jti вызывающий код
		// находит сессию и отличает истечение от повторного использования.
		if errors.Is(err, jwt.ErrTokenExpired) {
			return *claims, false, core_errors.ErrTokenExpired
		}

		return Claims{}, false, fmt.Errorf("%w: %w", core_errors.ErrUnauthorized, err)
	}

	if !token.Valid {
		return Claims{}, false, fmt.Errorf("%w: invalid token", core_errors.ErrUnauthorized)
	}

	if claims.TokenType != TokenTypeRefresh {
		return Claims{}, false, fmt.Errorf("%w: invalid token type", core_errors.ErrUnauthorized)
	}

	return *claims, true, nil
}

// ParseRefreshToken парсит RefreshTokenClaims из refresh токена
func (s *AuthService) ParseRefreshToken(
	tokenString string,
) (Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return s.jwtSecret, nil
	},
		jwt.WithIssuer(Issuer), // Проверка iss
		jwt.WithValidMethods([]string{"HS256"}),
	)

	if err != nil {
		return Claims{}, fmt.Errorf("%w: %w", core_errors.ErrUnauthorized, err)
	}

	if !token.Valid {
		return Claims{}, fmt.Errorf("%w: invalid token", core_errors.ErrUnauthorized)
	}

	if claims.TokenType != TokenTypeRefresh {
		return Claims{}, fmt.Errorf("%w: invalid token type", core_errors.ErrUnauthorized)
	}

	return *claims, nil
}
