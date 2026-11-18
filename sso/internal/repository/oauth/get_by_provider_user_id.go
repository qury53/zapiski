package oauth_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/qury53/GC_learn/sso/internal/core/domain"
	core_errors "github.com/qury53/GC_learn/sso/internal/core/errors"
)

// GetByProviderUserID ищет привязку по провайдеру и его идентификатору пользователя.
//
// Returns: domain.OAuthAccount, error

func (r *OAuthRepository) GetByProviderUserID(
	ctx context.Context,
	provider domain.Provider,
	providerUserID string,
) (domain.OAuthAccount, error) {
	const query = `
		SELECT ` + oauthAccountColumns + `
		FROM zapiski.oauth_accounts
		WHERE provider = $1 AND provider_user_id = $2
	`

	var m OAuthAccountModel
	err := r.pool.QueryRow(ctx, query, string(provider), providerUserID).Scan(
		&m.ID,
		&m.UserID,
		&m.Provider,
		&m.ProviderUserID,
		&m.Email,
		&m.EmailVerified,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.OAuthAccount{}, fmt.Errorf(
				"oauth account not found: %w",
				core_errors.ErrNotFound,
			)
		}

		return domain.OAuthAccount{}, fmt.Errorf("oauth_repository.GetByProviderUserID: %w", err)
	}

	return modelToDomain(m), nil
}
