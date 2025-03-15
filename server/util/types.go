package util

import "time"

type CreatePlayerReq struct {
	Name   string `json:"name"`
	UserId string `json:"user_id"`
}

type Player struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Avatar         string    `json:"avatar"`
	CorrectAnswers int       `json:"correct_answers"`
	TotalAttempts  int       `json:"total_attempts"`
	Score          float64   `json:"score"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Destination struct {
	ID      int      `json:"id"`
	City    string   `json:"city"`
	Country string   `json:"country"`
	Clues   []string `json:"clues"`
}

type NameOption struct {
	City    string `json:"city"`
	Country string `json:"country"`
}

type Question struct {
	QuestionID    int      `json:"question_id"`
	QuestionHints []string `json:"question_hints"`
	AnswerOptions []string `json:"answer_options"`
}
