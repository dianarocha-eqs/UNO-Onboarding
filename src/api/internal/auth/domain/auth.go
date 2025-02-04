package domain

import (
	"api/internal/users/domain"

	uuid "github.com/tentone/mssql-uuid"
)

// AuthToken represents a JWT token stored in the database
type AuthToken struct {
	UserID  uuid.UUID   `json:"user_id" gorm:"column:user_id;type:uniqueidentifier;not null;index"`
	User    domain.User `json:"user" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	Token   string      `json:"token" gorm:"column:token;type:nvarchar(max);unique;not null"`
	IsValid bool        `json:"is_valid" gorm:"column:is_valid;type:bit;default:1"`
}

// TableName overrides the default GORM table name
func (AuthToken) TableName() string {
	return "user_tokens"
}
