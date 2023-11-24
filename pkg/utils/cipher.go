package utils

import (
	"GoRestify/pkg/pkg_err"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

// Encrypt method is to encrypt or hide any classified text
func Encrypt(text, secret string, InitializationVector []byte) (encode string, err error) {

	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		err = pkg_err.New(fmt.Sprintf("Failed to created Cipher: %v", err), "E1121827").
			Custom(pkg_err.InternalServerErr).Message(pkg_err.SomethingWentWrong).Build()
		return
	}

	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, InitializationVector)

	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)

	encode = base64.StdEncoding.EncodeToString(cipherText)

	return
}

// Decrypt method is to extract back the encrypted text
func Decrypt(text, secret string, InitializationVector []byte) (decode string, err error) {

	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		err = pkg_err.New(fmt.Sprintf("Failed to created Cipher: %v", err), "E1169646").
			Custom(pkg_err.InternalServerErr).Message(pkg_err.SomethingWentWrong).Build()
		return
	}

	cipherText, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		err = pkg_err.New(fmt.Sprintf("Failed to decode text: %v", err), "E1141301").
			Custom(pkg_err.InternalServerErr).Message(pkg_err.SomethingWentWrong).Build()
		return
	}

	cfb := cipher.NewCFBDecrypter(block, InitializationVector)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)

	decode = string(plainText)

	return
}
