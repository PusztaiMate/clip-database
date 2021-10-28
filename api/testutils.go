package api

import (
	"math/rand"
	"strings"
)

const (
	stringCharSet = "abcdefghijklmnopqrstuvwy123456789"
)

func getRandomString(n int) string {
	builder := strings.Builder{}
	for i := 0; i < n; i++ {
		builder.WriteByte(stringCharSet[rand.Intn(len(stringCharSet))])
	}
	return builder.String()
}
