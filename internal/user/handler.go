package user

import (
	"context"
	"esi/internal/pkg/conf"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type Servicer interface {
	login(ctx context.Context, req *LoginReq) (*LoginResp, error)
	logout(ctx context.Context, access_token string) error
	register(ctx context.Context, req *RegReq) (*RegResp, error)
	refresh(ctx context.Context, token_str string) (*RefreshResp, error)
	profile(ctx context.Context) ([]*User, error)
}

// Interface check
var _ Servicer = (*Service)(nil)

/* ===== */ /* ===== */ /* ===== */

type Handler struct {
	cfg *conf.Config
	svc Servicer
}

func NewHandler(cfg *conf.Config, s Servicer) *Handler {
	return &Handler{
		cfg: cfg,
		svc: s,
	}
}

func (h *Handler) Login(c echo.Context) error {
	var req = &LoginReq{}
	if err := c.Bind(req); err != nil {
		return err
	}
	if err := c.Validate(req); err != nil {
		return err
	}
	resp, err := h.svc.login(c.Request().Context(), req)
	if err != nil {
		if err.Error() == "user not found" || err.Error() == "invalid password" {
			return echo.NewHTTPError(http.StatusConflict, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, map[string]any{
		"data": resp,
	})
}

/* ----- */

func (h *Handler) Logout(c echo.Context) error {
	var token = strings.Split(c.Request().Header.Get("Authorization"), " ")[1]
	if token == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Empty authorization token")
	}
	return h.svc.logout(c.Request().Context(), token)
}

/* ----- */

func (h *Handler) Register(c echo.Context) error {
	var req = &RegReq{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if err := c.Validate(req); err != nil {
		return err
	}
	resp, err := h.svc.register(c.Request().Context(), req)
	if err != nil {
		if err.Error() == "user exist" {
			return echo.NewHTTPError(http.StatusConflict, "User already exists")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, map[string]any{
		"data": resp,
	})
}

/* ----- */

func (h *Handler) Refresh(c echo.Context) error {
	var req = &RefreshReq{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if err := c.Validate(req); err != nil {
		return err
	}
	resp, err := h.svc.refresh(c.Request().Context(), req.RefreshToken)
	if err != nil {
		if err.Error() == "user exist" {
			return echo.NewHTTPError(http.StatusConflict, "User already exists")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, map[string]any{
		"data": resp,
	})
}

/* ----- */

func (h *Handler) Profile(c echo.Context) error {
	resp, err := h.svc.profile(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, map[string]any{
		"data": resp,
	})
}
