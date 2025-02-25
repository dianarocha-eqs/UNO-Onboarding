package domain

import (
	uuid "github.com/tentone/mssql-uuid"
)

// User represents a user in the system with personal and authentication details
type User struct {
	// Unique identifier for the user (UUID)
	ID uuid.UUID `json:"uuid"`
	// User's name (required)
	Name string `json:"name"`
	// User's email (unique and required)
	Email string `json:"email"`
	// User's hashed password
	Password string `json:"password"`
	// Profile picture (optional)
	Picture string `json:"picture"`
	// User's phone number (required)
	Phone string `json:"phone"`
	// User's role: admin (true), user (false)
	Role bool `json:"role"`
}
