package oauth_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/qury53/GC_learn/sso/internal/core/domain"
	core_errors "github.com/qury53/GC_learn/sso/internal/core/errors"
)

// Create создает привязку внешнего аккаунта к локальному пользователю.
//
// Returns: domain.OAuthAccount, error

func (r *OAuthRepository) Create(
	ctx context.Context,
	account domain.OAuthAccount,
) (domain.OAuthAccount, error) {
	const query = `
		INSERT INTO zapiski.oauth_accounts (
			user_id, provider, provider_user_id, email, email_verified
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING ` + oauthAccountColumns

	var m OAuthAccountModel
	err := r.pool.QueryRow(
		ctx,
		query,
		account.UserID,
		string(account.Provider),
		account.ProviderUserID,
		nullableString(account.Email),
		account.EmailVerified,
	).Scan(
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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.OAuthAccount{}, fmt.Errorf(
				"oauth account already linked: %w",
				core_errors.ErrConflict,
			)
		}

		return domain.OAuthAccount{}, fmt.Errorf("oauth_repository.Create: %w", err)
	}

	return modelToDomain(m), nil
}

// nullableString превращает пустую строку в NULL
func nullableString(value string) *string {
	if value == "" {
		return nil
	}

	return &value
}
