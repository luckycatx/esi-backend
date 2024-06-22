// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

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

// Injectors from wire.go:

// TODO: Implement cleanup func
func InitServer(app *bootstrap.App) (*echo.Echo, func(), error) {
	config := app.Cfg
	db := app.Mysql
	repo := user.NewRepo(db)
	confToken := config.Token
	client := app.Redis
	redisAdapter := token.NewRedisAdapter(client)
	tokenToken := token.NewToken(confToken, redisAdapter)
	service := user.NewService(config, repo, tokenToken)
	userHandler := user.NewHandler(config, service)
	handlerHandler := handler.NewHandler(userHandler)
	jwtAuth := auth.NewJWTAuth(tokenToken)
	middlewareMiddleware := middleware.NewMiddleware(jwtAuth)
	echoEcho := router.NewRouter(config, handlerHandler, middlewareMiddleware)
	return echoEcho, func() {
	}, nil
}

// wire.go:

var ProviderSet = wire.NewSet(token.NewRedisAdapter, token.NewToken, auth.NewJWTAuth, user.NewRepo, user.NewService, user.NewHandler, handler.NewHandler, middleware.NewMiddleware, router.NewRouter, wire.FieldsOf(new(*bootstrap.App), "Cfg", "Mysql", "Redis"), wire.FieldsOf(new(*conf.Config), "Token"), wire.Bind(new(token.Blacklist), new(*token.RedisAdapter)), wire.Bind(new(auth.TokenUtil), new(*token.Token)), wire.Bind(new(user.TokenUtil), new(*token.Token)), wire.Bind(new(user.Repoer), new(*user.Repo)), wire.Bind(new(user.Servicer), new(*user.Service)), wire.Bind(new(middleware.JWTAuthMW), new(*auth.JWTAuth)), wire.Bind(new(handler.UserHandler), new(*user.Handler)))