package domain

const (
	// Normal user
	ROLE_USER bool = false
	// Admin user
	ROLE_ADMIN bool = true
)

// User represents a user in the system with personal and authentication details
type User struct {
	// Unique identifier for the user (UUID)
	ID string `json:"uuid" gorm:"column:id;type:uniqueidentifier"`
	// User's name (required)
	Name string `json:"name" gorm:"type:nvarchar(255);not null"`
	// User's email (unique and required)
	Email string `json:"email" gorm:"type:nvarchar(255);unique;not null"`
	// User's hashed password
	Password string `json:"password" gorm:"type:nvarchar(64);not null"`
	// Profile picture (optional)
	Picture string `json:"picture" gorm:"type:nvarchar(max)"`
	// User's phone number (required)
	Phone string `json:"phone" gorm:"type:nvarchar(20)"`
	// User's role (admin or regular user)
	Role bool `json:"role" gorm:"type:bit;not null;default:0"`
}
