package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func initHandler(app *App, r *chi.Mux) {
	r.Get("/", func(rw http.ResponseWriter, r *http.Request) {
		sendResponse(rw, 200, nil, "Welcome to the Globetrotter API!")
	})
	r.With(AuthMiddleware).Mount("/players", HandlePlayerRoutes(app))
	r.Mount("/questions", HandleQuestionRoutes(app))
	r.With(AuthMiddleware).Mount("/friends", HandleFriendRoutes(app))
	r.With(AuthMiddleware).Mount("/sessions", HandleSessionsRoutes(app))
}
