package utils

import "math/rand"

func stringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func NewUniqueString(len int) string {
	alpha := "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numericChar := "0123456789"
	specialChar := "-_"
	charset := alpha + numericChar + specialChar
	return stringWithCharset(len, charset)
}
