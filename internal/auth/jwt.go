package auth

import (
	"context"
	"net/http"
	"strings"

	"esi/internal/pkg/token"

	"github.com/labstack/echo/v4"
)

type TokenUtil interface {
	ParseToken(auth string) (*token.TokenInfo, error)
	IsBlocked(ctx context.Context, token_id string) bool
}

// Interface check
var _ TokenUtil = (*token.Token)(nil)

/* ===== */ /* ===== */ /* ===== */

type JWTAuth struct {
	token TokenUtil
}

func NewJWTAuth(t TokenUtil) *JWTAuth {
	return &JWTAuth{token: t}
}

func (a *JWTAuth) Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var auth = strings.Split(c.Request().Header.Get("Authorization"), " ")[1]
		if auth == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "authorization token is required")
		}
		if a.token.IsBlocked(c.Request().Context(), auth) {
			return echo.NewHTTPError(http.StatusUnauthorized, "token has been blocked")
		}
		token_info, _ := a.token.ParseToken(auth)
		c.Set("user", token_info)
		c.Set("x-user-id", token_info.UID)
		return next(c)
	}
}