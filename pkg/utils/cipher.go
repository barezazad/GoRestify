package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"GoRestify/pkg/pkg_err"
)

// Encrypt encrypts text with AES-GCM.
// secret must be 16, 24, or 32 bytes. The nonce is generated and prepended to the ciphertext.
func Encrypt(text, secret string) (encode string, err error) {
	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		err = pkg_err.New(fmt.Sprintf("Failed to created Cipher: %v", err), "E1121827").
			Custom(pkg_err.InternalServerErr).Message(pkg_err.SomethingWentWrong).Build()
		return
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		err = pkg_err.New(fmt.Sprintf("Failed to created GCM: %v", err), "E1121828").
			Custom(pkg_err.InternalServerErr).Message(pkg_err.SomethingWentWrong).Build()
		return
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		err = pkg_err.New(fmt.Sprintf("Failed to generate nonce: %v", err), "E1121829").
			Custom(pkg_err.InternalServerErr).Message(pkg_err.SomethingWentWrong).Build()
		return
	}

	cipherText := gcm.Seal(nonce, nonce, []byte(text), nil)
	encode = base64.StdEncoding.EncodeToString(cipherText)
	return
}

// Decrypt decrypts AES-GCM ciphertext produced by Encrypt.
func Decrypt(text, secret string) (decode string, err error) {
	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		err = pkg_err.New(fmt.Sprintf("Failed to created Cipher: %v", err), "E1169646").
			Custom(pkg_err.InternalServerErr).Message(pkg_err.SomethingWentWrong).Build()
		return
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		err = pkg_err.New(fmt.Sprintf("Failed to created GCM: %v", err), "E1169647").
			Custom(pkg_err.InternalServerErr).Message(pkg_err.SomethingWentWrong).Build()
		return
	}

	cipherText, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		err = pkg_err.New(fmt.Sprintf("Failed to decode text: %v", err), "E1141301").
			Custom(pkg_err.InternalServerErr).Message(pkg_err.SomethingWentWrong).Build()
		return
	}

	nonceSize := gcm.NonceSize()
	if len(cipherText) < nonceSize {
		err = pkg_err.New("cipher text too short", "E1141302").
			Custom(pkg_err.InternalServerErr).Message(pkg_err.SomethingWentWrong).Build()
		return
	}

	nonce, payload := cipherText[:nonceSize], cipherText[nonceSize:]
	plainText, err := gcm.Open(nil, nonce, payload, nil)
	if err != nil {
		err = pkg_err.New(fmt.Sprintf("Failed to decrypt text: %v", err), "E1141303").
			Custom(pkg_err.InternalServerErr).Message(pkg_err.SomethingWentWrong).Build()
		return
	}

	decode = string(plainText)
	return
}
