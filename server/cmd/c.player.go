package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/meltedhyperion/globetrotter/server/util"
)

func HandlePlayerRoutes(app *App) http.Handler {
	r := chi.NewRouter()
	r.Get("/create", app.handleCreatePlayer)
	r.Get("/", app.handleGetPlayerById)
	return r
}

func (app *App) handleCreatePlayer(w http.ResponseWriter, r *http.Request) {
	playerId, err := GetUserID(r.Context())
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, nil, "User not authenticated")
		return
	}
	body, err := getBodyWithType[util.CreatePlayerReq](r)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, nil, err.Error())
		return
	}
	avatar := util.GenerateAvatar(playerId)

	player := util.Player{
		ID:             playerId,
		Name:           body.Name,
		Avatar:         avatar,
		CorrectAnswers: 0,
		TotalAttempts:  0,
		Score:          0.0,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	resp, _, err := app.DB.From("players").Insert(player, false, "", "representation", "").Execute()
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()}, "Error in creating player")
		return
	}

	sendResponse(w, http.StatusOK, resp, "Player created successfully")
}

func (app *App) handleGetPlayerById(w http.ResponseWriter, r *http.Request) {
	playerId, err := GetUserID(r.Context())
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, nil, "User not authenticated")
		return
	}
	var player []util.Player
	resp, _, err := app.DB.From("players").Select("id, name, avatar, correct_answers, total_attempts, score", "exact", false).Eq("id", playerId).Execute()
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()}, "Error in getting player")
		return
	}
	err = json.Unmarshal(resp, &player)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()}, "Error in getting player")
		return
	}
	if len(player) == 0 {
		sendErrorResponse(w, http.StatusNotFound, nil, "Player not found")
		return
	}
	sendResponse(w, http.StatusOK, player[0], "Player fetched successfully")
}
