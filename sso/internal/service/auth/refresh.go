package auth_service

import (
	"context"
	"crypto/subtle"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/qury53/GC_learn/sso/internal/core/domain"
	core_errors "github.com/qury53/GC_learn/sso/internal/core/errors"
)

// ErrRefreshTokenReuse — предъявлен refresh токен из уже ротированной сессии.
// Оборачивает ErrUnauthorized, чтобы транспорт отвечал клиенту 401.
var ErrRefreshTokenReuse = fmt.Errorf("%w: refresh token reuse", core_errors.ErrUnauthorized)

// Refresh обменивает refresh токен на новую пару access/refresh:
// старая сессия отзывается, новая создаётся дочерней (ротация).
//
// Returns: access токен, refresh токен, ошибку
func (s *AuthService) Refresh(
	ctx context.Context,
	refreshToken string,
) (string, string, error) {
	now := time.Now()

	claims, valid, err := s.ValidateRefreshToken(refreshToken)
	tokenExpired := errors.Is(err, core_errors.ErrTokenExpired)

	// Истёкший токен не отбрасываем сразу: его jti нужен, чтобы найти сессию
	// и отличить обычное истечение от повторного использования (см. п. 5).
	if err != nil && !tokenExpired {
		s.log.Warn("refresh: invalid token", "error", err)

		return "", "", fmt.Errorf("%w: invalid refresh token: %v", core_errors.ErrUnauthorized, err)
	}

	if !valid && !tokenExpired {
		return "", "", fmt.Errorf("%w: token is invalid", core_errors.ErrUnauthorized)
	}

	if claims.ID == "" {
		return "", "", fmt.Errorf("%w: token without jti", core_errors.ErrUnauthorized)
	}

	session, err := s.sessionRepository.GetByJTIForUpdate(ctx, claims.ID)
	if err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			return "", "", fmt.Errorf("%w: session not found", core_errors.ErrUnauthorized)
		}

		return "", "", fmt.Errorf("service: get session for update: %w", err)
	}

	// Страховка от подмены пары (userID из одного токена, jti из другого).
	if session.UserID != claims.UserID {
		s.log.Warn("refresh: session user mismatch",
			"jti", claims.ID,
			"token_user_id", claims.UserID,
			"session_user_id", session.UserID,
		)

		return "", "", fmt.Errorf("%w: user mismatch", core_errors.ErrUnauthorized)
	}

	// 3. Отзыв / ротация / срок / хэш токена.
	if stateErr := s.checkSessionState(session, refreshToken, now); stateErr != nil {
		if errors.Is(stateErr, ErrRefreshTokenReuse) {
			// 4. Компрометация: помечаем и рвём ветку.
			s.handleRefreshReuse(ctx, session)
		}

		return "", "", stateErr
	}

	// 5. Проверяем exp только ПОСЛЕ проверки состояния — иначе истёкший
	//    ротированный токен уйдёт как «просто истёк» и reuse не поймаем.
	if tokenExpired {
		return "", "", fmt.Errorf("%w: %w", core_errors.ErrUnauthorized, core_errors.ErrTokenExpired)
	}

	newJTI := uuid.NewString()

	newRefreshToken, err := s.GenerateRefreshToken(session.UserID, newJTI)
	if err != nil {
		return "", "", fmt.Errorf("service: generate refresh token: %w", err)
	}

	rotated, err := s.sessionRepository.Rotate(ctx, session.JTI, domain.RotateParams{
		NewJTI:              newJTI,
		NewRefreshTokenHash: s.encryptor.HashToken(newRefreshToken),
		NewExpiresAt:        now.Add(s.RefreshTokenTTL),
		// Переносим контекст устройства, иначе у новой сессии он потеряется.
		UserAgent: session.UserAgent,
		IPAddress: session.IPAddress,
	})
	if err != nil {
		// Сессию могли отозвать между чтением и ротацией — для клиента это не 500.
		if errors.Is(err, core_errors.ErrSessionRevoked) ||
			errors.Is(err, core_errors.ErrSessionExpired) {
			return "", "", fmt.Errorf("%w: %w", core_errors.ErrUnauthorized, err)
		}

		return "", "", fmt.Errorf("service: rotate session: %w", err)
	}

	// 6. Access выпускаем после успешной ротации: если Rotate упал,
	//    клиент не получит рабочую пару, а сессия останется прежней.
	accessToken, err := s.GenerateAccessToken(session.UserID)
	if err != nil {
		return "", "", fmt.Errorf("service: generate access token: %w", err)
	}

	s.log.Debug("refresh: session rotated",
		"user_id", session.UserID,
		"old_jti", session.JTI,
		"new_jti", newJTI,
		"new_session_id", rotated.ID,
		"ip", session.IPAddress,
	)

	return accessToken, newRefreshToken, nil
}

// checkSessionState проверяет состояние сессии: отзыв/ротацию, срок жизни и хэш токена.
// Порядок важен: reuse проверяем раньше срока, иначе истёкший ротированный токен
// выглядит как «просто истёк» и подмену не поймать.
//
// Returns: ErrRefreshTokenReuse при повторном использовании ротированного токена
func (s *AuthService) checkSessionState(
	session domain.Session,
	refreshToken string,
	now time.Time,
) error {
	if session.RevokedAt != nil {
		// Ротированный токен пришёл повторно — значит, им уже воспользовались.
		if session.RevokedReason != nil && *session.RevokedReason == domain.RevokedReasonRotated {
			return fmt.Errorf("%w: session %s already rotated", ErrRefreshTokenReuse, session.JTI)
		}

		return fmt.Errorf("%w: %w", core_errors.ErrUnauthorized, core_errors.ErrSessionRevoked)
	}

	if !session.ExpiresAt.After(now) {
		return fmt.Errorf("%w: %w: expired at %s",
			core_errors.ErrUnauthorized,
			core_errors.ErrSessionExpired,
			session.ExpiresAt.Format(time.RFC3339),
		)
	}

	// Хэш должен совпадать: иначе токен не от этой сессии.
	if subtle.ConstantTimeCompare(
		[]byte(s.encryptor.HashToken(refreshToken)),
		[]byte(session.RefreshTokenHash),
	) != 1 {
		return fmt.Errorf("%w: refresh token hash mismatch", core_errors.ErrUnauthorized)
	}

	return nil
}

// handleRefreshReuse помечает скомпрометированной всю ветку сессий.
// Ошибку не возвращаем: клиенту в любом случае отдаём отказ.
func (s *AuthService) handleRefreshReuse(ctx context.Context, session domain.Session) {
	s.log.Warn("refresh: token reuse detected, revoking session chain",
		"jti", session.JTI,
		"user_id", session.UserID,
		"ip", session.IPAddress,
	)

	if err := s.sessionRepository.MarkCompromised(
		ctx,
		session.JTI,
		domain.RevokedReasonRefreshReuse,
	); err != nil {
		s.log.Error("refresh: mark session chain compromised",
			"jti", session.JTI,
			"user_id", session.UserID,
			"error", err,
		)
	}
}
