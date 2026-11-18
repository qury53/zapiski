package core_transport_grpc_response

import (
	"errors"

	core_errors "github.com/qury53/GC_learn/sso/internal/core/errors"
	"google.golang.org/grpc/codes"
)

func RespondCode(err error) codes.Code {
	var code codes.Code
	switch {
	case errors.Is(err, core_errors.ErrInvalidArgument):
		code = codes.InvalidArgument
	case errors.Is(err, core_errors.ErrNotFound):
		code = codes.NotFound
	case errors.Is(err, core_errors.ErrTooManyRequest):
		code = codes.ResourceExhausted
	case errors.Is(err, core_errors.ErrUnauthorized):
		code = codes.Unauthenticated
	case errors.Is(err, core_errors.ErrForbidden):
		code = codes.PermissionDenied
	case errors.Is(err, core_errors.ErrInvalidCredentials):
		code = codes.Unauthenticated
	case errors.Is(err, core_errors.ErrUnsupportedProvider):
		code = codes.InvalidArgument
	default:
		code = codes.Internal
	}

	return code
}
