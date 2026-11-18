// Package core_oauth описывает контракт с внешним OAuth 2.0 / OIDC провайдером.
// Конкретная реализация (Google) живёт в repository слое, а сервис зависит
// только от этого интерфейса и легко подменяется в тестах.
package core_oauth

import "context"

// Identity — нормализованный профиль пользователя от провайдера.
// ProviderUserID — стабильный идентификатор (у Google это sub),
// email у провайдера может меняться.
type Identity struct {
	ProviderUserID string
	Email          string
	EmailVerified  bool
	Name           string
}

// Client — минимальный набор операций, который нужен сервису.
type Client interface {
	// ConnectURL собирает URL страницы согласия с привязанным state и PKCE challenge.
	ConnectURL(state string, codeChallenge string) string

	// Exchange меняет authorization code на профиль пользователя.
	// codeVerifier обязателен: PKCE поддерживается только в режиме S256.
	Exchange(ctx context.Context, code string, codeVerifier string) (Identity, error)
}
