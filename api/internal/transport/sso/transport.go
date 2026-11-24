package sso

import (
	"context"
	"net/http"
	"strings"

	"github.com/qury53/GC_learn/api/internal/core/ports"
	"github.com/qury53/GC_learn/platform/logger"
	"google.golang.org/grpc/metadata"
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthService interface {
	Register(ctx context.Context, req RegisterRequest) (RegisterResponse, error)
	Login(ctx context.Context, req LoginRequest, headers ports.HeaderCarrier) (LoginResponse, error)
	Refresh(ctx context.Context, refreshToken string, headers ports.HeaderCarrier) (TokenPair, error)
	Logout(ctx context.Context, refreshToken string, headers ports.HeaderCarrier) error
	VerifyEmail(ctx context.Context, req VerifyEmailRequest, headers ports.HeaderCarrier) (VerifyEmailResponse, error)
	HandleOAuthCallback(ctx context.Context, req OauthCallbackRequest, headers ports.HeaderCarrier) (OauthCallbackResponse, error)
	GetOAuthConnectUrl(ctx context.Context, req OauthConnectRequest, headers ports.HeaderCarrier) (OauthConnectResponse, error)
}

type AuthHandler struct {
	service AuthService
	log     *logger.Logger
}

func NewAuthHandler(service AuthService, log *logger.Logger) *AuthHandler {
	return &AuthHandler{service: service, log: log}
}

func extractRefreshFromCookie(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	var rawCookies []string
	rawCookies = append(rawCookies, md.Get("cookie")...)

	if len(rawCookies) == 0 {
		return ""
	}

	req := &http.Request{Header: http.Header{}}
	for _, line := range rawCookies {
		req.Header.Add("Cookie", strings.TrimSpace(line))
	}

	if cookie, err := req.Cookie("refresh_token"); err == nil {
		return cookie.Value
	}
	return ""
}

func extractCodeVerifierFromCookie(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	var rawCookies []string
	rawCookies = append(rawCookies, md.Get("cookie")...)

	if len(rawCookies) == 0 {
		return ""
	}

	req := &http.Request{Header: http.Header{}}
	for _, line := range rawCookies {
		req.Header.Add("Cookie", strings.TrimSpace(line))
	}

	if cookie, err := req.Cookie("code_verifier"); err == nil {
		return cookie.Value
	}
	return ""
}
