package middleware

import (
	"esi/internal/auth"

	"github.com/labstack/echo/v4"
)

type JWTAuthMW interface {
	Auth(next echo.HandlerFunc) echo.HandlerFunc
}

var _ JWTAuthMW = (*auth.JWTAuth)(nil)

/* ===== */ /* ===== */ /* ===== */

type Middleware struct {
	JWTAuth echo.MiddlewareFunc
}

func NewMiddleware(a JWTAuthMW) *Middleware {
	return &Middleware{
		JWTAuth: a.Auth,
	}
}
