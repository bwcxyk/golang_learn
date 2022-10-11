package utils

import (
	"log"
	"math/rand"
)

const (
	RandDigit  = "01"
	RandLetter = "02"
	RandMix    = "03"
)

func Rand(length int, t string) string {
	var s string
	switch t {
	case RandMix:
		s = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	case RandDigit:
		s = "0123456789"
	case RandLetter:
		s = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}
	l := len(s)
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		log.Panicf("get random number fail. %v", err)
	}
	for i, v := range b {
		b[i] = s[v%byte(l)]
	}
	return string(b)
}
