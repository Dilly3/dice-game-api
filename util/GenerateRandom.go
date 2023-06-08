package util

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func GenerateRandomString(cap int) string {
	rand.Seed(time.Now().UnixNano())
	alps := "abcdefghijklmnopqrstuvwxyz"
	var result = strings.Builder{}

	for i := 0; i < cap; i++ {
		result.WriteByte(alps[rand.Intn(len(alps))])

	}
	return result.String()
}

func GenerateRandomInt64(cap int64) int64 {
	rand.Seed(time.Now().UnixNano())
	return rand.Int63n(cap)
}
func GenerateRandomInt(cap int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(cap)
}

func GenerateRandomFloat(cap float64) float64 {
	rand.Seed(time.Now().UnixNano())
	return (rand.Float64() * cap) / 100
}

func GenerateRandomEmail(length int) string {
	dom := GenerateRandomString(length)
	return dom + "@gmail.com"
}

func GenerateRandomUsername(length int) string {
	dom := GenerateRandomString(length)
	num := strconv.Itoa(rand.Intn(50))

	return dom + num
}
