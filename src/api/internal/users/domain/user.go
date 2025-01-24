package domain

import (
	"time"
)

type User struct {
	ID        string    `json:"uuid" gorm:"type:uniqueidentifier;primaryKey;default:NEWID()"`
	Name      string    `json:"name" gorm:"type:nvarchar(255);not null"`
	Email     string    `json:"email" gorm:"type:nvarchar(255);unique;not null"`
	Password  string    `json:"password" gorm:"type:nvarchar(64);not null"`
	Picture   string    `json:"picture" gorm:"type:nvarchar(max)"`
	Phone     string    `json:"phone" gorm:"type:nvarchar(20)"`
	Role      bool      `json:"role" gorm:"type:bit;not null;default:0"`
	CreatedAt time.Time `json:"created_at" gorm:"type:datetime;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:datetime;autoUpdateTime"`
}

// func (u *User) MarshalJSON() ([]byte, error) {
// 	type Alias User
// 	return json.Marshal(&struct {
// 		Role string `json:"role"`
// 		*Alias
// 	}{
// 		Role:  boolToString(u.Role),
// 		Alias: (*Alias)(u),
// 	})
// }

// func boolToString(b bool) string {
// 	if b {
// 		return "admin"
// 	}
// 	return "user"
// }
