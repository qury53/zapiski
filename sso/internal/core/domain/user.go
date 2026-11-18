package domain

import (
	"fmt"
	"net/mail"
	"time"
	"unicode/utf8"

	core_errors "github.com/qury53/GC_learn/sso/internal/core/errors"
)

type VerificationStatus int32

const (
	StatusUnknown             VerificationStatus = iota // 0 — обязательный нулевой вариант
	StatusActive                                        // 1
	StatusPendingVerification                           // 2
	StatusBlocked                                       // 3
)

// доменная сущность пользователя
type User struct {
	ID        int64
	Version   int64
	Username  string
	Password  string
	Email     string
	State     VerificationStatus
	CreatedAt time.Time
	UpdatedAt *time.Time
}

// Validate валидация на уровне транспорта, синтакс и тд
// Username between 3 an d 32,
// Password between 6 and 72,
// email should be valid RFC 5322 format
// Returns: error
func (u *User) Validate() error {
	usernameLen := utf8.RuneCountInString(u.Username)
	if usernameLen < 3 || usernameLen > 32 {
		return fmt.Errorf("username length should be between 3 and 32: %w", core_errors.ErrInvalidArgument)
	}

	// Пароль может отсутствовать: OAuth-пользователи создаются без него.
	if u.Password != "" {
		passwordLen := utf8.RuneCountInString(u.Password)
		if passwordLen < 6 || passwordLen > 72 {
			return fmt.Errorf("password length should be between 6 and 72: %w", core_errors.ErrInvalidArgument)
		}
	}

	_, err := mail.ParseAddress(u.Email)
	if err != nil {
		return fmt.Errorf("incorrect email format: %w", core_errors.ErrInvalidArgument)
	}

	return nil
}

// CanAuthorize проверяет, имеет ли право пользователь войти в систему
func (u *User) CanAuthorize() bool {
	return u.State == StatusActive
}

// IsPending проверяет, ожидает ли пользователь подтверждения SMTP
func (u *User) IsPending() bool {
	return u.State == StatusPendingVerification
}

// NewUserUninitialized создает  пользователя, у которого
// поля ID, Version, UpdatedAt не заданы
//
// Returns: пользователя
func NewUserUninitialized(
	username string,
	password string,
	email string,
) User {
	return User{
		ID:       UserIDUnititialized,
		Version:  UserVersionUnititialized,
		Username: username,
		Password: password,
		Email:    email,
	}
}
