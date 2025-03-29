package main

import (
	"context"
	"net/http"
	"sort"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/meltedhyperion/globetrotter/server/db/pg_db"
	"github.com/meltedhyperion/globetrotter/server/util"
)

func HandlePlayerRoutes(app *App) http.Handler {
	r := chi.NewRouter()
	r.Post("/create", app.handleCreatePlayer)
	r.Get("/", app.handleGetPlayerById)
	r.Get("/leaderboard", app.handleGetLeaderboard)
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
		sendHerrorResponse(w, err)
		return
	}
	avatar := util.GenerateAvatar(playerId)

	player := &pg_db.CreateNewPlayerParams{
		ID:     uuid.MustParse(playerId),
		Name:   body.Name,
		Avatar: avatar,
	}
	err = app.store.CreateNewPlayer(context.Background(), player)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			sendErrorResponse(w, http.StatusConflict, nil, "Player already exists")
			return
		}
		sendErrorResponse(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()}, "Error in creating player")
		return
	}

	sendResponse(w, http.StatusOK, nil, "Player created successfully")
}

func (app *App) handleGetPlayerById(w http.ResponseWriter, r *http.Request) {
	playerId, err := GetUserID(r.Context())
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, nil, "User not authenticated")
		return
	}
	player, err := app.store.GetPlayerById(context.Background(), uuid.MustParse(playerId))
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err, "Error in getting player data")
		return
	}
	if player == nil {
		sendErrorResponse(w, http.StatusNotFound, nil, "Player Data Not Found")
		return
	}
	sendResponse(w, http.StatusOK, player, "Player fetched successfully")
}

func (app *App) handleGetLeaderboard(w http.ResponseWriter, r *http.Request) {
	playerId, err := GetUserID(r.Context())
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, nil, "User not authenticated")
		return
	}
	leaderboard, err := app.store.GetLeaderboardForFriends(context.Background(), uuid.MustParse(playerId))
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err, "Error in getting player data")
		return
	}
	if leaderboard == nil || len(leaderboard) == 1 {
		sendErrorResponse(w, http.StatusNotFound, nil, "Sorry buddy you are friendless")
		return
	}
	sort.Slice(leaderboard, func(i, j int) bool {
		return leaderboard[i].Score > leaderboard[j].Score
	})
	playerStats := util.Leaderboard{
		PlayerStats: leaderboard,
	}
	sendResponse(w, http.StatusOK, playerStats, "Leaderboard fetched successfully")
}
