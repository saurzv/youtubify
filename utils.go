package main

import "math/rand"

func GenerateRandomString(l int) string {
	allowedCharacters := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	s := make([]rune, l)
	for i := range s {
		s[i] = allowedCharacters[rand.Intn(l)]
	}

	return string(s)
}
