package user_transport

import (
	userv1 "shared/pkg/proto/user/v1"
	"time"

	"github.com/qury53/GC_learn/user/internal/core/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// UserProfileDomainToProto конвертирует доменную модель в proto-сообщение.
func UserProfileDomainToProto(u *domain.UserProfile) *userv1.UserProfile {
	if u == nil {
		return nil
	}
	return &userv1.UserProfile{
		Id:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		State:     verificationStatusDomainToProto(u.State),
		CreatedAt: timeDomainToProto(u.CreatedAt),
	}
}

// UserProfilesDomainToProto конвертирует срез доменных моделей в GetUsersResponse.
func UserProfilesDomainToProto(users []domain.UserProfile) *userv1.GetUsersResponse {
	if len(users) == 0 {
		return nil
	}
	resp := &userv1.GetUsersResponse{
		Users: make([]*userv1.UserProfile, len(users)),
	}
	for i := range users {
		resp.Users[i] = UserProfileDomainToProto(&users[i])
	}
	return resp
}

// verificationStatusDomainToProto маппит доменный статус в proto-enum.
func verificationStatusDomainToProto(s domain.VerificationStatus) userv1.VerificationStatus {
	switch s {
	case domain.StatusActive:
		return userv1.VerificationStatus_VERIFICATION_STATUS_ACTIVE
	case domain.StatusPendingVerification:
		return userv1.VerificationStatus_VERIFICATION_STATUS_PENDING_VERIFICATION
	case domain.StatusBlocked:
		return userv1.VerificationStatus_VERIFICATION_STATUS_BLOCKED
	case domain.StatusUnknown:
		fallthrough
	default:
		return userv1.VerificationStatus_VERIFICATION_STATUS_UNSPECIFIED
	}
}

// timeDomainToProto безопасно конвертирует time.Time в *timestamppb.Timestamp.
// Нулевое время возвращается как nil, чтобы не сериализовать 1970-01-01.
func timeDomainToProto(t time.Time) *timestamppb.Timestamp {
	if t.IsZero() {
		return nil
	}
	return timestamppb.New(t)
}
