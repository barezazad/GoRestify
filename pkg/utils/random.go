package utils

import (
	"crypto/rand"
	"math/big"
)

const charset = "LJGDEXvgausbaeavauGTRWJONBDDDEgafakITDRYKmaloajssvrxKDWACJOJKLMNPQRS"

// RandomStringWithCharset create random string using crypto/rand.
func RandomStringWithCharset(length int, charset string) string {
	if length <= 0 || charset == "" {
		return ""
	}

	b := make([]byte, length)
	max := big.NewInt(int64(len(charset)))
	for i := range b {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			// Extremely unlikely; fall back to first charset rune rather than panic.
			b[i] = charset[0]
			continue
		}
		b[i] = charset[n.Int64()]
	}
	return string(b)
}

// RandomString convert random to alphabetic word
func RandomString(length int) string {
	return RandomStringWithCharset(length, charset)
}
