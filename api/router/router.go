package router

import (
	"esi/api/handler"
	"esi/api/middleware"
	"esi/internal/pkg/conf"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	mw "github.com/labstack/echo/v4/middleware"
)

func NewRouter(cfg *conf.Config, h *handler.Handler, m *middleware.Middleware) *echo.Echo {
	var e = echo.New()
	e.Use(mw.Logger())
	e.Use(mw.Recover())
	e.Use(mw.CORS())
	e.Validator = &Validator{v: validator.New()}

	var pub = e.Group("")
	pub.POST("/login", h.User.Login)
	pub.POST("/register", h.User.Register)
	pub.POST("/refresh", h.User.Refresh)
	// pub.GET("/gopy", util.GoPy)

	var auth = e.Group("", m.JWTAuth)
	auth.GET("/user/profile", h.User.Profile)
	auth.POST("/user/logout", h.User.Logout)

	return e
}

type Validator struct {
	v *validator.Validate
}

func (v *Validator) Validate(s any) error {
	if err := v.v.Struct(s); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
