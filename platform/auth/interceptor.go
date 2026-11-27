package auth

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Policy описывает, какие gRPC-методы доступны без токена.
// Всё, что не помечено как Public, требует валидный access-токен.
type Policy struct {
	public map[string]struct{}
}

func NewPolicy() *Policy {
	return &Policy{public: make(map[string]struct{})}
}

// Public помечает методы как доступные без токена.
func (p *Policy) Public(methods ...string) *Policy {
	for _, method := range methods {
		p.public[method] = struct{}{}
	}

	return p
}

func (p *Policy) IsPublic(method string) bool {
	_, ok := p.public[method]

	return ok
}

// UnaryServerInterceptor проверяет access-токен и кладёт claims в контекст.
func UnaryServerInterceptor(validator *Validator, policy *Policy) grpc.UnaryServerInterceptor {
	if policy == nil {
		policy = NewPolicy()
	}

	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		if policy.IsPublic(info.FullMethod) {
			return handler(ctx, req)
		}

		raw, ok := bearerFromMetadata(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing access token")
		}

		claims, err := validator.ValidateAccess(raw)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid access token")
		}

		return handler(WithClaims(ctx, claims), req)
	}
}

func bearerFromMetadata(ctx context.Context) (string, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", false
	}

	values := md.Get("authorization")
	if len(values) == 0 {
		return "", false
	}

	raw := strings.TrimSpace(values[0])
	raw = strings.TrimSpace(strings.TrimPrefix(raw, "Bearer "))
	if raw == "" {
		return "", false
	}

	return raw, true
}
