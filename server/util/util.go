package util

import (
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
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

func ConvertIntSliceToPostgresArray(slice []int) string {
	var postgresArray string
	for i, id := range slice {
		if i > 0 {
			postgresArray += ","
		}
		postgresArray += strconv.Itoa(id)
	}
	postgresArray = "(" + postgresArray + ")"
	return postgresArray
}

func GenerateQuestion(destinations []Destination, nameOptions []NameOption) []Question {
	questions := make([]Question, 0, 5)
	for _, dest := range destinations {
		correct := fmt.Sprintf("%s, %s", dest.City, dest.Country)

		wrongPool := make([]NameOption, len(nameOptions))
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
			QuestionID:    dest.ID,
			QuestionHints: dest.Clues,
			AnswerOptions: options,
		}
		questions = append(questions, q)
	}
	return questions
}
