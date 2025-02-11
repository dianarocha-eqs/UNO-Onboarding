package utils

import (
	"api/internal/users/domain"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()"

// Generates a random password of a given length.
func GenerateRandomPassword(length int) (string, error) {
	// If length is negative, return an error
	if length < 0 {
		return "", fmt.Errorf("password length cannot be negative")
	}

	// If length is zero, return an empty string and a warning message
	if length == 0 {
		fmt.Println("Warning: password length is 0, returning an empty string")
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
func GeneratePasswordHash(password string) (string, string, error) {
	var plainPassword string

	// If the user provided a password, use it, otherwise generate a random password
	if password == "" {
		var err error
		plainPassword, err = GenerateRandomPassword(12) // 12-character password
		if err != nil {
			return "", "", err
		}
	} else {
		plainPassword = password
	}

	// Hash the password
	hasher := sha256.New()
	hasher.Write([]byte(plainPassword))
	hashedPassword := hex.EncodeToString(hasher.Sum(nil))

	return plainPassword, hashedPassword, nil
}

// Creates the password and sends it to the user via email
func CreatePassword(user *domain.User, password string) (string, error) {
	// Generate password hash (either random or user-provided)
	var plainPassword string
	var hashedPassword string
	var err error

	plainPassword, hashedPassword, err = GeneratePasswordHash(password)
	if err != nil {
		return "", err
	}
	// Assign the hashed password to the user
	user.Password = hashedPassword

	return plainPassword, nil
}
