package util

import (
	"fmt"
	"math/rand"
	"net/url"
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
