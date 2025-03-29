package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/meltedhyperion/globetrotter/server/db/pg_db"
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
	friendId, err := uuid.Parse(chi.URLParam(r, "friend_id"))
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, nil, "Friend Not Found")
		return
	}

	friend, err := app.store.GetPlayerById(context.Background(), friendId)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, nil, "Error Finding Friend")
		return
	}
	if friend == nil {
		sendErrorResponse(w, http.StatusNotFound, nil, "Friend not found")
		return
	}

	err = app.store.AddFriend(context.Background(), &pg_db.AddFriendParams{
		Player1ID: uuid.MustParse(playerId),
		Player2ID: friendId,
	})
	if err != nil {
		sendErrorResponse(w, http.StatusConflict, nil, "Friend already exists")
		return
	}
	sendResponse(w, http.StatusOK, nil, "Friend added successfully")

}
