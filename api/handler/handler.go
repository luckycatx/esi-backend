package handler

import (
	"esi/internal/user"

	"github.com/labstack/echo/v4"
)

type UserHandler interface {
	Login(c echo.Context) error
	Logout(c echo.Context) error
	Register(c echo.Context) error
	Refresh(c echo.Context) error
	Profile(c echo.Context) error
}

// Interface check
var _ UserHandler = (*user.Handler)(nil)

/* ===== */ /* ===== */ /* ===== */

type Handler struct {
	User UserHandler
}

func NewHandler(u UserHandler) *Handler {
	return &Handler{User: u}
}
