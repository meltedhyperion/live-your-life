package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/meltedhyperion/globetrotter/server/util"
)

func HandleQuestionRoutes(app *App) http.Handler {
	r := chi.NewRouter()
	r.Get("/", app.handleGetQuestions)
	return r
}

func (app *App) handleGetQuestions(w http.ResponseWriter, r *http.Request) {
	var destinations []util.Destination

	result, _, err := app.DB.From("destinations").
		Select("*", "RANDOM()", false).Limit(5, "").Execute()

	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()}, "failed to fetch destinations")
		return
	}

	err = json.Unmarshal(result, &destinations)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, nil, "failed to unmarshal destinations")
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
	result, _, err = app.DB.From("destinations").
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
