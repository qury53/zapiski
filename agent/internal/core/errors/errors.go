package core_errors

import "errors"

var (
	ErrInvalidArgument     = errors.New("invalid argument")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrNotFound            = errors.New("not found")
	ErrUnprocessableEntity = errors.New("unprocessable entity")
	ErrTooMantRequest      = errors.New("too many request")
	ErrUndefinedError      = errors.New("error undefined")

	InsufficientBalance    = errors.New("insufficient balance")
	ErrInternalServerError = errors.New("internal server error")
)
