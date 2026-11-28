package domain

import "time"

type AuthResponse struct {
	ExpiredAt   time.Time
	AccessToken string
}
