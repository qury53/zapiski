package transport_grpc_response

import (
	"errors"

	core_errors "github.com/qury53/GC_learn/user/internal/core/errors"
	"google.golang.org/grpc/codes"
)

func RespondCode(err error) codes.Code {
	var code codes.Code
	switch {
	case errors.Is(err, core_errors.ErrInvalidArgument):
		code = codes.InvalidArgument
	case errors.Is(err, core_errors.ErrNotFound):
		code = codes.NotFound
	default:
		code = codes.Internal
	}

	return code
}
