package router

import (
	"esi/api/handler"
	"esi/api/middleware"
	"esi/internal/pkg/conf"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	mw "github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/acme/autocert"
)

func NewRouter(cfg *conf.Config, h *handler.Handler, m *middleware.Middleware) *echo.Echo {
	var e = echo.New()
	e.HTTPErrorHandler = HTTPErrorHandler
	e.AutoTLSManager.Cache = autocert.DirCache("/var/www/.cache")
	e.Use(mw.LoggerWithConfig(mw.LoggerConfig{
		Format:           "[ECHO] ${time_custom} | ${status} | ${latency_human} | ${remote_ip} | ${method} ${uri} | ${error}\n\n",
		CustomTimeFormat: "2006/01/02 - 15:04:05",
		Output:           e.Logger.Output(),
	}))
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
	auth.Any("/files", echo.WrapHandler(http.StripPrefix("/files", h.Upload)))
	auth.Any("/files/*", echo.WrapHandler(http.StripPrefix("/files/", h.Upload)))

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

func HTTPErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}
	var code int
	var msg any
	var ierr error
	if he, ok := err.(*echo.HTTPError); ok {
		if he.Internal != nil {
			if ihe, ok := he.Internal.(*echo.HTTPError); ok {
				he = ihe // Override if internal error is an HTTP error
			} else {
				ierr = he.Internal // Else add to error list
			}
		}
		code, msg = he.Code, he.Message
		switch m := msg.(type) {
		case error:
			msg = m.Error()
		}
	} else {
		if c.Echo().Debug {
			// Only response the message in debug mode if the error returned is non-HTTPError
			code, msg = http.StatusInternalServerError, err.Error()
		} else {
			code, msg = http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)
		}
	}

	// Only response internal error in debug mode (in errors list)
	if ierr != nil && c.Echo().Debug {
		msg = map[string]any{
			"error": map[string]any{
				"code":    code,
				"message": "Internal error occurred",
				"errors": []map[string]any{
					{"message": msg}, {"message": ierr.Error()},
				},
			},
		}
	} else {
		msg = map[string]any{
			"error": map[string]any{
				"code":    code,
				"message": msg,
			},
		}
	}

	if c.Request().Method == http.MethodHead {
		err = c.NoContent(code)
	} else {
		err = c.JSON(code, msg)
	}
	if err != nil {
		c.Logger().Error(err)
	}
}
