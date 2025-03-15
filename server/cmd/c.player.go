package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"time"

	"github.com/go-chi/chi/v5"
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
	res, _, err := app.DB.From("players").Select("id", "exact", false).Eq("id", playerId).Execute()
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()}, "Error in fetching player")
		return
	}
	var p []map[string]interface{}
	err = json.Unmarshal(res, &p)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()}, "Error in getting player")
		return
	}
	if len(p) > 0 {
		sendErrorResponse(w, http.StatusConflict, nil, "Player already exists")
		return
	}
	body, err := getBodyWithType[util.CreatePlayerReq](r)
	if err != nil {
		sendHerrorResponse(w, err)
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

	_, _, err = app.DB.From("players").Insert(player, false, "", "representation", "").Execute()
	if err != nil {
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

func (app *App) handleGetLeaderboard(w http.ResponseWriter, r *http.Request) {
	playerId, err := GetUserID(r.Context())
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, nil, "User not authenticated")
		return
	}
	var leaderboard util.Leaderboard
	var listOfFriendIDs []util.FriendsIDs
	var leaderboardPlayers []string
	resp, _, err := app.DB.From("friends").Select("player2_id", "exact", false).Eq("player1_id", playerId).Execute()
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()}, "Error in getting friends")
		return
	}
	err = json.Unmarshal(resp, &listOfFriendIDs)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()}, "Error in getting friends")
		return
	}
	if len(listOfFriendIDs) == 0 {
		sendErrorResponse(w, http.StatusNotFound, nil, "Sorry buddy you are friendless")
		return
	}
	for _, friend := range listOfFriendIDs {
		leaderboardPlayers = append(leaderboardPlayers, friend.FriendID)
	}
	leaderboardPlayers = append(leaderboardPlayers, playerId)
	resp, _, err = app.DB.From("players").Select("name, avatar, correct_answers, total_attempts, score", "exact", false).In("id", leaderboardPlayers).Execute()
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()}, "Error in getting leaderboard")
		return
	}
	err = json.Unmarshal(resp, &leaderboard.PlayerStats)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()}, "Error in getting leaderboard")
		return
	}
	sort.Slice(leaderboard.PlayerStats, func(i, j int) bool {
		return leaderboard.PlayerStats[i].Score > leaderboard.PlayerStats[j].Score
	})
	sendResponse(w, http.StatusOK, leaderboard, "Leaderboard fetched successfully")
}
