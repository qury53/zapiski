package oauth_google_repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	core_oauth "github.com/qury53/GC_learn/sso/internal/core/oauth"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// userInfoURL — OIDC userinfo endpoint Google. Профиль забираем server-to-server
// по access token, полученному при обмене code.
const userInfoURL = "https://openidconnect.googleapis.com/v1/userinfo"

// userInfo — ответ userinfo Google.
type userInfo struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
}

// Client реализует core_oauth.Client поверх golang.org/x/oauth2.
type Client struct {
	config     *oauth2.Config
	httpClient *http.Client
}

// NewClient конструктор клиента Google.
//
// Returns: *Client
func NewClient(
	config *Config,
) *Client {
	return &Client{
		config: &oauth2.Config{
			ClientID:     config.ClientID,
			ClientSecret: config.ClientSecret,
			RedirectURL:  config.RedirectURI,
			Endpoint:     google.Endpoint,
			Scopes:       []string{"openid", "email", "profile"},
		},
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// ConnectURL собирает URL страницы согласия. challenge приходит уже посчитанным
// на клиенте: verifier остаётся у клиента и вернётся в callback.
//
// Returns: string
func (c *Client) ConnectURL(state string, codeChallenge string) string {
	return c.config.AuthCodeURL(
		state,
		oauth2.AccessTypeOnline,
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
	)
}

// Exchange меняет code на профиль пользователя.
//
// Returns: core_oauth.Identity, error
func (c *Client) Exchange(
	ctx context.Context,
	code string,
	codeVerifier string,
) (core_oauth.Identity, error) {
	ctx = context.WithValue(ctx, oauth2.HTTPClient, c.httpClient)

	token, err := c.config.Exchange(ctx, code, oauth2.VerifierOption(codeVerifier))
	if err != nil {
		return core_oauth.Identity{}, fmt.Errorf("google: exchange code: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, userInfoURL, nil)
	if err != nil {
		return core_oauth.Identity{}, fmt.Errorf("google: build userinfo request: %w", err)
	}

	token.SetAuthHeader(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return core_oauth.Identity{}, fmt.Errorf("google: userinfo request: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		return core_oauth.Identity{}, fmt.Errorf("google: userinfo status %d", resp.StatusCode)
	}

	var info userInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return core_oauth.Identity{}, fmt.Errorf("google: decode userinfo: %w", err)
	}

	if info.Sub == "" {
		return core_oauth.Identity{}, fmt.Errorf("google: userinfo without sub")
	}

	return core_oauth.Identity{
		ProviderUserID: info.Sub,
		Email:          info.Email,
		EmailVerified:  info.EmailVerified,
		Name:           info.Name,
	}, nil
}
