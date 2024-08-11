package main

import (
	"esi/bootstrap"
	"log"
)

func main() {
	var app = bootstrap.NewApp()
	if app == nil {
		log.Fatal("Boot failed")
	}
	e, cleanup, err := InitServer(app)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	e.Logger.Fatal(e.StartAutoTLS(app.Cfg.Server.Host + ":" + app.Cfg.Server.Port))
}
