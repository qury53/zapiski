package user_repository

import (
	"time"

	"github.com/qury53/GC_learn/user/internal/core/domain"
)

const userColumns = `
	id,
	username,
	email,
	state,
	created_at
`

type UserProfileModel struct {
	ID        int64     `db:"id"`
	Username  string    `db:"username"`
	Email     string    `db:"email"`
	State     string    `db:"state"`
	CreatedAt time.Time `db:"created_at"`
}

func modelToDomain(model UserProfileModel) domain.UserProfile {
	return domain.UserProfile{
		ID:        model.ID,
		Username:  model.Username,
		Email:     model.Email,
		State:     mapDBVerificationStatusToDomain(model.State),
		CreatedAt: model.CreatedAt,
	}
}

func modelsToDomain(models []UserProfileModel) []domain.UserProfile {
	domains := make([]domain.UserProfile, len(models))
	for i, model := range models {
		domains[i] = modelToDomain(model)
	}

	return domains
}

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
