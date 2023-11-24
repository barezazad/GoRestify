package utils

import (
	"math/rand"
	"time"
)

const charset = "LJGDEXvgausbaeavauGTRWJONBDDDEgafakITDRYKmaloajssvrxKDWACJOJKLMNPQRS"

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

// RandomStringWithCharset create random string
func RandomStringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// RandomString convert random to alphabetic word
func RandomString(length int) string {
	return RandomStringWithCharset(length, charset)
}
