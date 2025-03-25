package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/meltedhyperion/globetrotter/server/util"
)

func HandleQuestionRoutes(app *App) http.Handler {
	r := chi.NewRouter()
	r.Get("/", app.handleGetQuestions)
	r.Post("/check", app.handleCheckAnswer)
	return r
}

func (app *App) handleGetQuestions(w http.ResponseWriter, r *http.Request) {
	var destinations []util.Destination

	res := app.DB.Rpc("get_random_destinations", "", nil)

	destinations, err := util.ParseDestinations(res)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()}, "failed to parse destinations")
		return
	}

	if len(destinations) < 5 {
		http.Error(w, "Not enough destinations available", http.StatusInternalServerError)
		sendErrorResponse(w, http.StatusInternalServerError, nil, "Not enough destinations available")
		return
	}

	var excludeIDs []int
	for _, d := range destinations {
		excludeIDs = append(excludeIDs, d.ID)
	}
	excludeIDsStr := util.ConvertIntSliceToPostgresArray(excludeIDs)

	var nameOptions []util.NameOption
	result, _, err := app.DB.From("destinations").
		Select("city, country", "RANDOM()", false).
		Not("id", "in", excludeIDsStr).
		Limit(15, "").
		Execute()
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()}, "failed to fetch name options")
		return
	}

	err = json.Unmarshal(result, &nameOptions)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, nil, "failed to unmarshal name options")
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

	var destination []util.Destination
	result, _, err := app.DB.From("destinations").
		Select("id, city, country, clues, fun_facts, trivia", "", false).
		Eq("id", strconv.Itoa(body.QuestionID)).
		Execute()
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, nil, "Failed to fetch destination")
		return
	}
	if err := json.Unmarshal(result, &destination); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, nil, "Failed to unmarshal destination")
		return
	}

	if destination[0].ID == 0 {
		sendErrorResponse(w, http.StatusNotFound, nil, "Question not found")
		return
	}

	correctAnswer := fmt.Sprintf("%s, %s", destination[0].City, destination[0].Country)

	isCorrect := (body.Answer == correctAnswer)

	var player []util.Player
	result, _, err = app.DB.From("players").
		Select("id,score,avatar,name,correct_answers,total_attempts", "", false).
		Eq("id", playerID).
		Execute()
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, nil, "Failed to fetch player record")
		return
	}
	if err := json.Unmarshal(result, &player); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, nil, "Failed to unmarshal player record")
		return
	}
	if len(player) == 0 {
		sendErrorResponse(w, http.StatusNotFound, nil, "Player not found")
		return
	}
	var correctAnswers int
	player[0].TotalAttempts++
	if isCorrect {
		player[0].CorrectAnswers++
		correctAnswers = player[0].CorrectAnswers
	}
	player[0].Score = util.CalculateWilsonScore(player[0].CorrectAnswers, player[0].TotalAttempts)
	updateScore := util.UpdatePlayer{
		CorrectAnswers: correctAnswers,
		TotalAttempts:  player[0].TotalAttempts,
		Score:          player[0].Score,
		UpdatedAt:      time.Now(),
	}

	_, _, err = app.DB.From("players").
		Update(updateScore, "", "").
		Eq("id", playerID).
		Execute()
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, nil, "Failed to update player record")
		return
	}

	resp := util.CheckAnswerResponse{
		Correct:        isCorrect,
		FunFacts:       destination[0].FunFacts,
		Trivia:         destination[0].Trivia,
		CorrectAnswer:  correctAnswer,
		CorrectAnswers: player[0].CorrectAnswers,
		TotalAttempts:  player[0].TotalAttempts,
		Score:          player[0].Score,
	}

	sendResponse(w, http.StatusOK, resp, "Answer checked successfully")

}
