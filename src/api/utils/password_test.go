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

// Test for hashed password function
func TestGenerateRandomPasswordHash(t *testing.T) {
	password, hashedPassword, err := GenerateRandomPasswordHash()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check if password is not empty
	if len(password) == 0 {
		t.Errorf("expected a non-empty password")
	}

	// Check if hashed password is not empty
	if len(hashedPassword) == 0 {
		t.Errorf("expected a non-empty hashed password")
	}

	// Compare password with the hashed version
	if password == hashedPassword {
		t.Errorf("expected hashed password to be different from plain password")
	}

	// Verify if the hash is correct
	hasher := sha256.New()
	hasher.Write([]byte(password))
	expectedHash := hex.EncodeToString(hasher.Sum(nil))

	if hashedPassword != expectedHash {
		t.Errorf("expected hashed password %s, got %s", expectedHash, hashedPassword)
	}
}
