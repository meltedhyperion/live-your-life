package util

import (
	"encoding/json"
	"math"
	"math/rand"
	"strings"
	"testing"

	"github.com/meltedhyperion/globetrotter/server/db/pg_db"
)

func TestPadStringTo(t *testing.T) {
	tests := []struct {
		input    string
		length   int
		expected string
	}{
		{"hello", 10, "hello     "},
		{"hello world", 5, "hello"},
		{"", 3, "   "},
		{"abc", 3, "abc"},
		{"abcd", 3, "abc"},
	}

	for _, tt := range tests {
		got := PadStringTo(tt.input, tt.length)
		if got != tt.expected {
			t.Errorf("PadStringTo(%q, %d) = %q; want %q", tt.input, tt.length, got, tt.expected)
		}
	}
}

func TestGenerateAvatar(t *testing.T) {
	seed := "test seed"
	avatar := GenerateAvatar(seed)

	expectedPrefix := "https://api.dicebear.com/7.x/pixel-art/svg?seed=test+seed&backgroundColor="
	if !strings.HasPrefix(avatar, expectedPrefix) {
		t.Errorf("GenerateAvatar(%q) = %q; expected to start with %q", seed, avatar, expectedPrefix)
	}
	if !strings.HasSuffix(avatar, "&size=128") {
		t.Errorf("GenerateAvatar(%q) = %q; expected to end with &size=128", seed, avatar)
	}
}

func stringToRawMessage(cluesData []string) json.RawMessage {
	cluesJSON, _ := json.Marshal(cluesData)
	return cluesJSON
}
func TestGenerateQuestion(t *testing.T) {
	rand.Seed(100)
	destinations := []*pg_db.GetRandomDestinationsForQuestionsRow{
		{ID: 1, City: "Paris", Country: "France", Clues: stringToRawMessage([]string{"Eiffel Tower", "Louvre"})},
		{ID: 2, City: "Tokyo", Country: "Japan", Clues: stringToRawMessage([]string{"Sushi", "Cherry Blossoms"})},
	}
	nameOptions := []*pg_db.GetRandomDestinationsRow{
		{City: "Paris", Country: "France"},
		{City: "London", Country: "UK"},
		{City: "Berlin", Country: "Germany"},
		{City: "Madrid", Country: "Spain"},
		{City: "Rome", Country: "Italy"},
	}

	questions := GenerateQuestion(destinations, nameOptions)
	if len(questions) != len(destinations) {
		t.Errorf("Expected %d questions, got %d", len(destinations), len(questions))
	}

	for _, q := range questions {
		var correct string
		for _, dest := range destinations {
			if dest.ID == int32(q.QuestionID) {
				correct = dest.City + ", " + dest.Country
				break
			}
		}
		found := false
		for _, option := range q.AnswerOptions {
			if option == correct {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Question with ID %d: correct option %q not found in options %v", q.QuestionID, correct, q.AnswerOptions)
		}
		if len(q.AnswerOptions) != 4 {
			t.Errorf("Question with ID %d: expected 4 options, got %d", q.QuestionID, len(q.AnswerOptions))
		}
		if len(q.QuestionHints) == 0 {
			t.Errorf("Question with ID %d: expected non-empty hints", q.QuestionID)
		}
	}
}

func TestCalculateWilsonScore(t *testing.T) {
	score := CalculateWilsonScore(0, 0)
	if score != 0.0 {
		t.Errorf("CalculateWilsonScore(0, 0) = %f; want 0.0", score)
	}

	score = CalculateWilsonScore(3, 10)
	expected := 0.107789
	if math.Abs(score-expected) > 1e-5 {
		t.Errorf("CalculateWilsonScore(3, 10) = %f; want approx %f", score, expected)
	}

	score = CalculateWilsonScore(5, 5)
	expected = 0.565509
	if math.Abs(score-expected) > 1e-5 {
		t.Errorf("CalculateWilsonScore(5, 5) = %f; want approx %f", score, expected)
	}

	score = CalculateWilsonScore(0, 5)
	expected = 0.0
	if math.Abs(score-expected) > 1e-5 {
		t.Errorf("CalculateWilsonScore(0, 5) = %f; want %f", score, expected)
	}
}
