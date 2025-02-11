package domain

import (
	"api/internal/users/domain"

	uuid "github.com/tentone/mssql-uuid"
)

// AuthToken represents a JWT token stored in the database
type AuthToken struct {
	// Foreign key that links to the User table.
	UserID uuid.UUID `json:"user_id" gorm:"column:user_id;type:uniqueidentifier;not null;index"`
	// This field establishes a relationship between AuthToken and User using the foreign key.
	User domain.User `json:"user" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	// Token is the JWT token string associated with the user
	Token string `json:"token" gorm:"column:token;type:nvarchar(max);unique;not null"`
	// IsValid is a flag that indicates whether the token is still valid.
	IsValid bool `json:"is_valid" gorm:"column:is_valid;type:bit;default:1"`
}
