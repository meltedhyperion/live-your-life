package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/meltedhyperion/globetrotter/server/db/pg_db"
	"github.com/meltedhyperion/globetrotter/server/util"
)

func HandleSessionsRoutes(app *App) http.Handler {
	r := chi.NewRouter()
	r.Post("/create", app.handleCreateUserSession)
	r.Get("/{id}", app.handleGetUSerSessionById)
	r.Put("/{id}", app.UpdateUserSessionById)
	return r
}

func (app *App) handleCreateUserSession(w http.ResponseWriter, r *http.Request) {
	playerId, err := GetUserID(r.Context())
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, nil, "User not authenticated")
		return
	}
	question_ids, err := app.store.GetRandomDestinationsForSessionQuestions(r.Context())
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, nil, "Error creating questions")
		return
	}
	user_sesion := &pg_db.CreateUserSessionParams{
		UserID:       uuid.MustParse(playerId),
		Destinations: question_ids,
	}
	err = app.store.CreateUserSession(r.Context(), user_sesion)
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, nil, "Error creating session")
		return
	}
	sendResponse(w, 200, nil, "Session created successfully")
}

func (app *App) handleGetUSerSessionById(w http.ResponseWriter, r *http.Request) {
	playerId, err := GetUserID(r.Context())
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, nil, "User not authenticated")
		return
	}
	id, _ := app.store.GetAllUserSessionByID(r.Context(), uuid.MustParse(playerId))

	session, err := app.store.GetUserSessionByID(r.Context(), &pg_db.GetUserSessionByIDParams{
		UserID: uuid.MustParse(playerId),
		ID:     id,
	})
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, nil, "Session not found")
		return
	}
	sendResponse(w, 200, session, "Session sent successfuly")
}

func (app *App) UpdateUserSessionById(w http.ResponseWriter, r *http.Request) {
	playerId, err := GetUserID(r.Context())
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, nil, "User not authenticated")
		return
	}

	id := chi.URLParam(r, "id")
	body, err := getBodyWithType[util.CheckAnswerRequest](r)
	if err != nil {
		sendHerrorResponse(w, err)
		return
	}

	isCorrect, _, err := util.CheckAnswerToQuestionID(app.store, int32(body.QuestionID), body.Answer)
	if err != nil {
		sendHerrorResponse(w, err)
		return
	}
	idINt, _ := strconv.Atoi(id)

	session, _ := app.store.GetUserSessionByID(r.Context(), &pg_db.GetUserSessionByIDParams{
		ID:     int32(idINt),
		UserID: uuid.MustParse(playerId),
	})
	totalAttempted := session.TotalAttempted

	correct := session.Correct
	if isCorrect {
		correct.Int32++
	}
	totalAttempted.Int32++

	updatedSession := pg_db.UpdateUserSessionParams{
		Score:          float64(correct.Int32),
		TotalAttempted: totalAttempted,
		Correct:        correct,
		ID:             session.ID,
	}
	err = app.store.UpdateUserSession(r.Context(), &updatedSession)
	if err != nil {
		sendHerrorResponse(w, err)
		return
	}

}
