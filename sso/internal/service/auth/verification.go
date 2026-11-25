package auth_service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	core_errors "github.com/qury53/GC_learn/sso/internal/core/errors"
)

const (
	// cooldownDuration — минимальный интервал между отправками кода на один email
	cooldownDuration = 5 * time.Minute
	// verificationCodeTTL — срок жизни кода. Задаётся явно и не зависит от общего
	// CacheTTL из конфига Redis. Должен быть больше cooldownDuration, иначе запись
	// исчезнет раньше, чем закончится кулдаун.
	verificationCodeTTL = 15 * time.Minute
)

// verificationCacheKey строит ключ кэша по email.
// Email нормализуется, чтобы один и тот же адрес в разном регистре не порождал
// разные записи — иначе кулдаун обходится сменой регистра.
func (s *AuthService) verificationCacheKey(email string) string {
	return s.cacheClient.GenerateCacheKey("verification", "email", normalizeEmail(email))
}

// verificationCooldownKey — ключ кулдауна повторной отправки. Живёт отдельно
// от кода и удаляется вместе с ним: ограничение держит Redis (SET NX PX),
// а не арифметика по локальным часам инстанса.
func (s *AuthService) verificationCooldownKey(email string) string {
	return s.cacheClient.GenerateCacheKey("verificationCooldown", "email", normalizeEmail(email))
}

func (s *AuthService) verificationCodeAttemptsCacheKey(email string) string {
	return s.cacheClient.GenerateCacheKey("verificationAttempts", "email", normalizeEmail(email))
}

// normalizeEmail приводит email к виду, в котором он используется в ключах кэша.
// Один и тот же адрес в разном регистре/с пробелами должен попадать в одну запись.
func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// dropVerificationCode удаляет код, счётчик попыток и кулдаун. Нужен, когда код
// заведомо не дойдёт до пользователя: иначе кулдаун заблокирует повторную отправку.
func (s *AuthService) dropVerificationCode(ctx context.Context, email string) error {
	return errors.Join(
		s.cacheClient.Del(ctx, s.verificationCacheKey(email)),
		s.cacheClient.Del(ctx, s.verificationCodeAttemptsCacheKey(email)),
		s.cacheClient.Del(ctx, s.verificationCooldownKey(email)),
	)
}

// CacheVerificationCode атомарно выдаёт новый код верификации: проверяет кулдаун,
// сохраняет хэш кода и сбрасывает счётчик неудачных попыток.
//
// В Redis кладётся только HMAC-хэш кода, а не сам код: дамп кэша не должен
// давать возможность подтвердить чужой email.
//
// Returns: ошибку ErrTooManyRequest, если кулдаун ещё не истёк
func (s *AuthService) CacheVerificationCode(ctx context.Context, email string, code string) error {
	keys := []string{
		s.verificationCacheKey(email),             // KEYS: code
		s.verificationCodeAttemptsCacheKey(email), // KEYS: attempts
		s.verificationCooldownKey(email),          // KEYS: cooldown
	}

	res, err := issueCodeScript.Run(
		ctx,
		s.cacheClient,
		keys,
		s.encryptor.HashToken(code),
		verificationCodeTTL.Milliseconds(),
		cooldownDuration.Milliseconds(),
	).Int()
	if err != nil {
		return fmt.Errorf("failed to execute issue code script: %w", err)
	}

	if res == 0 {
		return fmt.Errorf(
			"verification code cooldown: %w",
			core_errors.ErrTooManyRequest,
		)
	}

	return nil
}
