//go:build wireinject
// +build wireinject

package main

import (
	"esi/api/handler"
	"esi/api/middleware"
	"esi/api/router"
	"esi/bootstrap"
	"esi/internal/auth"
	"esi/internal/pkg/conf"
	"esi/internal/pkg/token"
	"esi/internal/user"

	"github.com/google/wire"
	"github.com/labstack/echo/v4"
)

var ProviderSet = wire.NewSet(
	// pkg
	token.NewRedisAdapter, token.NewToken,
	// auth
	auth.NewJWTAuth,
	// user
	user.NewRepo, user.NewService, user.NewHandler,
	// api
	handler.NewHandler, middleware.NewMiddleware, router.NewRouter,

	// Required fields binding
	wire.FieldsOf(new(*bootstrap.App), "Cfg", "Mysql", "Redis"),
	wire.FieldsOf(new(*conf.Config), "Token"),
	// Interface implementation binding
	wire.Bind(new(token.Blacklist), new(*token.RedisAdapter)),
	wire.Bind(new(auth.TokenUtil), new(*token.Token)),
	wire.Bind(new(user.TokenUtil), new(*token.Token)),
	wire.Bind(new(user.Repoer), new(*user.Repo)),
	wire.Bind(new(user.Servicer), new(*user.Service)),
	wire.Bind(new(middleware.JWTAuthMW), new(*auth.JWTAuth)),
	wire.Bind(new(handler.UserHandler), new(*user.Handler)),
)

// TODO: Implement cleanup func
func InitServer(app *bootstrap.App) (*echo.Echo, func(), error) {
	panic(wire.Build(ProviderSet))
}
