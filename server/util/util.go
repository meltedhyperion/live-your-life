package util

import (
	"fmt"
	"math"
	"math/rand"
	"net/url"

	"github.com/meltedhyperion/globetrotter/server/db/pg_db"
)

func PadStringTo(v string, n int) string {
	if len(v) >= n {
		return v[:n]
	}
	return fmt.Sprintf("%-*s", n, v)
}

func GenerateAvatar(seed string) string {
	style := "pixel-art"
	backgroundColors := []string{"b6e3f4", "c0aede", "d1d4f9", "ffd5dc", "ffdfbf"}
	randomBg := backgroundColors[rand.Intn(len(backgroundColors))]
	seedEncoded := url.QueryEscape(seed)
	return fmt.Sprintf("https://api.dicebear.com/7.x/%s/svg?seed=%s&backgroundColor=%s&size=128", style, seedEncoded, randomBg)
}

func GenerateQuestion(destinations []*pg_db.GetRandomDestinationsForQuestionsRow, nameOptions []*pg_db.GetRandomDestinationsRow) []Question {
	questions := make([]Question, 0, 5)
	for _, dest := range destinations {
		correct := fmt.Sprintf("%s, %s", dest.City, dest.Country)

		wrongPool := make([]*pg_db.GetRandomDestinationsRow, len(nameOptions))
		copy(wrongPool, nameOptions)
		rand.Shuffle(len(wrongPool), func(i, j int) {
			wrongPool[i], wrongPool[j] = wrongPool[j], wrongPool[i]
		})
		selectedWrong := wrongPool[:3]
		wrongOptions := make([]string, 0, 3)
		for _, opt := range selectedWrong {
			wrongOptions = append(wrongOptions, fmt.Sprintf("%s, %s", opt.City, opt.Country))
		}

		options := append([]string{correct}, wrongOptions...)
		rand.Shuffle(len(options), func(i, j int) {
			options[i], options[j] = options[j], options[i]
		})

		q := Question{
			QuestionID:    int(dest.ID),
			QuestionHints: dest.Clues,
			AnswerOptions: options,
		}
		questions = append(questions, q)
	}
	return questions
}

func CalculateWilsonScore(correct, total int32) float64 {
	if total == 0 {
		return 0.0
	}
	p := float64(correct) / float64(total)
	z := 1.96
	numerator := p + (z*z)/(2*float64(total)) - z*math.Sqrt((p*(1-p)+z*z/(4*float64(total)))/float64(total))
	denom := 1 + z*z/float64(total)
	return numerator / denom
}
