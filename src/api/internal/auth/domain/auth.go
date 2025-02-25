package domain

import (
	"time"

	uuid "github.com/tentone/mssql-uuid"
)

// AuthToken represents a JWT token stored in the database
type AuthToken struct {
	// Unique identifier for each token
	ID uuid.UUID `json:"uuid"`
	// Foreign key linking to the Users table
	UserID uuid.UUID `json:"userUuid"`
	// JWT token
	Token string `json:"token"`
	// Indicates if the token is still valid
	IsValid bool `json:"is_valid"`
	// Timestamp for when the token was created
	CreatedAt time.Time `json:"created_at,omitempty"`
	// Timestamp for when the token expired
	ExpiredAt time.Time `json:"expired_at,omitempty"`
}
