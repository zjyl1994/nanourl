package util

import (
	"math/rand"
)

var randCharset = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandString(num int) string {
	lenCharset := len(randCharset)
	result := make([]rune, num)
	for i := 0; i < num; i++ {
		result[i] = randCharset[rand.Intn(lenCharset)]
	}
	return string(result)
}
