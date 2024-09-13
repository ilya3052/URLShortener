package random

import (
	"math/rand"
)

func NewRandomString(length int) string {
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ��ÄÖ" +
		"abcdefghijklmnopqrstuvwxyzåäö" +
		"0123456789")
	b := make([]rune, length)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}
