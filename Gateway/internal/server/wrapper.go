package server

import "context"

// в контексте лучше хранить значение по ключу не примитивного типа, на это линтер ругается
type (
	bodyKey   struct{}
	claimsKey struct{}
)

// Это для более удобной работы, вместо ctx.Value("body").(T)
func WithBody[T any](ctx context.Context, body T) context.Context {
	return context.WithValue(ctx, bodyKey{}, body)
}

func GetBody[T any](ctx context.Context) T {
	return ctx.Value(bodyKey{}).(T)
}

func WithClaims(ctx context.Context, claims AuthClaims) context.Context {
	return context.WithValue(ctx, claimsKey{}, claims)
}

func GetClaims(ctx context.Context) (AuthClaims, bool) {
	claims, ok := ctx.Value(claimsKey{}).(AuthClaims)
	return claims, ok
}
