package ports

import "context"

type HeaderCarrier interface {
	Add(key, value string)
	Values(key string) []string
}

type AuthStatus string

const (
	AuthStatusActive              AuthStatus = "AUTH_STATUS_ACTIVE"
	AuthStatusPendingVerification AuthStatus = "AUTH_STATUS_PENDING_VERIFICATION"
	AuthStatusBlocked             AuthStatus = "AUTH_STATUS_BLOCKED"
)

type RegisterRequest struct {
	Username string
	Password string
	Email    string
}

type RegisterResponse struct {
	UserID int64
	Status AuthStatus
}

type LoginRequest struct {
	Username string
	Password string
}

type LoginResponse struct {
	AccessToken string
	UserID      int64
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type VerifyEmailRequest struct {
	Email string
	Code  string
}

type VerifyEmailResponse struct {
	Status AuthStatus
}

type OauthConnectRequest struct {
	Provider            string
	RedirectURI         string
	CodeChallenge       string
	CodeChallengeMethod string
}

type OauthConnectResponse struct {
	AuthorizationURL string
	OAuthState       string
}

type OauthCallbackRequest struct {
	Provider     string
	Code         string
	State        string
	CodeVerifier string
}

type OauthCallbackResponse struct {
	AccessToken string
	UserID      int64
}

type SSOClient interface {
	Register(ctx context.Context, req RegisterRequest) (RegisterResponse, error)
	Login(ctx context.Context, req LoginRequest, headers HeaderCarrier) (LoginResponse, error)
	Refresh(ctx context.Context, refreshToken string, headers HeaderCarrier) (TokenPair, error)
	Logout(ctx context.Context, refreshToken string, headers HeaderCarrier) error
	VerifyEmail(ctx context.Context, req VerifyEmailRequest, headers HeaderCarrier) (VerifyEmailResponse, error)
	HandleOAuthCallback(ctx context.Context, req OauthCallbackRequest, headers HeaderCarrier) (OauthCallbackResponse, error)
	GetOAuthConnectUrl(ctx context.Context, req OauthConnectRequest, headers HeaderCarrier) (OauthConnectResponse, error)
}
