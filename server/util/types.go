package util

import (
	"encoding/json"

	"github.com/meltedhyperion/globetrotter/server/db/pg_db"
)

type CreatePlayerReq struct {
	Name string `json:"name"`
}

type Question struct {
	QuestionID    int             `json:"question_id"`
	QuestionHints json.RawMessage `json:"question_hints"`
	AnswerOptions []string        `json:"answer_options"`
}

type CheckAnswerRequest struct {
	QuestionID int    `json:"question_id"`
	Answer     string `json:"answer"`
}

type CheckAnswerResponse struct {
	Correct        bool            `json:"correct"`
	FunFacts       json.RawMessage `json:"fun_facts"`
	Trivia         json.RawMessage `json:"trivia"`
	CorrectAnswer  string          `json:"correct_answer"`
	CorrectAnswers int32           `json:"correct_answers"`
	TotalAttempts  int32           `json:"total_attempts"`
	Score          float64         `json:"score"`
}

type Leaderboard struct {
	PlayerStats []*pg_db.GetLeaderboardForFriendsRow `json:"player_stats"`
}

type FunFactsAndTrivia struct {
	CorrectAnswer string          `json:"correct_answer"`
	FunFacts      json.RawMessage `json:"fun_facts"`
	Trivia        json.RawMessage `json:"trivia"`
}
