package auth_service

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"net/mail"
	"net/netip"
	"net/url"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
	"github.com/qury53/GC_learn/sso/internal/core/domain"
	core_errors "github.com/qury53/GC_learn/sso/internal/core/errors"
	"golang.org/x/oauth2"
)

// oauthStatePrefix — префикс ключа состояния OAuth-флоу в кэше.
const oauthStatePrefix = "oauth:state"

// oauthState — данные, привязанные к конкретной попытке авторизации.
type oauthState struct {
	Provider    domain.Provider `json:"provider"`
	RedirectURI string          `json:"redirect_uri"`
	Challenge   string          `json:"challenge"`
	CreatedAt   time.Time       `json:"created_at"`
}

// GetConnectURI создает одноразовый state и возвращает URL страницы согласия.
// redirectURI проверяется точным совпадением с allowlist, а не префиксом,
// иначе сервис превращается в open redirect.
//
// Returns: authorization URL, state, error
func (s *AuthService) GetConnectURI(
	ctx context.Context,
	provider domain.Provider,
	redirectURI string,
	challenge string,
) (string, string, error) {
	client, ok := s.oauthClients[provider]
	if !ok {
		return "", "", fmt.Errorf("%w: %s", core_errors.ErrUnsupportedProvider, provider)
	}

	if err := validateRedirectURI(redirectURI); err != nil {
		return "", "", err
	}

	if _, ok := s.allowedRedirects[redirectURI]; !ok {
		return "", "", fmt.Errorf("%w: redirect uri is not allowed", core_errors.ErrInvalidArgument)
	}

	// S256 challenge — base64url от sha256, ровно 43 символа.
	if len(challenge) != 43 || !isBase64URL(challenge) {
		return "", "", fmt.Errorf("%w: invalid S256 code challenge", core_errors.ErrInvalidArgument)
	}

	rawState := make([]byte, 32)
	if _, err := rand.Read(rawState); err != nil {
		return "", "", fmt.Errorf("oauth: generate state: %w", err)
	}

	state := base64.RawURLEncoding.EncodeToString(rawState)
	item := oauthState{
		Provider:    provider,
		RedirectURI: redirectURI,
		Challenge:   challenge,
		CreatedAt:   time.Now(),
	}

	key := s.cacheClient.GenerateCacheKey(oauthStatePrefix, state)
	if err := s.cacheClient.SetWithTTL(ctx, key, item, s.oauthStateTTL); err != nil {
		return "", "", fmt.Errorf("oauth: save state: %w", err)
	}

	return client.ConnectURL(state, challenge), state, nil
}

// ExchangeCallback одноразово потребляет state, проверяет PKCE, обменивает code
// на профиль провайдера и выпускает собственные токены сервиса.
//
// Returns: access token, refresh token, user ID, error
func (s *AuthService) ExchangeCallback(
	ctx context.Context,
	provider domain.Provider,
	code string,
	state string,
	codeVerifier string,
	userAgent *string,
	ipAddr *netip.Addr,
) (string, string, int64, error) {
	client, ok := s.oauthClients[provider]
	if !ok {
		return "", "", 0, fmt.Errorf("%w: %s", core_errors.ErrUnsupportedProvider, provider)
	}

	if code == "" || state == "" || len(codeVerifier) < 43 || len(codeVerifier) > 128 {
		return "", "", 0, fmt.Errorf("%w: invalid callback parameters", core_errors.ErrInvalidArgument)
	}

	key := s.cacheClient.GenerateCacheKey(oauthStatePrefix, state)

	var stateData oauthState
	found, err := s.cacheClient.GetDel(ctx, key, &stateData)
	if err != nil {
		return "", "", 0, fmt.Errorf("oauth: consume state: %w", err)
	}

	// State неизвестен, просрочен, уже использован или принадлежит другому
	// провайдеру — наружу отдаём одну и ту же ошибку.
	if !found || stateData.Provider != provider {
		return "", "", 0, fmt.Errorf("%w: invalid oauth state", core_errors.ErrUnauthorized)
	}

	if subtle.ConstantTimeCompare(
		[]byte(oauth2.S256ChallengeFromVerifier(codeVerifier)),
		[]byte(stateData.Challenge),
	) != 1 {
		return "", "", 0, fmt.Errorf("%w: PKCE verification failed", core_errors.ErrUnauthorized)
	}

	identity, err := client.Exchange(ctx, code, codeVerifier)
	if err != nil {
		return "", "", 0, fmt.Errorf("oauth: exchange callback: %w", err)
	}

	// Автопривязка к существующему аккаунту допустима только для подтвержденного email,
	// иначе любой сможет занять чужой аккаунт, создав такой же адрес у провайдера.
	if identity.ProviderUserID == "" || identity.Email == "" || !identity.EmailVerified {
		return "", "", 0, fmt.Errorf("%w: verified provider identity required", core_errors.ErrUnauthorized)
	}

	identity.Email = strings.ToLower(strings.TrimSpace(identity.Email))
	if _, err := mail.ParseAddress(identity.Email); err != nil {
		return "", "", 0, fmt.Errorf("%w: invalid provider email", core_errors.ErrUnauthorized)
	}

	account, err := s.oauthRepository.GetByProviderUserID(ctx, provider, identity.ProviderUserID)
	if err != nil && !errors.Is(err, core_errors.ErrNotFound) {
		return "", "", 0, fmt.Errorf("oauth: find account: %w", err)
	}

	var user domain.User

	switch {
	case err == nil:
		// Привязка уже есть: идентифицируем по user_id, а не по изменяемому email.
		user, err = s.userRepository.GetUserByID(ctx, account.UserID)
		if err != nil {
			return "", "", 0, fmt.Errorf("oauth: get linked user: %w", err)
		}

		if !user.CanAuthorize() {
			return "", "", 0, fmt.Errorf("%w: user cannot authorize", core_errors.ErrForbidden)
		}

		if err := s.oauthRepository.Touch(ctx, account.ID, identity.Email, identity.EmailVerified); err != nil {
			return "", "", 0, fmt.Errorf("oauth: update account profile: %w", err)
		}

	default:
		user, err = s.userRepository.GetUserByEmail(ctx, identity.Email)
		if err != nil && !errors.Is(err, core_errors.ErrNotFound) {
			return "", "", 0, fmt.Errorf("oauth: find user by email: %w", err)
		}

		if err == nil {
			if !user.CanAuthorize() {
				return "", "", 0, fmt.Errorf("%w: user cannot authorize", core_errors.ErrForbidden)
			}
		} else {
			user, err = s.userRepository.CreateOAuthUser(
				ctx,
				oauthUsername(identity.Email, identity.Name),
				identity.Email,
			)
			if err != nil {
				// Параллельный callback мог создать этого же пользователя.
				if errors.Is(err, core_errors.ErrConflict) {
					user, err = s.userRepository.GetUserByEmail(ctx, identity.Email)
				}

				if err != nil {
					return "", "", 0, fmt.Errorf("oauth: create user: %w", err)
				}
			}
		}

		_, err = s.oauthRepository.Create(ctx, domain.OAuthAccount{
			UserID:         user.ID,
			Provider:       provider,
			ProviderUserID: identity.ProviderUserID,
			Email:          identity.Email,
			EmailVerified:  identity.EmailVerified,
		})
		if err != nil {
			if !errors.Is(err, core_errors.ErrConflict) {
				return "", "", 0, fmt.Errorf("oauth: link account: %w", err)
			}

			// Повторная привязка того же provider+sub к тому же пользователю допустима.
			linked, findErr := s.oauthRepository.GetByProviderUserID(ctx, provider, identity.ProviderUserID)
			if findErr != nil || linked.UserID != user.ID {
				return "", "", 0, fmt.Errorf("%w: oauth account already linked", core_errors.ErrConflict)
			}
		}
	}

	return s.issueSession(ctx, user.ID, userAgent, ipAddr)
}

// issueSession создает refresh-сессию и выпускает пару локальных токенов.
// Используется и в Login, чтобы логика выдачи токенов не разъезжалась.
//
// Returns: access token, refresh token, user ID, error
func (s *AuthService) issueSession(
	ctx context.Context,
	userID int64,
	userAgent *string,
	ipAddr *netip.Addr,
) (string, string, int64, error) {
	jti := uuid.NewString()
	now := time.Now()

	refreshToken, err := s.GenerateRefreshToken(userID, jti)
	if err != nil {
		return "", "", 0, fmt.Errorf("service: generate refresh token: %w", err)
	}

	session := domain.NewSessionUninitialized(
		userID,
		jti,
		s.encryptor.HashToken(refreshToken),
		userAgent,
		ipAddr,
		now.Add(s.RefreshTokenTTL),
		nil,
	)

	if _, err := s.sessionRepository.Create(ctx, session); err != nil {
		return "", "", 0, fmt.Errorf("service: create session: %w", err)
	}

	accessToken, err := s.GenerateAccessToken(userID)
	if err != nil {
		return "", "", 0, fmt.Errorf("service: generate access token: %w", err)
	}

	return accessToken, refreshToken, userID, nil
}

// validateRedirectURI проверяет, что redirect_uri — абсолютный https URL
// без userinfo и фрагмента. Локальный http разрешён только для localhost.
//
// Returns: error
func validateRedirectURI(raw string) error {
	parsed, err := url.ParseRequestURI(raw)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" ||
		parsed.User != nil || parsed.Fragment != "" {
		return fmt.Errorf("%w: invalid redirect uri", core_errors.ErrInvalidArgument)
	}

	if parsed.Scheme != "https" &&
		parsed.Hostname() != "localhost" &&
		parsed.Hostname() != "127.0.0.1" {
		return fmt.Errorf("%w: redirect uri must use https", core_errors.ErrInvalidArgument)
	}

	return nil
}

// isBase64URL проверяет, что строка состоит только из символов base64url без паддинга.
//
// Returns: bool
func isBase64URL(value string) bool {
	if value == "" {
		return false
	}

	for _, r := range value {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-' && r != '_' {
			return false
		}
	}

	return true
}

// oauthUsername собирает username из локальной части email, добивая суффиксом,
// чтобы уложиться в CHECK(char_length BETWEEN 3 AND 32) и не конфликтовать по UNIQUE.
//
// Returns: string
func oauthUsername(email string, name string) string {
	local := strings.SplitN(email, "@", 2)[0]
	if local == "" {
		local = name
	}

	var builder strings.Builder
	for _, r := range strings.ToLower(local) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			builder.WriteRune(r)
		}
	}

	base := strings.Trim(builder.String(), "_")
	if len([]rune(base)) < 3 {
		base = "user"
	}

	if len([]rune(base)) > 22 {
		base = string([]rune(base)[:22])
	}

	return base + "_" + uuid.NewString()[:8]
}
