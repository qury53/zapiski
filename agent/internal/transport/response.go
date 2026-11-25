package agent_transport

import (
	"context"
	"errors"

	core_errors "github.com/qury53/GC_learn/agent/internal/errors"
	"google.golang.org/grpc/codes"
)

func RespondCode(err error) codes.Code {
	var code codes.Code
	switch {
	case errors.Is(err, core_errors.ErrInvalidArgument):
		code = codes.InvalidArgument
	case errors.Is(err, core_errors.ErrNotFound):
		code = codes.NotFound
	case errors.Is(err, core_errors.ErrTooMantRequest):
		code = codes.ResourceExhausted
	case errors.Is(err, core_errors.InsufficientBalance):
		code = codes.ResourceExhausted
	case errors.Is(err, core_errors.ErrUnauthorized):
		code = codes.Unauthenticated
	case errors.Is(err, core_errors.ErrUnprocessableEntity):
		code = codes.InvalidArgument
	case errors.Is(err, core_errors.ErrUndefinedError):
		code = codes.Unknown
	case errors.Is(err, core_errors.ErrInternalServerError):
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
