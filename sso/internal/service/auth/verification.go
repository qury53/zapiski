package auth_service

import (
	"context"
	"fmt"
	"strings"
	"time"

	core_errors "github.com/qury53/GC_learn/sso/internal/core/errors"
)

// VerificationCacheItem структура кода верификации, записывается в кэш,
// хранит код и дату создания
type VerificationCacheItem struct {
	Code      string    `json:"code"`
	CreatedAt time.Time `json:"created_at"`
}

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
	normalized := strings.ToLower(strings.TrimSpace(email))

	return s.cacheClient.GenerateCacheKey("verification", "email", normalized)
}

func (s *AuthService) verificationCodeAttemptsCacheKey(email string) string {
	normalized := strings.ToLower(strings.TrimSpace(email))
	return s.cacheClient.GenerateCacheKey("verificationAttempts", "email", normalized)
}

// dropVerificationCode удаляет код из кэша. Нужен, когда код заведомо
// не дойдёт до пользователя: иначе кулдаун заблокирует повторную отправку.
func (s *AuthService) dropVerificationCode(ctx context.Context, email string) error {
	return s.cacheClient.Del(ctx, s.verificationCacheKey(email))
}

// CacheVerificationCode проверяет кулдаун и сохраняет сгенерированный код в кэш
func (s *AuthService) CacheVerificationCode(ctx context.Context, email string, code string) error {
	cacheKey := s.verificationCacheKey(email)

	var existingItem VerificationCacheItem
	found, err := s.cacheClient.Get(ctx, cacheKey, &existingItem)
	if err != nil {
		return fmt.Errorf("failed to check cache: %w", err)
	}

	now := time.Now()

	if found && !existingItem.CreatedAt.IsZero() {
		// Отрицательную разницу тоже считаем кулдауном: расхождение часов
		// не должно позволять обойти ограничение.
		if timePassed := now.Sub(existingItem.CreatedAt); timePassed < cooldownDuration {
			timeLeft := cooldownDuration - timePassed

			return fmt.Errorf(
				"verification code cooldown, time left:%v %w",
				timeLeft,
				core_errors.ErrTooManyRequest,
			)
		}
	}

	newItem := VerificationCacheItem{
		Code:      code,
		CreatedAt: now,
	}

	// TTL задаём явно: обычный Set берёт общий CacheTTL из конфига Redis,
	// а он может быть не задан (0) — тогда код не истечёт никогда.
	if err := s.cacheClient.SetWithTTL(ctx, cacheKey, newItem, verificationCodeTTL); err != nil {
		return fmt.Errorf("failed to save code to cache: %w", err)
	}

	return nil
}
