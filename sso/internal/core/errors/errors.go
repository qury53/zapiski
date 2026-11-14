package core_errors

import "errors"

var (
	ErrConflict           = errors.New("conflict")
	ErrInvalidArgument    = errors.New("invalid argument")
	ErrNotFound           = errors.New("not found")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrTooManyRequest     = errors.New("too many requests")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrForbidden          = errors.New("forbidden")

	ErrSessionRevoked = errors.New("session revoked")
	ErrSessionExpired = errors.New("session expired")

	ErrTokenExpired = errors.New("token expired")
)
