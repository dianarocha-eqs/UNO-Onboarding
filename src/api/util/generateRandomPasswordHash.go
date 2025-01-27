package util

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"math/big"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()"

// Generates a random password of a given length.
func GenerateRandomPassword(length int) (string, error) {
	if length <= 0 {
		return "", nil
	}

	password := make([]byte, length)
	for i := range password {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		password[i] = charset[index.Int64()]
	}

	return string(password), nil
}

// Generates the hashed version of the random password (in order to not save the plain text of the user's password, mantaining his security)
func GenerateRandomPasswordHash() (string, string, error) {
	plainPassword, err := GenerateRandomPassword(12) // 12-character password
	if err != nil {
		return "", "", err
	}

	hasher := sha256.New()
	hasher.Write([]byte(plainPassword))
	hashedPassword := hex.EncodeToString(hasher.Sum(nil))

	return plainPassword, hashedPassword, nil
}
