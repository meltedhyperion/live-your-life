package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/meltedhyperion/globetrotter/server/util"
)

func HandlePlayerRoutes(app *App) http.Handler {
	r := chi.NewRouter()
	r.Get("/create", app.handleCreatePlayer)
	r.Get("/{player_id}", app.handleGetPlayerById)
	return r
}

func (app *App) handleCreatePlayer(w http.ResponseWriter, r *http.Request) {
	body, err := getBodyWithType[util.CreatePlayerReq](r)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, nil, err.Error())
		return
	}
	avatar := util.GenerateAvatar(body.UserId)

	player := util.Player{
		ID:             body.UserId,
		Name:           body.Name,
		Avatar:         avatar,
		CorrectAnswers: 0,
		TotalAttempts:  0,
		Score:          0.0,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	filterBuilder := app.DB.From("players").Insert(player, false, "", "representation", "")
	resp, _, err := filterBuilder.Execute()
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, nil, "Error in creating player")
		return
	}

	sendResponse(w, http.StatusOK, resp, "Player created successfully")
}

func (app *App) handleGetPlayerById(w http.ResponseWriter, r *http.Request) {
	playerId := chi.URLParam(r, "player_id")

	filterBuilder := app.DB.From("players").Select("*", "exact", false).Eq("id", playerId)
	resp, _, err := filterBuilder.Execute()
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, nil, "Error in getting player")
		return
	}

	sendResponse(w, http.StatusOK, resp, "Player fetched successfully")
}
