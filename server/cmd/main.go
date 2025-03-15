package main

import (
	"net/http"

	"github.com/meltedhyperion/globetrotter/server/logger"
	"github.com/supabase-community/supabase-go"
)

type App struct {
	Srv *http.Server
	DB  *supabase.Client
}

func main() {
	app := &App{}

	InitConfig()
	InitDB(app)
	InitServer(app)

	logger.Log.Info("api running on", app.Srv.Addr)
	logger.Log.Fatal(app.Srv.ListenAndServe())

}
