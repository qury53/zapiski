package domain

import "time"

type VerificationStatus int32

const (
	StatusUnknown             VerificationStatus = iota // 0 — обязательный нулевой вариант
	StatusActive                                        // 1
	StatusPendingVerification                           // 2
	StatusBlocked                                       // 3
)

type UserProfile struct {
	ID        int64
	Username  string
	Email     string
	State     VerificationStatus
	CreatedAt time.Time
}
