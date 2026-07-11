package middleware

import (
	"net/http"
	"strings"

	"github.com/Hettank/habit-tracker/internal/response"
	"github.com/Hettank/habit-tracker/internal/utils"
)

const bearerPrefix = "Bearer "

func AuthMiddleware(
	jwtManager *utils.JWTManager,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				response.Unauthorized(
					w,
					"authorization header missing",
				)
				return
			}

			if !strings.HasPrefix(authHeader, bearerPrefix) {
				response.Unauthorized(
					w,
					"invalid authorization header",
				)
				return
			}

			token := strings.TrimSpace(
				strings.TrimPrefix(authHeader, bearerPrefix),
			)

			claims, err := jwtManager.ValidateAccessToken(token)

			if err != nil {
				response.Unauthorized(
					w,
					"invalid or expired token",
				)
				return
			}

			ctx := SetClaims(
				r.Context(),
				claims,
			)

			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
