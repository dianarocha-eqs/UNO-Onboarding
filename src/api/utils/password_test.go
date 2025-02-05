package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

// Test for random password function
func TestGenerateRandomPassword(t *testing.T) {
	tests := []struct {
		length    int
		expectErr bool
	}{
		{12, false}, // test a valid length for password
		{0, false},  // test password with length 0 (print of warning, return empty)
		{-1, false}, // test password with negative length (return an error)
	}

	for _, tt := range tests {
		t.Run("GenerateRandomPassword", func(t *testing.T) {
			password, err := GenerateRandomPassword(tt.length)
			if (err != nil) != tt.expectErr {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
			}

			// Test for non-empty password
			if tt.length > 0 && len(password) != tt.length {
				t.Errorf("expected password length %d, got %d", tt.length, len(password))
			}
		})
	}
}

func TestGeneratePasswordHash(t *testing.T) {
	// Set a random password to test
	password := "randompassword"

	// Test GeneratePasswordHash(password)
	_, hashedPassword, err := GeneratePasswordHash(password)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(hashedPassword) == 0 {
		t.Errorf("expected a non-empty hashed password")
	}

	// Verify if the hash is correct
	hasher := sha256.New()
	hasher.Write([]byte(password))
	expectedHash := hex.EncodeToString(hasher.Sum(nil))

	// Compare the expected hash with the actual hashed password
	if hashedPassword != expectedHash {
		t.Errorf("expected hashed password %s, got %s", expectedHash, hashedPassword)
	}

	// Test for empty password (simulating random password generation)
	emptyPassword := ""

	// Generate the hash using an empty password (this should generate a random password)
	plainPassword, hashedPassword, err := GeneratePasswordHash(emptyPassword)
	if err != nil {
		t.Fatalf("unexpected error for empty password: %v", err)
	}

	if len(plainPassword) == 0 {
		t.Errorf("expected a non-empty random password")
	}

	if len(hashedPassword) == 0 {
		t.Errorf("expected a non-empty hashed password for random password")
	}

	// Verify if the hash is correct
	hasher.Reset()
	hasher.Write([]byte(plainPassword))
	expectedRandomHash := hex.EncodeToString(hasher.Sum(nil))

	// Compare the expected hash with the actual hashed password
	if hashedPassword != expectedRandomHash {
		t.Errorf("expected hashed password %s, got %s for random password", expectedRandomHash, hashedPassword)
	}
}
