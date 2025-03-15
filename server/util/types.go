package util

import "time"

type CreatePlayerReq struct {
	Name string `json:"name"`
}

type Player struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Avatar         string    `json:"avatar"`
	CorrectAnswers int       `json:"correct_answers"`
	TotalAttempts  int       `json:"total_attempts"`
	Score          float64   `json:"score"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
}

type Destination struct {
	ID        int       `json:"id"`
	City      string    `json:"city"`
	Country   string    `json:"country"`
	Clues     []string  `json:"clues"`
	FunFacts  []string  `json:"fun_facts,omitempty"`
	Trivia    []string  `json:"trivia,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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

type CheckAnswerRequest struct {
	QuestionID int    `json:"question_id"`
	Answer     string `json:"answer"`
}

type CheckAnswerResponse struct {
	Correct        bool     `json:"correct"`
	FunFacts       []string `json:"fun_facts"`
	Trivia         []string `json:"trivia"`
	CorrectAnswer  string   `json:"correct_answer"`
	CorrectAnswers int      `json:"correct_answers"`
	TotalAttempts  int      `json:"total_attempts"`
	Score          float64  `json:"score"`
}

type FriendRequest struct {
	FriendID string `json:"friend_id"`
}

type AddFriend struct {
	Player1ID string `json:"player1_id"`
	Player2ID string `json:"player2_id"`
}

type PlayerStats struct {
	Name           string  `json:"name"`
	Avatar         string  `json:"avatar"`
	Score          float64 `json:"score"`
	CorrectAnswers int     `json:"correct_answers"`
	TotalAttempts  int     `json:"total_attempts"`
}

type Leaderboard struct {
	PlayerStats []PlayerStats `json:"player_stats"`
}

type FriendsIDs struct {
	FriendID string `json:"player2_id"`
}
