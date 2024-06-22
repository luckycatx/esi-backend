package bootstrap

import (
	"context"
	"database/sql"
	"esi/internal/pkg/conf"
	"esi/internal/pkg/db"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type App struct {
	Cfg   *conf.Config
	Mysql *sql.DB
	Redis *redis.Client
}

func NewApp() *App {
	var app = &App{}

	app.Cfg = conf.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	app.Mysql = db.InitMysql(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		app.Cfg.Mysql.User, app.Cfg.Mysql.Pwd,
		app.Cfg.Mysql.Host, app.Cfg.Mysql.Port, app.Cfg.Mysql.DBName,
	))

	app.Redis = db.InitRedis(ctx, &redis.Options{
		Addr:     app.Cfg.Redis.Addr,
		Password: app.Cfg.Redis.Pwd,
		DB:       0,
	})

	return app
}
