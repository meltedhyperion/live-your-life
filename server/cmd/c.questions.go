package main

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/meltedhyperion/globetrotter/server/db/pg_db"
	"github.com/meltedhyperion/globetrotter/server/util"
)

func HandleQuestionRoutes(app *App) http.Handler {
	r := chi.NewRouter()
	r.Get("/", app.handleGetQuestions)
	r.With(AuthMiddleware).Post("/check", app.handleCheckAnswer)
	r.Post("/check/guest", app.handleCheckAnswerForGuest)
	return r
}

func (app *App) handleGetQuestions(w http.ResponseWriter, r *http.Request) {
	destinations, err := app.store.GetRandomDestinationsForQuestions(context.Background())
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err, "Error in getting questions")
		return
	}

	if len(destinations) < 5 {
		http.Error(w, "Not enough destinations available", http.StatusInternalServerError)
		sendErrorResponse(w, http.StatusInternalServerError, nil, "Not enough destinations available")
		return
	}

	var excludeIDs []int32
	for _, d := range destinations {
		excludeIDs = append(excludeIDs, d.ID)
	}
	nameOptions, err := app.store.GetRandomDestinations(context.Background(), excludeIDs)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err, "Error in getting questions")
		return
	}
	if len(nameOptions) < 3 {
		sendErrorResponse(w, http.StatusInternalServerError, nil, "Not enough name options available")
		return
	}

	questions := util.GenerateQuestion(destinations, nameOptions)

	sendResponse(w, http.StatusOK, questions, "Questions fetched successfully")
}

func (app *App) handleCheckAnswer(w http.ResponseWriter, r *http.Request) {
	playerID, err := GetUserID(r.Context())
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, nil, "User not authenticated")
		return
	}

	body, err := getBodyWithType[util.CheckAnswerRequest](r)
	if err != nil {
		sendHerrorResponse(w, err)
		return
	}

	isCorrect, destination, err := util.CheckAnswerToQuestionID(app.store, int32(body.QuestionID), body.Answer)
	if err != nil {
		sendHerrorResponse(w, err)
		return
	}

	player, err := app.store.GetPlayerById(context.Background(), uuid.MustParse(playerID))
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, nil, "Failed to fetch player record")
		return
	}
	if player == nil {
		sendErrorResponse(w, http.StatusNotFound, nil, "Player not found")
		return
	}
	player.TotalAttempts++
	if isCorrect {
		player.CorrectAnswers++
	}
	player.Score = util.CalculateWilsonScore(player.CorrectAnswers, player.TotalAttempts)
	updateScore := &pg_db.UpdatePlayerScoreParams{
		CorrectAnswers: player.CorrectAnswers,
		TotalAttempts:  player.TotalAttempts,
		Score:          player.Score,
		UpdatedAt:      time.Now(),
		ID:             player.ID,
	}

	err = app.store.UpdatePlayerScore(context.Background(), updateScore)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, nil, "Failed to update player record")
		return
	}

	resp := util.CheckAnswerResponse{
		Correct:        isCorrect,
		FunFacts:       destination.FunFacts,
		Trivia:         destination.Trivia,
		CorrectAnswer:  destination.CorrectAnswer,
		CorrectAnswers: player.CorrectAnswers,
		TotalAttempts:  player.TotalAttempts,
		Score:          player.Score,
	}

	sendResponse(w, http.StatusOK, resp, "Answer checked successfully")

}

func (app *App) handleCheckAnswerForGuest(w http.ResponseWriter, r *http.Request) {
	body, err := getBodyWithType[util.CheckAnswerRequest](r)
	if err != nil {
		sendHerrorResponse(w, err)
		return
	}

	isCorrect, destination, err := util.CheckAnswerToQuestionID(app.store, int32(body.QuestionID), body.Answer)
	if err != nil {
		sendHerrorResponse(w, err)
		return
	}
	resp := util.CheckAnswerResponse{
		Correct:       isCorrect,
		FunFacts:      destination.FunFacts,
		Trivia:        destination.Trivia,
		CorrectAnswer: destination.CorrectAnswer,
	}
	sendResponse(w, http.StatusOK, resp, "Answer checked successfully")
}
