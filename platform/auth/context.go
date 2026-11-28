package auth

import "context"

type contextKey struct{}

// WithClaims кладёт проверенные claims в контекст.
func WithClaims(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, contextKey{}, claims)
}

// ClaimsFrom достаёт claims из контекста.
func ClaimsFrom(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(contextKey{}).(*Claims)

	return claims, ok
}

// UserIDFrom возвращает id пользователя из контекста.
func UserIDFrom(ctx context.Context) (int64, bool) {
	claims, ok := ClaimsFrom(ctx)
	if !ok || claims.UserID <= 0 {
		return 0, false
	}

	return claims.UserID, true
}
