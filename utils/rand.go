package utils

import (
	"math/rand"
	"time"
)

const letters = "abcdef0123456789"

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func RandStr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
