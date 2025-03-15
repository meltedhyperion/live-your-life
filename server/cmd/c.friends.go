package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/meltedhyperion/globetrotter/server/util"
)

func HandleFriendRoutes(app *App) http.Handler {
	r := chi.NewRouter()
	r.Post("/{friend_id}", app.handleMakeFriend)
	return r
}

func (app *App) handleMakeFriend(w http.ResponseWriter, r *http.Request) {
	playerId, err := GetUserID(r.Context())
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, nil, "User not authenticated")
		return
	}
	friendId := chi.URLParam(r, "friend_id")

	resp, _, err := app.DB.From("players").Select("name", "", false).Eq("id", friendId).Execute()
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()}, "Error in getting player")
		return
	}
	var player []util.Player
	err = json.Unmarshal(resp, &player)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()}, "Error in getting player")
		return
	}
	if len(player) == 0 {
		sendErrorResponse(w, http.StatusNotFound, nil, "Friend not found")
		return
	}

	addFriend := util.AddFriend{
		Player1ID: playerId,
		Player2ID: friendId,
	}

	resp, _, err = app.DB.From("friends").Insert(addFriend, false, "", "", "").Execute()
	if err != nil {
		sendErrorResponse(w, http.StatusConflict, nil, "Friend already exists")
		return
	}

	sendResponse(w, http.StatusOK, resp, "Friend added successfully")

}
