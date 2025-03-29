package main

import (
	"net/http"

	"github.com/meltedhyperion/globetrotter/server/db/pg_db"
	"github.com/meltedhyperion/globetrotter/server/logger"
)

type App struct {
	Srv   *http.Server
	store *pg_db.Store
}

func main() {
	app := &App{}

	InitConfig()
	InitDB(app)
	InitServer(app)

	logger.Log.Info("api running on", app.Srv.Addr)
	logger.Log.Fatal(app.Srv.ListenAndServe())

}
