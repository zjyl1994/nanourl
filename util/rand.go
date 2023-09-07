package util

import (
	"math/rand"
)

var randCharset = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandString(num int) string {
	result := make([]rune, num)
	for i := 0; i < num; i++ {
		result[i] = randCharset[rand.Intn(num)]
	}
	return string(result)
}
