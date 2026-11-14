package domain

import (
	"errors"
	"fmt"
	"net/netip"
	"time"
)

// Provider определяет тип OAuth провайдера
type Provider string

const (
	ProviderGoogle Provider = "google"
	ProviderGithub Provider = "github"
	ProviderYandex Provider = "yandex"
)

const (
	// revokedReasonRotated — причина отзыва, которую проставляет Rotate у старой сессии.
	RevokedReasonRotated = "rotated"
	// revokedReasonRefreshReuse — причина отзыва при компрометации всей ветки.
	RevokedReasonRefreshReuse = "refresh_reuse"
	// revokedReasonLogout - причина отзыва при выходе из аккаунта
	RevokedReasonLogout = "logout"
)

// CodeChallengeMethod определяет метод хеширования PKCE
type CodeChallengeMethod string

const (
	PKCEMethodS256 CodeChallengeMethod = "S256"
)

var ErrUnsupportedProvider = errors.New("unsupported oauth provider")

// Session - сессия пользователя
type Session struct {
	ID               int64
	UserID           int64
	JTI              string
	RefreshTokenHash string
	UserAgent        *string
	IPAddress        *netip.Addr
	ExpiresAt        time.Time
	RevokedAt        *time.Time
	RevokedReason    *string
	ParentSessionID  *int64
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// RotatelParams - структура, которая нужня для атомарной замены старой сессии на новую
type RotateParams struct {
	NewJTI              string
	NewRefreshTokenHash string
	NewExpiresAt        time.Time
	UserAgent           *string
	IPAddress           *netip.Addr
}

func NewSessionUninitialized(
	userID int64,
	jti string,
	refreshTokenHash string,
	userAgent *string,
	IPAddress *netip.Addr,
	expiresAt time.Time,
	parentSessionID *int64,
) Session {
	return Session{
		UserID:           userID,
		JTI:              jti,
		RefreshTokenHash: refreshTokenHash,
		UserAgent:        userAgent,
		IPAddress:        IPAddress,
		ExpiresAt:        expiresAt.UTC(),
		ParentSessionID:  parentSessionID,
	}
}

func (s *Session) Validate() error {
	if s.UserID <= 0 {
		return fmt.Errorf("session: user_id must be positive")
	}

	if s.JTI == "" {
		return fmt.Errorf("session: jti is required")
	}

	if s.RefreshTokenHash == "" {
		return fmt.Errorf("session: refresh_token_hash is required")
	}

	if len(s.RefreshTokenHash) > 255 {
		return fmt.Errorf("session: refresh_token_hash exceeds 255 chars")
	}

	if s.ExpiresAt.IsZero() {
		return fmt.Errorf("session: expires_at is required")
	}

	if !s.ExpiresAt.After(time.Now()) {
		return fmt.Errorf("session: expires_at must be in the future")
	}

	if s.ParentSessionID != nil && *s.ParentSessionID <= 0 {
		return fmt.Errorf("session: parent_session_id must be positive")
	}

	return nil
}
