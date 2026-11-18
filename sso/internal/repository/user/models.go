package users_postgres_repository

import (
	"time"

	"github.com/qury53/GC_learn/sso/internal/core/domain"
)

// UserModel модель пользователя для базы данных
type UserModel struct {
	ID        int64      `db:"id"`
	Version   int64      `db:"version"`
	Username  string     `db:"username"`
	Password  *string    `db:"password"`
	Email     string     `db:"email"`
	State     string     `db:"state"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

// modelToDomain конвертер модели для базы в доменную сущность
//
// Returns: domain.User

func modelToDomain(dto UserModel) domain.User {
	var password string
	if dto.Password != nil {
		password = *dto.Password
	}

	return domain.User{
		ID:       dto.ID,
		Version:  dto.Version,
		Username: dto.Username,
		Password: password,
		Email:    dto.Email,
		State: domain.VerificationStatus(
			mapDBVerificationStatusToDomain(dto.State),
		),
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
	}
}

// mapDBVerificationStatusToDomain мапит статус верификации из бд в доменный статус
//
// Returns: domain.VerificationStatus

func mapDBVerificationStatusToDomain(dbStatus string) domain.VerificationStatus {
	var verification_status domain.VerificationStatus = domain.StatusUnknown
	switch dbStatus {
	case "active":
		verification_status = domain.StatusActive
	case "pending_verification":
		verification_status = domain.StatusPendingVerification
	case "blocked":
		verification_status = domain.StatusBlocked
	}

	return verification_status
}
