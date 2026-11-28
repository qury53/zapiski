package core_errors

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
)

// RespondCode маппит доменные ошибки в gRPC-коды.
func RespondCode(err error) codes.Code {
	var code codes.Code

	switch {
	case errors.Is(err, ErrInvalidArgument):
		code = codes.InvalidArgument
	case errors.Is(err, ErrNotFound):
		code = codes.NotFound
	case errors.Is(err, ErrTooMantRequest):
		code = codes.ResourceExhausted
	case errors.Is(err, InsufficientBalance):
		code = codes.ResourceExhausted
	case errors.Is(err, ErrUnauthorized):
		code = codes.Unauthenticated
	case errors.Is(err, ErrUnprocessableEntity):
		code = codes.InvalidArgument
	case errors.Is(err, ErrUndefinedError):
		code = codes.Unknown
	case errors.Is(err, ErrInternalServerError):
		code = codes.Internal
	case errors.Is(err, context.DeadlineExceeded):
		code = codes.DeadlineExceeded
	case errors.Is(err, context.Canceled):
		code = codes.Canceled
	default:
		code = codes.Internal
	}

	return code
}
