package adapters_sso

import (
	"fmt"
	ssov1 "shared/pkg/proto/sso/v1"
	"strings"

	"github.com/qury53/GC_learn/api/internal/core/ports"
	"google.golang.org/grpc/metadata"
)

func copySetCookies(headers ports.HeaderCarrier, mds ...metadata.MD) {
	if headers == nil {
		return
	}
	for _, md := range mds {
		for _, v := range md.Get("set-cookie") {
			headers.Add("Set-Cookie", v)
		}
	}
}

func authStatusFromPB(s ssov1.AuthStatus) (ports.AuthStatus, error) {
	switch s {
	case ssov1.AuthStatus_AUTH_STATUS_ACTIVE:
		return ports.AuthStatusActive, nil
	case ssov1.AuthStatus_AUTH_STATUS_PENDING_VERIFICATION:
		return ports.AuthStatusPendingVerification, nil
	case ssov1.AuthStatus_AUTH_STATUS_BLOCKED:
		return ports.AuthStatusBlocked, nil
	case ssov1.AuthStatus_AUTH_STATUS_UNSPECIFIED:
		return "", fmt.Errorf("auth status is unspecified")
	default:
		return "", fmt.Errorf("unsupported auth status: %v", s)
	}
}

func oauthProviderToPB(provider string) (ssov1.OAuthProvider, error) {
	switch strings.ToLower(strings.TrimSpace(provider)) {
	case "google":
		return ssov1.OAuthProvider_OAUTH_PROVIDER_GOOGLE, nil
	case "github":
		return ssov1.OAuthProvider_OAUTH_PROVIDER_GITHUB, nil
	case "yandex":
		return ssov1.OAuthProvider_OAUTH_PROVIDER_YANDEX, nil
	default:
		return ssov1.OAuthProvider_OAUTH_PROVIDER_UNSPECIFIED, fmt.Errorf("unsupported oauth provider: %s", provider)
	}
}

func codeChallengeMethodToPB(method string) (ssov1.CodeChallengeMethod, error) {
	switch strings.ToUpper(strings.TrimSpace(method)) {
	case "", "S256":
		return ssov1.CodeChallengeMethod_CODE_CHALLENGE_METHOD_S256, nil
	default:
		return ssov1.CodeChallengeMethod_CODE_CHALLENGE_METHOD_UNSPECIFIED, fmt.Errorf("unsupported code challenge method: %s", method)
	}
}
