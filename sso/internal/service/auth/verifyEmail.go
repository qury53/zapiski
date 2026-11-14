package auth_service

import (
	"context"

	"github.com/qury53/GC_learn/sso/internal/core/domain"
)

func (s *AuthService) VerifyEmail(
	ctx context.Context,
	email string,
	code string,
) domain.VerificationStatus {

}
