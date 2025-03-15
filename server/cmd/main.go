package main

import (
	"net/http"

	"github.com/meltedhyperion/globetrotter/server/logger"
)

type App struct {
	Srv *http.Server
}

func main() {
	app := &App{}

	InitConfig()
	InitServer(app)

	logger.Log.Info("api running on", app.Srv.Addr)
	logger.Log.Fatal(app.Srv.ListenAndServe())

}
