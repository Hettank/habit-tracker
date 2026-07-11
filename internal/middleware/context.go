package middleware

import (
	"context"

	"github.com/Hettank/habit-tracker/internal/utils"
)

type contextKey string

const claimsContextKey contextKey = "claims"

func SetClaims(
	ctx context.Context,
	claims *utils.Claims,
) context.Context {
	return context.WithValue(
		ctx,
		claimsContextKey,
		claims,
	)
}

func GetClaims(
	ctx context.Context,
) (*utils.Claims, bool) {

	claims, ok := ctx.Value(
		claimsContextKey,
	).(*utils.Claims)

	return claims, ok
}
