package domain

import (
	"encoding/json"
)

// Sensor represents a device that collects and transmits data about its environment.
type User struct {
	// ID of the users
	ID uint `json:"id"`
	// UUID of the user
	UserID string `gorm:"column:userid" json:"user_id"`
	// Name of the user
	Name string `json:"name"`
	// Email of the user
	Email string `json:"email"`
	// Password of the user
	Password string `json:"password"`
	// Phone number of the user
	Phone string `json:"phoney"`
	// Role of the user
	Role bool `json:"role"`
	// Picture path of the user
	Picture string `json:"picture"`
}

// // Base contains common columns for all tables.
// type Base struct {
// 	ID        string     `gorm:"type:uuid;primary_key;"`
// 	CreatedAt time.Time  `json:"created_at"`
// 	UpdatedAt time.Time  `json:"updated_at"`
// 	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
// }

// // BeforeCreate will set a UUID rather than numeric ID.
// func (b *Base) BeforeCreate(tx *gorm.DB) (err error) {
// 	b.ID = uuid.New().String()
// 	return
// }

func (u *User) MarshalJSON() ([]byte, error) {
	type Alias User
	return json.Marshal(&struct {
		Role string `json:"role"`
		*Alias
	}{
		Role:  boolToString(u.Role),
		Alias: (*Alias)(u),
	})
}

func boolToString(b bool) string {
	if b {
		return "admin"
	}
	return "user"
}
