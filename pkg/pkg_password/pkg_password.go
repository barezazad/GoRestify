package pkg_password

import "golang.org/x/crypto/bcrypt"

const bcryptCost = 12

// Hash is used for hashing the password
func Hash(password, salt string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password+salt), bcryptCost)
	return string(bytes), err
}

// Verify is used for verifying the hash and password
func Verify(password, hash, salt string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password+salt))
	return err == nil
}
