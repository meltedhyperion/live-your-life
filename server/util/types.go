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
